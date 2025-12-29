package handlers

import (
	"net/http"
	"strconv"
)

// PaginationParams holds pagination parameters from a request
type PaginationParams struct {
	Page    int
	PerPage int
	Offset  int
}

// DefaultPerPage is the default number of items per page
const DefaultPerPage = 10

// MaxPerPage is the maximum allowed items per page
const MaxPerPage = 100

// GetPaginationParams extracts pagination parameters from request query string
func GetPaginationParams(r *http.Request) PaginationParams {
	params := PaginationParams{
		Page:    1,
		PerPage: DefaultPerPage,
	}

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil && perPage > 0 {
			if perPage > MaxPerPage {
				perPage = MaxPerPage
			}
			params.PerPage = perPage
		}
	}

	params.Offset = (params.Page - 1) * params.PerPage
	return params
}

// Paginate returns a slice of items for the given page
func Paginate[T any](items []T, params PaginationParams) []T {
	if len(items) == 0 {
		return items
	}

	start := params.Offset
	if start >= len(items) {
		return []T{}
	}

	end := start + params.PerPage
	if end > len(items) {
		end = len(items)
	}

	return items[start:end]
}

// PaginationMeta returns pagination metadata
func PaginationMeta(total int, params PaginationParams) *APIMeta {
	return &APIMeta{
		Total:   total,
		Page:    params.Page,
		PerPage: params.PerPage,
	}
}
