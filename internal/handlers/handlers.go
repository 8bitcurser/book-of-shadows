package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"book-of-shadows/internal/errors"
	"book-of-shadows/models"
	"book-of-shadows/storage"
	"book-of-shadows/views"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	store  storage.Store
	logger *log.Logger
}

// New creates a new Handler with dependencies
func New(store storage.Store, logger *log.Logger) *Handler {
	if logger == nil {
		logger = log.New(io.Discard, "", 0)
	}
	return &Handler{
		store:  store,
		logger: logger,
	}
}

// respondJSON sends a JSON response
func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Printf("Failed to encode response: %v", err)
	}
}

// respondError sends an error response
func (h *Handler) respondError(w http.ResponseWriter, err error) {
	var httpErr errors.HTTPError
	if e, ok := err.(errors.HTTPError); ok {
		httpErr = e
	} else {
		// Map known errors to HTTP status codes
		switch err {
		case errors.ErrNotFound, errors.ErrCookieNotFound:
			httpErr = errors.NewHTTPError(http.StatusNotFound, "Resource not found", err)
		case errors.ErrInvalidData, errors.ErrInvalidAttribute, errors.ErrInvalidSkill:
			httpErr = errors.NewHTTPError(http.StatusBadRequest, "Invalid request", err)
		case errors.ErrAlreadyExists:
			httpErr = errors.NewHTTPError(http.StatusConflict, "Resource already exists", err)
		case errors.ErrCookieTooLarge:
			httpErr = errors.NewHTTPError(http.StatusRequestEntityTooLarge, "Data too large", err)
		default:
			httpErr = errors.NewHTTPError(http.StatusInternalServerError, "Internal server error", err)
		}
	}

	h.logger.Printf("HTTP %d: %v", httpErr.Code, httpErr)

	response := map[string]interface{}{
		"error": httpErr.Message,
	}

	h.respondJSON(w, httpErr.Code, response)
}

// Home handles the home page
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	component := views.Home()
	if err := component.Render(r.Context(), w); err != nil {
		h.logger.Printf("Failed to render home: %v", err)
		h.respondError(w, err)
	}
}

// Generate creates a random investigator
func (h *Handler) Generate(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	modeParam := r.URL.Query().Get("mode")
	mode := models.Pulp
	if modeParam == "classic" {
		mode = models.Classic
	}

	// Generate investigator
	investigator := models.RandomInvestigator(mode)

	// Save to cookie
	_, err := h.store.SaveInvestigator(w, investigator)
	if err != nil {
		h.respondError(w, err)
		return
	}

	// Render response
	component := views.CharacterSheet(investigator)
	if err := component.Render(r.Context(), w); err != nil {
		h.logger.Printf("Failed to render character sheet: %v", err)
		h.respondError(w, err)
	}
}

// ListInvestigators returns all investigators
func (h *Handler) ListInvestigators(w http.ResponseWriter, r *http.Request) {
	investigators, err := h.store.ListInvestigators(r)
	if err != nil {
		h.respondError(w, err)
		return
	}

	component := views.InvestigatorsList(investigators)
	if err := component.Render(r.Context(), w); err != nil {
		h.logger.Printf("Failed to render investigators list: %v", err)
		h.respondError(w, err)
	}
}

// GetInvestigator retrieves a specific investigator
func (h *Handler) GetInvestigator(w http.ResponseWriter, r *http.Request) {
	// Extract ID from context (set by router)
	params := r.Context().Value("params").([]string)
	if len(params) == 0 {
		h.respondError(w, errors.NewHTTPError(http.StatusBadRequest, "Missing investigator ID", nil))
		return
	}
	id := params[0]

	investigator, err := h.store.GetInvestigator(r, id)
	if err != nil {
		h.respondError(w, err)
		return
	}

	component := views.CharacterSheet(investigator)
	if err := component.Render(r.Context(), w); err != nil {
		h.logger.Printf("Failed to render character sheet: %v", err)
		h.respondError(w, err)
	}
}

// CreateInvestigator creates a new investigator
func (h *Handler) CreateInvestigator(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.respondError(w, errors.NewHTTPError(http.StatusBadRequest, "Failed to read request body", err))
		return
	}
	defer r.Body.Close()

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		h.respondError(w, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON", err))
		return
	}

	// Convert and validate payload
	processedPayload := h.processInvestigatorPayload(payload)

	// Create investigator
	investigator := models.InvestigatorBaseCreate(processedPayload)

	// Save to cookie
	key, err := h.store.SaveInvestigator(w, investigator)
	if err != nil {
		h.respondError(w, err)
		return
	}

	// Send response
	h.respondJSON(w, http.StatusCreated, map[string]string{
		"Key": key,
	})
}

// UpdateInvestigator updates an existing investigator
func (h *Handler) UpdateInvestigator(w http.ResponseWriter, r *http.Request) {
	// Extract ID from context
	params := r.Context().Value("params").([]string)
	if len(params) == 0 {
		h.respondError(w, errors.NewHTTPError(http.StatusBadRequest, "Missing investigator ID", nil))
		return
	}
	id := params[0]

	// Get existing investigator
	investigator, err := h.store.GetInvestigator(r, id)
	if err != nil {
		h.respondError(w, err)
		return
	}

	// Parse update request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.respondError(w, errors.NewHTTPError(http.StatusBadRequest, "Failed to read request body", err))
		return
	}
	defer r.Body.Close()

	var updateReq UpdateRequest
	if err := json.Unmarshal(body, &updateReq); err != nil {
		h.respondError(w, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON", err))
		return
	}

	// Apply updates
	if err := h.applyInvestigatorUpdate(investigator, &updateReq); err != nil {
		h.respondError(w, err)
		return
	}

	// Save updated investigator
	if err := h.store.UpdateInvestigator(w, id, investigator); err != nil {
		h.respondError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteInvestigator deletes an investigator
func (h *Handler) DeleteInvestigator(w http.ResponseWriter, r *http.Request) {
	// Extract ID from context
	params := r.Context().Value("params").([]string)
	if len(params) == 0 {
		h.respondError(w, errors.NewHTTPError(http.StatusBadRequest, "Missing investigator ID", nil))
		return
	}
	id := params[0]

	if err := h.store.DeleteInvestigator(w, id); err != nil {
		h.respondError(w, err)
		return
	}

	w.Header().Set("HX-Trigger", "deleted")
	w.WriteHeader(http.StatusOK)
}

// ExportInvestigatorsList exports all investigators for sharing
func (h *Handler) ExportInvestigatorsList(w http.ResponseWriter, r *http.Request) {
	exportCode, err := h.store.ExportInvestigatorsList(r)
	if err != nil {
		h.respondError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, exportCode)
}

// ImportInvestigatorsList imports investigators from a share code
func (h *Handler) ImportInvestigatorsList(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.respondError(w, errors.NewHTTPError(http.StatusBadRequest, "Failed to read request body", err))
		return
	}
	defer r.Body.Close()

	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		h.respondError(w, errors.NewHTTPError(http.StatusBadRequest, "Invalid JSON", err))
		return
	}

	importCode, ok := payload["ImportCode"]
	if !ok || importCode == "" {
		h.respondError(w, errors.NewHTTPError(http.StatusBadRequest, "Missing import code", nil))
		return
	}

	if err := h.store.ImportInvestigatorsList(w, importCode); err != nil {
		h.respondError(w, err)
		return
	}

	w.Header().Set("HX-Trigger", "import")
	w.WriteHeader(http.StatusCreated)
}