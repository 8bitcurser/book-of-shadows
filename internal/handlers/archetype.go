package handlers

import (
	"net/http"

	"book-of-shadows/internal/errors"
	"book-of-shadows/models"
)

// OccupationInfo represents occupation data for API responses
type OccupationInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ArchetypeOccupationsResponse is the response structure for archetype occupations
type ArchetypeOccupationsResponse struct {
	Suggested []OccupationInfo `json:"suggested"`
	Others    []OccupationInfo `json:"others"`
}

// GetArchetypeOccupations returns occupations for a given archetype
func (h *Handler) GetArchetypeOccupations(w http.ResponseWriter, r *http.Request) {
	// Extract archetype name from URL path
	params := r.Context().Value("params").([]string)
	if len(params) == 0 {
		h.respondError(w, errors.NewHTTPError(http.StatusBadRequest, "Archetype name is required", nil))
		return
	}
	archetypeName := params[0]

	if archetypeName == "" {
		h.respondError(w, errors.NewHTTPError(http.StatusBadRequest, "Archetype name is required", nil))
		return
	}

	// Get the archetype from the models
	archetype, exists := models.Archetypes[archetypeName]
	if !exists {
		h.respondError(w, errors.NewHTTPError(http.StatusNotFound, "Archetype not found", nil))
		return
	}

	// Build response
	response := ArchetypeOccupationsResponse{
		Suggested: make([]OccupationInfo, 0, len(archetype.SuggestedOccupations)),
		Others:    make([]OccupationInfo, 0),
	}

	// Create a set of suggested occupations for O(1) lookup
	suggestedSet := make(map[string]struct{}, len(archetype.SuggestedOccupations))
	for _, name := range archetype.SuggestedOccupations {
		suggestedSet[name] = struct{}{}
	}

	// Add suggested occupations
	for _, suggestedOccName := range archetype.SuggestedOccupations {
		if occupation, exists := models.Occupations[suggestedOccName]; exists {
			response.Suggested = append(response.Suggested, OccupationInfo{
				Name:        occupation.Name,
				Description: occupation.GetDescription(),
			})
		}
	}

	// Add all other occupations
	for _, occupationName := range models.OccupationsList {
		if _, isSuggested := suggestedSet[occupationName]; !isSuggested {
			if occupation, exists := models.Occupations[occupationName]; exists {
				response.Others = append(response.Others, OccupationInfo{
					Name:        occupation.Name,
					Description: occupation.GetDescription(),
				})
			}
		}
	}

	h.respondJSON(w, http.StatusOK, response)
}