package main

import (
	"book-of-shadows/storage"
	"log"
	"net/http"
)

func main() {
	conn := storage.SQLiteDB{}
	conn.Init()
	conn.StartCleanupRoutine()
	defer conn.DB.Close()

	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	// routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/generate-random", handleGenerate)
	http.HandleFunc("/api/investigator/export/", handleExportPDF)
	http.HandleFunc("/api/investigator/list", handleListInvestigators)
	http.HandleFunc("/api/investigator/list/export", handleListInvestigatorsExport)
	http.HandleFunc("/api/investigator/list/import/", handleListInvestigatorsImport)
	http.HandleFunc("/api/generate-step", handleCreateStepInvestigator)
	http.HandleFunc("/api/investigator/confirm-attributes/", handleConfirmAttrStepInvestigator)
	http.HandleFunc("/api/investigator/confirm-archetype/", handleConfirmArchSkillStepInvestigator)
	http.HandleFunc("/api/investigator/confirm-occupation/", handleConfirmOccSkillStepInvestigator)
	http.HandleFunc("/api/investigator/create/", handleCreateBaseInvestigator)
	http.HandleFunc("/api/investigator/delete/", handleDeleteInvestigator)
	http.HandleFunc("/api/investigator/update/", handleUpdateInvestigator)
	http.HandleFunc("/api/investigator/", handleGetInvestigator)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
