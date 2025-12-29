package storage

import (
	"book-of-shadows/models"
	"net/http"
)

// Store defines the interface for data storage operations
type Store interface {
	ExportStore
	InvestigatorStore
}

// ExportStore handles export/import operations
type ExportStore interface {
	SaveExport(data string) (string, error)
	GetExport(id string) (string, error)
	DeleteExpiredExports() error
}

// InvestigatorStore handles investigator CRUD operations
type InvestigatorStore interface {
	SaveInvestigator(w http.ResponseWriter, inv *models.Investigator) (string, error)
	GetInvestigator(r *http.Request, id string) (*models.Investigator, error)
	UpdateInvestigator(w http.ResponseWriter, id string, inv *models.Investigator) error
	DeleteInvestigator(w http.ResponseWriter, id string) error
	ListInvestigators(r *http.Request) (map[string]*models.Investigator, error)
	ExportInvestigatorsList(r *http.Request) (string, error)
	ImportInvestigatorsList(w http.ResponseWriter, uuid string) error
}