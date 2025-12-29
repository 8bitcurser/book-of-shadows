package storage

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"book-of-shadows/internal/config"
	"book-of-shadows/internal/errors"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStore implements the ExportStore interface using SQLite
type SQLiteStore struct {
	db     *sql.DB
	config *config.DatabaseConfig
	mu     sync.RWMutex

	// Cleanup goroutine management
	cleanupCtx    context.Context
	cleanupCancel context.CancelFunc
}

// NewSQLiteStore creates a new SQLiteStore instance
func NewSQLiteStore(cfg *config.DatabaseConfig) (*SQLiteStore, error) {
	if cfg == nil {
		return nil, fmt.Errorf("database config is required")
	}

	db, err := sql.Open("sqlite3", cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(1) // SQLite doesn't handle concurrent writes well
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0) // Connections don't expire

	store := &SQLiteStore{
		db:     db,
		config: cfg,
	}

	if err := store.initialize(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Start cleanup routine
	store.startCleanupRoutine()

	return store, nil
}

// initialize creates the necessary tables
func (s *SQLiteStore) initialize() error {
	query := `
		CREATE TABLE IF NOT EXISTS exports (
			id TEXT PRIMARY KEY,
			data TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_exports_created_at ON exports(created_at);
	`

	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

// SaveExport saves export data and returns a unique ID
func (s *SQLiteStore) SaveExport(data string) (string, error) {
	if data == "" {
		return "", errors.ErrInvalidData
	}

	id := uuid.New().String()

	query := `INSERT INTO exports (id, data, created_at) VALUES (?, ?, ?)`
	if _, err := s.db.Exec(query, id, data, time.Now()); err != nil {
		return "", fmt.Errorf("failed to save export: %w", err)
	}

	return id, nil
}

// GetExport retrieves export data by ID
func (s *SQLiteStore) GetExport(id string) (string, error) {
	if id == "" {
		return "", errors.ErrInvalidData
	}

	var data string
	query := `SELECT data FROM exports WHERE id = ?`

	err := s.db.QueryRow(query, id).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.ErrNotFound
		}
		return "", fmt.Errorf("failed to get export: %w", err)
	}

	return data, nil
}

// DeleteExpiredExports removes exports older than the retention period
func (s *SQLiteStore) DeleteExpiredExports() error {
	cutoff := time.Now().Add(-s.config.RetentionPeriod)

	query := `DELETE FROM exports WHERE created_at < ?`
	result, err := s.db.Exec(query, cutoff)
	if err != nil {
		return fmt.Errorf("failed to delete expired exports: %w", err)
	}

	if rowsAffected, err := result.RowsAffected(); err == nil && rowsAffected > 0 {
		// Log cleanup for monitoring (you might want to use a proper logger)
		fmt.Printf("Cleaned up %d expired exports\n", rowsAffected)
	}

	return nil
}

// startCleanupRoutine starts a background goroutine to clean up expired exports
func (s *SQLiteStore) startCleanupRoutine() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create cancellable context
	s.cleanupCtx, s.cleanupCancel = context.WithCancel(context.Background())

	go func() {
		ticker := time.NewTicker(s.config.CleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := s.DeleteExpiredExports(); err != nil {
					// Log error (you might want to use a proper logger)
					fmt.Printf("cleanup error: %v\n", err)
				}
			case <-s.cleanupCtx.Done():
				return
			}
		}
	}()
}

// Close gracefully shuts down the store
func (s *SQLiteStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Cancel cleanup routine
	if s.cleanupCancel != nil {
		s.cleanupCancel()
	}

	// Close database connection
	if s.db != nil {
		return s.db.Close()
	}

	return nil
}

// Ping checks if the database is accessible
func (s *SQLiteStore) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return s.db.PingContext(ctx)
}