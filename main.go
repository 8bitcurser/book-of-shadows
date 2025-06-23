package main

import (
	"book-of-shadows/storage"
	"fmt"
	"log"
	"net/http"
)

func setupRoutes() *Router {
	router := NewRouter()

	fs := http.FileServer(http.Dir("static"))
	fsHandler := http.StripPrefix("/static/", fs)
	router.HandleStatic("/static/", fsHandler)
	fmt.Printf("Static handlers: %+v\n", router.statics)

	// Define routes
	router.GET("/", handleHome)
	router.GET("/api/generate", handleGenerate)

	// Investigator CRUD operations
	router.GET("/api/investigator", handleListInvestigators)
	router.GET("/api/investigator/", handleGetInvestigator)
	router.POST("/api/investigator/", handleCreateBaseInvestigator)
	router.PUT("/api/investigator/", handleUpdateInvestigator)
	router.DELETE("/api/investigator/", handleDeleteInvestigator)

	router.GET("/api/investigator/PDF/", handleExportPDF)
	router.GET("/api/investigator/list/export", handleListInvestigatorsExport)

	router.POST("/api/investigator/list/import/", handleListInvestigatorsImport)
	router.POST("/api/generate-step/", handleCreateStepInvestigator)
	router.POST("/api/investigator/confirm-attributes/", handleConfirmAttrStepInvestigator)

	router.GET("/api/archetype/", handleArchetypeOccupations)
	router.POST("/api/report-issue", handleReportIssue)

	return router
}

func main() {
	conn := storage.SQLiteDB{}
	conn.Init()
	conn.StartCleanupRoutine()
	defer conn.DB.Close()

	router := setupRoutes()
	// routes
	// http.HandleFunc("/", handleHome)
	// http.HandleFunc("/api/generate", handleGenerate)
	// http.HandleFunc("/api/investigator/export/", handleExportPDF)
	// http.HandleFunc("/api/investigator/list", handleListInvestigators)
	// http.HandleFunc("/api/investigator/list/export", handleListInvestigatorsExport)
	// http.HandleFunc("/api/investigator/list/import/", handleListInvestigatorsImport)
	// http.HandleFunc("/api/generate-step/", handleCreateStepInvestigator)
	// http.HandleFunc("/api/investigator/confirm-attributes/", handleConfirmAttrStepInvestigator)
	// http.HandleFunc("/api/investigator/create/", handleCreateBaseInvestigator)
	// http.HandleFunc("/api/investigator/delete/", handleDeleteInvestigator)
	// http.HandleFunc("/api/investigator/update/", handleUpdateInvestigator)
	// http.HandleFunc("/api/investigator/", handleGetInvestigator)
	// http.HandleFunc("/api/archetype/", handleArchetypeOccupations)
	// http.HandleFunc("/api/report-issue", handleReportIssue)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
