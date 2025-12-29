package storage

import (
	"book-of-shadows/internal/config"
	"fmt"
)

// AppStore combines all storage functionality
type AppStore struct {
	*SQLiteStore
	*CookieStore
}

// NewAppStore creates a new combined store instance
func NewAppStore(cfg *config.Config) (*AppStore, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	// Create SQLite store
	sqliteStore, err := NewSQLiteStore(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to create SQLite store: %w", err)
	}

	// Create cookie store with SQLite as export store
	cookieStore := NewCookieStore(&cfg.Cookie, sqliteStore)

	return &AppStore{
		SQLiteStore: sqliteStore,
		CookieStore: cookieStore,
	}, nil
}

// Close gracefully shuts down the store
func (s *AppStore) Close() error {
	return s.SQLiteStore.Close()
}

// Ensure AppStore implements the Store interface
var _ Store = (*AppStore)(nil)

// The AppStore now implements all methods from both SQLiteStore and CookieStore:
// From SQLiteStore (ExportStore):
// - SaveExport(data string) (string, error)
// - GetExport(id string) (string, error)
// - DeleteExpiredExports() error
//
// From CookieStore (InvestigatorStore):
// - SaveInvestigator(w http.ResponseWriter, inv *models.Investigator) (string, error)
// - GetInvestigator(r *http.Request, id string) (*models.Investigator, error)
// - UpdateInvestigator(w http.ResponseWriter, id string, inv *models.Investigator) error
// - DeleteInvestigator(w http.ResponseWriter, id string) error
// - ListInvestigators(r *http.Request) (map[string]*models.Investigator, error)
// - ExportInvestigatorsList(r *http.Request) (string, error)
// - ImportInvestigatorsList(w http.ResponseWriter, uuid string) error