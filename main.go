package main

import (
	"book-of-shadows/storage"
	"log"
	"net/http"
)

func setupRoutes() *RadixTree {
	router := NewRouter()

	fs := http.FileServer(http.Dir("static"))
	fsHandler := http.StripPrefix("/static/", fs)
	router.HandleStatic("/static/", fsHandler)

	// Define routes
	router.GET("/", handleHome)
	router.GET("api/generate/", handleGenerate)

	// Investigator CRUD operations
	router.GET("api/investigator", handleListInvestigators)
	router.GET("api/investigator/{:id}", handleGetInvestigator)
	router.POST("api/investigator/", handleCreateBaseInvestigator)
	router.PUT("api/investigator/", handleUpdateInvestigator)
	router.DELETE("api/investigator/", handleDeleteInvestigator)

	router.GET("api/investigator/PDF/", handleExportPDF)
	router.GET("api/investigator/list/export", handleListInvestigatorsExport)

	router.POST("api/investigator/list/import/", handleListInvestigatorsImport)
	router.GET("api/generate-step/", handleCreateStepInvestigator)
	router.POST("api/investigator/confirm-attributes/", handleConfirmAttrStepInvestigator)

	router.GET("api/archetype/", handleArchetypeOccupations)
	router.POST("api/report-issue", handleReportIssue)

	return router
}

func main() {
	conn := storage.SQLiteDB{}
	conn.Init()
	conn.StartCleanupRoutine()
	defer conn.DB.Close()

	router := setupRoutes()
	router.PrintRoutes()

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
