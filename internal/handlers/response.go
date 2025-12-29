package handlers

import (
	"encoding/json"
	"net/http"
)

// APIResponse is the standard response format for all API endpoints
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *APIMeta    `json:"meta,omitempty"`
}

// APIError represents a standardized error response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// APIMeta contains optional metadata about the response
type APIMeta struct {
	Total   int `json:"total,omitempty"`
	Page    int `json:"page,omitempty"`
	PerPage int `json:"per_page,omitempty"`
}

// respondSuccess sends a successful JSON response
func (h *Handler) respondSuccess(w http.ResponseWriter, status int, data interface{}, meta *APIMeta) {
	response := APIResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
	h.sendJSONResponse(w, status, response)
}

// respondAPIError sends a standardized error response
func (h *Handler) respondAPIError(w http.ResponseWriter, status int, code, message string) {
	response := APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	}
	h.sendJSONResponse(w, status, response)
}

// sendJSONResponse is a helper to send JSON responses
func (h *Handler) sendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Printf("Failed to encode response: %v", err)
	}
}

// Error codes for consistent API responses
const (
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeBadRequest     = "BAD_REQUEST"
	ErrCodeInvalidJSON    = "INVALID_JSON"
	ErrCodeMissingField   = "MISSING_FIELD"
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeConflict       = "CONFLICT"
	ErrCodePayloadTooLarge = "PAYLOAD_TOO_LARGE"
	ErrCodeInternalError  = "INTERNAL_ERROR"
)
