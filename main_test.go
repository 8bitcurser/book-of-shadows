package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"book-of-shadows/internal/errors"
	"book-of-shadows/internal/handlers"
	"book-of-shadows/models"
)

// TestServer wraps the routing and handler setup for integration tests
type TestServer struct {
	router   *RadixTree
	handlers *handlers.Handler
	store    *MockAppStore
}

// MockAppStore implements a complete in-memory store for integration testing
type MockAppStore struct {
	investigators map[string]*models.Investigator
	exports       map[string]string
}

func NewMockAppStore() *MockAppStore {
	return &MockAppStore{
		investigators: make(map[string]*models.Investigator),
		exports:       make(map[string]string),
	}
}

// ExportStore methods
func (m *MockAppStore) SaveExport(data string) (string, error) {
	id := "test-export-id"
	m.exports[id] = data
	return id, nil
}

func (m *MockAppStore) GetExport(id string) (string, error) {
	data, ok := m.exports[id]
	if !ok {
		return "", nil
	}
	return data, nil
}

func (m *MockAppStore) DeleteExpiredExports() error {
	return nil
}

// InvestigatorStore methods
func (m *MockAppStore) SaveInvestigator(w http.ResponseWriter, inv *models.Investigator) (string, error) {
	id := "test-inv-id"
	inv.ID = id
	m.investigators[id] = inv
	return id, nil
}

func (m *MockAppStore) GetInvestigator(r *http.Request, id string) (*models.Investigator, error) {
	inv, ok := m.investigators[id]
	if !ok {
		return nil, errors.ErrNotFound
	}
	return inv, nil
}

func (m *MockAppStore) UpdateInvestigator(w http.ResponseWriter, id string, inv *models.Investigator) error {
	if _, ok := m.investigators[id]; !ok {
		return errors.ErrNotFound
	}
	m.investigators[id] = inv
	return nil
}

func (m *MockAppStore) DeleteInvestigator(w http.ResponseWriter, id string) error {
	if _, ok := m.investigators[id]; !ok {
		return errors.ErrNotFound
	}
	delete(m.investigators, id)
	return nil
}

func (m *MockAppStore) ListInvestigators(r *http.Request) (map[string]*models.Investigator, error) {
	return m.investigators, nil
}

func (m *MockAppStore) ExportInvestigatorsList(r *http.Request) (string, error) {
	return "test-export-code", nil
}

func (m *MockAppStore) ImportInvestigatorsList(w http.ResponseWriter, uuid string) error {
	return nil
}

// Close is a no-op for the mock store
func (m *MockAppStore) Close() error {
	return nil
}

// Setup test server
func newTestServer() *TestServer {
	store := NewMockAppStore()
	logger := log.New(io.Discard, "", 0)
	h := handlers.New(store, logger)

	router := NewRouter()

	// Register routes matching main.go
	router.GET("api/investigator", h.ListInvestigators)
	router.POST("api/investigator/", h.CreateInvestigator)
	router.GET("api/investigator/{:id}", h.GetInvestigator)
	router.PUT("api/investigator/{:id}", h.UpdateInvestigator)
	router.DELETE("api/investigator/{:id}", h.DeleteInvestigator)
	router.GET("api/investigator/list/export", h.ExportInvestigatorsList)
	router.POST("api/investigator/list/import/", h.ImportInvestigatorsList)
	router.GET("api/archetype/{:name}/occupations/", h.GetArchetypeOccupations)

	return &TestServer{
		router:   router,
		handlers: h,
		store:    store,
	}
}

