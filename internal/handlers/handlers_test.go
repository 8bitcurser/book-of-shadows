package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	stderrors "errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"book-of-shadows/internal/errors"
	"book-of-shadows/models"
)

// MockStore implements storage.Store for testing
type MockStore struct {
	investigators map[string]*models.Investigator
	exports       map[string]string
	saveError     error
	getError      error
}

func NewMockStore() *MockStore {
	return &MockStore{
		investigators: make(map[string]*models.Investigator),
		exports:       make(map[string]string),
	}
}

// ExportStore methods
func (m *MockStore) SaveExport(data string) (string, error) {
	if m.saveError != nil {
		return "", m.saveError
	}
	id := "test-export-id"
	m.exports[id] = data
	return id, nil
}

func (m *MockStore) GetExport(id string) (string, error) {
	if m.getError != nil {
		return "", m.getError
	}
	data, ok := m.exports[id]
	if !ok {
		return "", errors.ErrNotFound
	}
	return data, nil
}

func (m *MockStore) DeleteExpiredExports() error {
	return nil
}

// InvestigatorStore methods
func (m *MockStore) SaveInvestigator(w http.ResponseWriter, inv *models.Investigator) (string, error) {
	if m.saveError != nil {
		return "", m.saveError
	}
	id := "test-inv-id"
	inv.ID = id
	m.investigators[id] = inv
	return id, nil
}

func (m *MockStore) GetInvestigator(r *http.Request, id string) (*models.Investigator, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	inv, ok := m.investigators[id]
	if !ok {
		return nil, errors.ErrNotFound
	}
	return inv, nil
}

func (m *MockStore) UpdateInvestigator(w http.ResponseWriter, id string, inv *models.Investigator) error {
	if _, ok := m.investigators[id]; !ok {
		return errors.ErrNotFound
	}
	m.investigators[id] = inv
	return nil
}

func (m *MockStore) DeleteInvestigator(w http.ResponseWriter, id string) error {
	if _, ok := m.investigators[id]; !ok {
		return errors.ErrNotFound
	}
	delete(m.investigators, id)
	return nil
}

func (m *MockStore) ListInvestigators(r *http.Request) (map[string]*models.Investigator, error) {
	return m.investigators, nil
}

func (m *MockStore) ExportInvestigatorsList(r *http.Request) (string, error) {
	return "test-export-code", nil
}

func (m *MockStore) ImportInvestigatorsList(w http.ResponseWriter, uuid string) error {
	return nil
}

// Helper to create a test handler
func newTestHandler() (*Handler, *MockStore) {
	store := NewMockStore()
	logger := log.New(io.Discard, "", 0)
	return New(store, logger), store
}

// Helper to create a request with params context
func requestWithParams(method, path string, body []byte, params []string) *http.Request {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}
	req := httptest.NewRequest(method, path, bodyReader)
	if params != nil {
		ctx := context.WithValue(req.Context(), "params", params)
		req = req.WithContext(ctx)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestNew(t *testing.T) {
	t.Run("creates handler with store and logger", func(t *testing.T) {
		store := NewMockStore()
		logger := log.New(io.Discard, "", 0)
		h := New(store, logger)
		if h == nil {
			t.Error("expected non-nil handler")
		}
	})

	t.Run("creates handler with nil logger", func(t *testing.T) {
		store := NewMockStore()
		h := New(store, nil)
		if h == nil {
			t.Error("expected non-nil handler")
		}
	})
}

func TestRespondJSON(t *testing.T) {
	h, _ := newTestHandler()

	t.Run("sends JSON response", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"key": "value"}

		h.respondJSON(w, http.StatusOK, data)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", contentType)
		}

		var result map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if result["key"] != "value" {
			t.Errorf("expected key=value, got %s", result["key"])
		}
	})
}

