package main

import (
	"book-of-shadows/bugreporting"
	"book-of-shadows/storage"
	"book-of-shadows/wizard"
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
	router.POST("api/investigator/", handleCreateBaseInvestigator)
	router.GET("api/investigator/{:id}", handleGetInvestigator)
	router.PUT("api/investigator/{:id}", handleUpdateInvestigator)
	router.DELETE("api/investigator/{:id}", handleDeleteInvestigator)

	router.GET("api/investigator/PDF/", handleExportPDF)
	router.GET("api/investigator/list/export", handleListInvestigatorsExport)

	router.POST("api/investigator/list/import/", handleListInvestigatorsImport)

	router.GET("api/archetype/{:name}/occupations/", handleArchetypeOccupations)
	router.POST("api/report-issue", bugreporting.HandleReportIssue)

	router.GET("wizard/base/{:key}", wizard.HandleBaseStep)
	router.GET("wizard/attributes/{:key}", wizard.HandleAttrStep)
	router.GET("wizard/skills/{:key}", wizard.HandleSkillForm)

	return router
}

func main() {
	conn := storage.SQLiteDB{}
	conn.Init()
	conn.StartCleanupRoutine()
	defer conn.DB.Close()

	router := setupRoutes()

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