func TestIntegrationCreateAndGetInvestigator(t *testing.T) {
	ts := newTestServer()

	t.Run("create investigator returns 201", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":       "Test Investigator",
			"age":        "30",
			"residence":  "Boston",
			"birthplace": "New York",
			"archetype":  "Adventurer",
			"occupation": "Antiquarian",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/investigator/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ts.router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
		}

		var result map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if result["Key"] == "" {
			t.Error("expected non-empty Key in response")
		}
	})

	t.Run("get investigator returns 200", func(t *testing.T) {
		// Pre-populate the store
		inv := models.RandomInvestigator(models.Pulp)
		inv.ID = "existing-id"
		ts.store.investigators["existing-id"] = inv

		req := httptest.NewRequest("GET", "/api/investigator/existing-id", nil)
		w := httptest.NewRecorder()

		ts.router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("get non-existent investigator returns 404", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/investigator/nonexistent", nil)
		w := httptest.NewRecorder()

		ts.router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestIntegrationUpdateInvestigator(t *testing.T) {
	ts := newTestServer()

	t.Run("update investigator returns 200", func(t *testing.T) {
		// Pre-populate the store
		inv := models.RandomInvestigator(models.Pulp)
		inv.ID = "test-update-id"
		ts.store.investigators["test-update-id"] = inv

		payload := map[string]interface{}{
			"section": "personalInfo",
			"field":   "Name",
			"value":   "Updated Name",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("PUT", "/api/investigator/test-update-id", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ts.router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("update non-existent investigator returns 404", func(t *testing.T) {
		payload := map[string]interface{}{
			"section": "personalInfo",
			"field":   "Name",
			"value":   "Updated Name",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("PUT", "/api/investigator/nonexistent", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ts.router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestIntegrationDeleteInvestigator(t *testing.T) {
	ts := newTestServer()

	t.Run("delete investigator returns 200", func(t *testing.T) {
		// Pre-populate the store
		inv := models.RandomInvestigator(models.Pulp)
		inv.ID = "test-delete-id"
		ts.store.investigators["test-delete-id"] = inv

		req := httptest.NewRequest("DELETE", "/api/investigator/test-delete-id", nil)
		w := httptest.NewRecorder()

		ts.router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}

		// Verify it was deleted
		if _, exists := ts.store.investigators["test-delete-id"]; exists {
			t.Error("expected investigator to be deleted")
		}
	})

	t.Run("delete non-existent investigator returns 404", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/investigator/nonexistent", nil)
		w := httptest.NewRecorder()

		ts.router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestIntegrationListInvestigators(t *testing.T) {
	ts := newTestServer()

	t.Run("list investigators returns 200", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/investigator", nil)
		w := httptest.NewRecorder()

		ts.router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("list with investigators returns data", func(t *testing.T) {
		// Add some investigators
		inv1 := models.RandomInvestigator(models.Pulp)
		inv1.ID = "inv-1"
		inv1.Name = "Investigator 1"
		ts.store.investigators["inv-1"] = inv1

		inv2 := models.RandomInvestigator(models.Pulp)
		inv2.ID = "inv-2"
		inv2.Name = "Investigator 2"
		ts.store.investigators["inv-2"] = inv2

		req := httptest.NewRequest("GET", "/api/investigator", nil)
		w := httptest.NewRecorder()

		ts.router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestIntegrationExportImport(t *testing.T) {
	ts := newTestServer()

	t.Run("export investigators list returns 200", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/investigator/list/export", nil)
		w := httptest.NewRecorder()

		ts.router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("import investigators list returns 201", func(t *testing.T) {
		payload := map[string]string{"ImportCode": "test-code"}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/investigator/list/import/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ts.router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
		}
	})

	t.Run("import without code returns 400", func(t *testing.T) {
		payload := map[string]string{}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/investigator/list/import/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		ts.router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestIntegrationArchetypeOccupations(t *testing.T) {
	ts := newTestServer()

	t.Run("get archetype occupations returns 200", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/archetype/Adventurer/occupations/", nil)
		w := httptest.NewRecorder()

		ts.router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}
	})
}

func TestIntegrationInvalidRoutes(t *testing.T) {
	ts := newTestServer()

	t.Run("invalid route returns 404", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/invalid/route", nil)
		w := httptest.NewRecorder()

		ts.router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("wrong method returns 404", func(t *testing.T) {
		req := httptest.NewRequest("PATCH", "/api/investigator/test-id", nil)
		w := httptest.NewRecorder()

		ts.router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}
