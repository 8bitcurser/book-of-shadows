package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"book-of-shadows/internal/config"
	"book-of-shadows/internal/handlers"
	"book-of-shadows/internal/middleware"
	"book-of-shadows/storage"
	"book-of-shadows/wizard"
)

// Server represents the application server with all dependencies
type Server struct {
	config   *config.Config
	store    *storage.AppStore
	handlers *handlers.Handler
	wizard   *wizard.Handler
	logger   *log.Logger
}

// NewServer creates a new server instance with all dependencies
func NewServer() (*Server, error) {
	// Load configuration
	cfg := config.New()

	// Setup logger
	logger := log.New(os.Stdout, "[book-of-shadows] ", log.LstdFlags|log.Lshortfile)

	// Create store
	store, err := storage.NewAppStore(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	// Create handlers with dependencies
	h := handlers.New(store, logger)

	// Create wizard handler with dependencies
	wizardHandler := wizard.New(store, logger)

	return &Server{
		config:   cfg,
		store:    store,
		handlers: h,
		wizard:   wizardHandler,
		logger:   logger,
	}, nil
}

// setupRoutes configures all application routes
func (s *Server) setupRoutes() *RadixTree {
	router := NewRouter()

	// Static files
	fs := http.FileServer(http.Dir("static"))
	fsHandler := http.StripPrefix("/static/", fs)
	router.HandleStatic("/static/", fsHandler)

	// Main routes
	router.GET("/", s.handlers.Home)
	router.GET("api/generate/", s.handlers.Generate)

	// Investigator CRUD operations
	router.GET("api/investigator", s.handlers.ListInvestigators)
	router.POST("api/investigator/", s.handlers.CreateInvestigator)
	router.GET("api/investigator/{:id}", s.handlers.GetInvestigator)
	router.PUT("api/investigator/{:id}", s.handlers.UpdateInvestigator)
	router.DELETE("api/investigator/{:id}", s.handlers.DeleteInvestigator)

	// Export/Import operations
	router.POST("api/investigator/PDF/{:id}", s.handlers.ExportPDF)
	router.GET("api/investigator/list/export", s.handlers.ExportInvestigatorsList)
	router.POST("api/investigator/list/import/", s.handlers.ImportInvestigatorsList)

	// Other routes
	router.GET("api/archetype/{:name}/occupations/", s.handlers.GetArchetypeOccupations)
	router.POST("api/report-issue", s.handlers.ReportIssue)

	// Wizard routes
	router.GET("wizard/base/{:key}", s.wizard.BaseStep)
	router.GET("wizard/attributes/{:key}", s.wizard.AttrStep)
	router.GET("wizard/skills/{:key}", s.wizard.SkillStep)

	return router
}

// Run starts the HTTP server with graceful shutdown support
func (s *Server) Run() error {
	router := s.setupRoutes()

	// Apply middleware chain
	handler := middleware.Chain(
		router,
		middleware.Recoverer(s.logger),
		middleware.Logger(s.logger),
		middleware.SecurityHeaders,
		middleware.RequestID,
	)

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         ":" + s.config.Server.Port,
		Handler:      handler,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
		IdleTimeout:  s.config.Server.IdleTimeout,
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run server in goroutine
	go func() {
		s.logger.Printf("Server starting on port %s", s.config.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	s.logger.Println("Server shutting down...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	// Close store connections
	if err := s.store.Close(); err != nil {
		return fmt.Errorf("failed to close store: %w", err)
	}

	s.logger.Println("Server stopped gracefully")
	return nil
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
