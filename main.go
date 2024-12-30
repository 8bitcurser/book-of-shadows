package main

import (
	"log"
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	// routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/generate", handleGenerate)
	http.HandleFunc("/api/investigator/export/", handleExportPDF)
	http.HandleFunc("/api/investigators/list", handleListInvestigators)
	http.HandleFunc("/api/investigator/delete/", handleDeleteInvestigator)
	http.HandleFunc("/api/investigator/update/", handleUpdateInvestigator)
	http.HandleFunc("/api/investigator/", handleGetInvestigator)
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
