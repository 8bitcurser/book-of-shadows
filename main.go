package main

import (
	"book-of-shadows/storage"
	"log"
	"net/http"
)

type Route struct {
	handlers map[string]http.HandlerFunc
}

type Router struct {
	routes map[string]*Route
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]*Route),
	}
}

func (r *Router) Handle(method, path string, handler http.HandlerFunc) {
	if r.routes[path] == nil {
		r.routes[path] = &Route{
			handlers: make(map[string]http.HandlerFunc),
		}
	}
	r.routes[path].handlers[method] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if route, exists := r.routes[req.URL.Path]; exists {
		if handler, methodExists := route.handlers[req.Method]; methodExists {
			handler(w, req)
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
	http.NotFound(w, req)
}

func (r *Router) GET(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodPut, path, handler)
}

func (r *Router) DELETE(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodDelete, path, handler)
}

func setupRoutes() *Router {
	router := NewRouter()

	// Define routes
	router.GET("/", handleHome)
	router.GET("/api/generate", handleGenerate)

	// Investigator CRUD operations
	router.GET("/api/investigator", handleListInvestigators)
	router.GET("/api/investigator/", handleGetInvestigator)
	router.POST("/api/investigator/", handleCreateBaseInvestigator)
	router.PUT("/api/investigator/", handleUpdateInvestigator)
	router.DELETE("/api/investigator/", handleDeleteInvestigator)

	router.GET("/api/investigator/PDF/", handleExportPDF)
	router.GET("/api/investigator/list/export", handleListInvestigatorsExport)

	router.POST("/api/investigator/list/import/", handleListInvestigatorsImport)
	router.POST("/api/generate-step/", handleCreateStepInvestigator)
	router.POST("/api/investigator/confirm-attributes/", handleConfirmAttrStepInvestigator)

	router.GET("/api/archetype/", handleArchetypeOccupations)
	router.POST("/api/report-issue", handleReportIssue)

	return router
}

func main() {
	conn := storage.SQLiteDB{}
	conn.Init()
	conn.StartCleanupRoutine()
	defer conn.DB.Close()

	router := setupRoutes()
	// routes
	// http.HandleFunc("/", handleHome)
	// http.HandleFunc("/api/generate", handleGenerate)
	// http.HandleFunc("/api/investigator/export/", handleExportPDF)
	// http.HandleFunc("/api/investigator/list", handleListInvestigators)
	// http.HandleFunc("/api/investigator/list/export", handleListInvestigatorsExport)
	// http.HandleFunc("/api/investigator/list/import/", handleListInvestigatorsImport)
	// http.HandleFunc("/api/generate-step/", handleCreateStepInvestigator)
	// http.HandleFunc("/api/investigator/confirm-attributes/", handleConfirmAttrStepInvestigator)
	// http.HandleFunc("/api/investigator/create/", handleCreateBaseInvestigator)
	// http.HandleFunc("/api/investigator/delete/", handleDeleteInvestigator)
	// http.HandleFunc("/api/investigator/update/", handleUpdateInvestigator)
	// http.HandleFunc("/api/investigator/", handleGetInvestigator)
	// http.HandleFunc("/api/archetype/", handleArchetypeOccupations)
	// http.HandleFunc("/api/report-issue", handleReportIssue)
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
