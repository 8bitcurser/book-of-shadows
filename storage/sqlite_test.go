package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"book-of-shadows/internal/config"
	"book-of-shadows/internal/errors"
)

// testConfig returns a config for testing with a temporary database
func testConfig(t *testing.T) *config.DatabaseConfig {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()

	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	return &config.DatabaseConfig{
		Path:            tmpFile.Name(),
		RetentionPeriod: 24 * time.Hour,
		CleanupInterval: 1 * time.Hour,
	}
}

func TestNewSQLiteStore(t *testing.T) {
	t.Run("creates store with valid config", func(t *testing.T) {
		cfg := testConfig(t)
		store, err := NewSQLiteStore(cfg)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		defer store.Close()

		if store.db == nil {
			t.Error("expected db to be initialized")
		}
	})

	t.Run("returns error with nil config", func(t *testing.T) {
		_, err := NewSQLiteStore(nil)
		if err == nil {
			t.Error("expected error with nil config")
		}
	})

	t.Run("returns error with invalid path", func(t *testing.T) {
		cfg := &config.DatabaseConfig{
			Path:            "/nonexistent/path/to/db.db",
			RetentionPeriod: 24 * time.Hour,
			CleanupInterval: 1 * time.Hour,
		}
		_, err := NewSQLiteStore(cfg)
		if err == nil {
			t.Error("expected error with invalid path")
		}
	})
}

func TestSaveExport(t *testing.T) {
	cfg := testConfig(t)
	store, err := NewSQLiteStore(cfg)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	t.Run("saves valid data", func(t *testing.T) {
		data := `{"test": "data"}`
		id, err := store.SaveExport(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if id == "" {
			t.Error("expected non-empty ID")
		}
	})

	t.Run("returns error for empty data", func(t *testing.T) {
		_, err := store.SaveExport("")
		if err != errors.ErrInvalidData {
			t.Errorf("expected ErrInvalidData, got %v", err)
		}
	})

	t.Run("generates unique IDs", func(t *testing.T) {
		id1, _ := store.SaveExport("data1")
		id2, _ := store.SaveExport("data2")
		if id1 == id2 {
			t.Error("expected unique IDs")
		}
	})
}

func TestGetExport(t *testing.T) {
	cfg := testConfig(t)
	store, err := NewSQLiteStore(cfg)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	t.Run("retrieves saved data", func(t *testing.T) {
		data := `{"name": "test"}`
		id, err := store.SaveExport(data)
		if err != nil {
			t.Fatalf("failed to save: %v", err)
		}

		retrieved, err := store.GetExport(id)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if retrieved != data {
			t.Errorf("expected %q, got %q", data, retrieved)
		}
	})

	t.Run("returns error for non-existent ID", func(t *testing.T) {
		_, err := store.GetExport("nonexistent-id")
		if err != errors.ErrNotFound {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("returns error for empty ID", func(t *testing.T) {
		_, err := store.GetExport("")
		if err != errors.ErrInvalidData {
			t.Errorf("expected ErrInvalidData, got %v", err)
		}
	})
}

func TestDeleteExpiredExports(t *testing.T) {
	cfg := testConfig(t)
	cfg.RetentionPeriod = 1 * time.Millisecond // Very short for testing
	store, err := NewSQLiteStore(cfg)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	t.Run("deletes expired exports", func(t *testing.T) {
		// Save some data
		id, err := store.SaveExport("test data")
		if err != nil {
			t.Fatalf("failed to save: %v", err)
		}

		// Wait for it to expire
		time.Sleep(10 * time.Millisecond)

		// Delete expired
		if err := store.DeleteExpiredExports(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Should be gone
		_, err = store.GetExport(id)
		if err != errors.ErrNotFound {
			t.Errorf("expected ErrNotFound after deletion, got %v", err)
		}
	})
}

func TestPing(t *testing.T) {
	cfg := testConfig(t)
	store, err := NewSQLiteStore(cfg)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	t.Run("returns nil for healthy connection", func(t *testing.T) {
		err := store.Ping(context.Background())
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}

func TestClose(t *testing.T) {
	cfg := testConfig(t)
	store, err := NewSQLiteStore(cfg)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	t.Run("closes cleanly", func(t *testing.T) {
		err := store.Close()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("subsequent operations fail", func(t *testing.T) {
		// After closing, SaveExport should fail
		_, err := store.SaveExport("test")
		if err == nil {
			t.Error("expected error after close")
		}
	})
}