func TestRespondError(t *testing.T) {
	h, _ := newTestHandler()

	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{"not found error", errors.ErrNotFound, http.StatusNotFound},
		{"cookie not found", errors.ErrCookieNotFound, http.StatusNotFound},
		{"invalid data", errors.ErrInvalidData, http.StatusBadRequest},
		{"invalid attribute", errors.ErrInvalidAttribute, http.StatusBadRequest},
		{"already exists", errors.ErrAlreadyExists, http.StatusConflict},
		{"cookie too large", errors.ErrCookieTooLarge, http.StatusRequestEntityTooLarge},
		{"generic error", stderrors.New("unknown error"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			h.respondError(w, tt.err)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}
			if _, ok := result["error"]; !ok {
				t.Error("expected error field in response")
			}
		})
	}
}

func TestDeleteInvestigator(t *testing.T) {
	t.Run("deletes existing investigator", func(t *testing.T) {
		h, store := newTestHandler()

		// Create an investigator first
		inv := &models.Investigator{ID: "test-id", Name: "Test"}
		store.investigators["test-id"] = inv

		req := requestWithParams("DELETE", "/api/investigator/test-id", nil, []string{"test-id"})
		w := httptest.NewRecorder()

		h.DeleteInvestigator(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		if _, exists := store.investigators["test-id"]; exists {
			t.Error("expected investigator to be deleted")
		}
	})

	t.Run("returns error for missing params", func(t *testing.T) {
		h, _ := newTestHandler()

		req := requestWithParams("DELETE", "/api/investigator/", nil, []string{})
		w := httptest.NewRecorder()

		h.DeleteInvestigator(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("returns error for non-existent investigator", func(t *testing.T) {
		h, _ := newTestHandler()

		req := requestWithParams("DELETE", "/api/investigator/nonexistent", nil, []string{"nonexistent"})
		w := httptest.NewRecorder()

		h.DeleteInvestigator(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestCreateInvestigator(t *testing.T) {
	t.Run("creates investigator with valid data", func(t *testing.T) {
		h, _ := newTestHandler()

		payload := map[string]interface{}{
			"name":       "Test Investigator",
			"age":        "30",
			"residence":  "Boston",
			"birthplace": "New York",
			"archetype":  "Adventurer",
			"occupation": "Antiquarian",
		}
		body, _ := json.Marshal(payload)

		req := requestWithParams("POST", "/api/investigator/", body, nil)
		w := httptest.NewRecorder()

		h.CreateInvestigator(w, req)

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

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		h, _ := newTestHandler()

		req := requestWithParams("POST", "/api/investigator/", []byte("invalid json"), nil)
		w := httptest.NewRecorder()

		h.CreateInvestigator(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestUpdateInvestigator(t *testing.T) {
	t.Run("updates existing investigator", func(t *testing.T) {
		h, store := newTestHandler()

		// Create an investigator first
		inv := models.RandomInvestigator(models.Pulp)
		inv.ID = "test-id"
		store.investigators["test-id"] = inv

		payload := UpdateRequest{
			Section: "personalInfo",
			Field:   "Name",
			Value:   "Updated Name",
		}
		body, _ := json.Marshal(payload)

		req := requestWithParams("PUT", "/api/investigator/test-id", body, []string{"test-id"})
		w := httptest.NewRecorder()

		h.UpdateInvestigator(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}
	})

	t.Run("returns error for non-existent investigator", func(t *testing.T) {
		h, _ := newTestHandler()

		payload := UpdateRequest{
			Section: "personalInfo",
			Field:   "Name",
			Value:   "Updated Name",
		}
		body, _ := json.Marshal(payload)

		req := requestWithParams("PUT", "/api/investigator/nonexistent", body, []string{"nonexistent"})
		w := httptest.NewRecorder()

		h.UpdateInvestigator(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestImportInvestigatorsList(t *testing.T) {
	t.Run("imports with valid code", func(t *testing.T) {
		h, _ := newTestHandler()

		payload := map[string]string{"ImportCode": "valid-code"}
		body, _ := json.Marshal(payload)

		req := requestWithParams("POST", "/api/investigator/list/import/", body, nil)
		w := httptest.NewRecorder()

		h.ImportInvestigatorsList(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("returns error for missing import code", func(t *testing.T) {
		h, _ := newTestHandler()

		payload := map[string]string{}
		body, _ := json.Marshal(payload)

		req := requestWithParams("POST", "/api/investigator/list/import/", body, nil)
		w := httptest.NewRecorder()

		h.ImportInvestigatorsList(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}
