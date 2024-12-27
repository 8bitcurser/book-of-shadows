package main

import (
	"book-of-shadows/models"
	"book-of-shadows/storage"
	"book-of-shadows/views"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	component := views.Home()
	err := component.Render(r.Context(), w)
	if err != nil {
		log.Println(err)
	}
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	modeQParam := r.URL.Query().Get("mode")
	mode := models.Pulp
	if modeQParam == "classic" {
		mode = models.Classic
	}

	investigator := models.NewInvestigator(mode)
	cm := storage.NewInvestigatorCookieConfig()
	cm.SaveInvestigatorCookie(w, investigator)
	components := views.CharacterSheet(investigator)
	err := components.Render(r.Context(), w)
	if err != nil {
		log.Println(err)
	}
}

func handleExportPDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	data := make(map[string]string)
	key := strings.TrimPrefix(r.URL.Path, "/api/export-pdf/")
	if key == "" {
		http.Error(w, "No investigator Key passed", http.StatusBadRequest)
	}
	cm := storage.NewInvestigatorCookieConfig()
	investigator, err := cm.GetInvestigatorCookie(r, key)

	if err != nil {
		log.Println(err)
	}
	data = ConvertInvestigatorToMap(investigator)
	fileName := fmt.Sprintf("%s.pdf", strings.ReplaceAll(data["Investigators_Name"], " ", "_"))
	investigatorPDF := "./static/" + fileName

	err = PDFExport(
		"./static/modernSheet.pdf",
		investigatorPDF,
		data)

	if err != nil {
		http.Error(w, "Error generating PDF: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)

	http.ServeFile(w, r, investigatorPDF)

	defer os.Remove(investigatorPDF)

}

func handleListInvestigators(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cm := storage.NewInvestigatorCookieConfig()
	investigators, err := cm.ListInvestigators(r)
	if err != nil {
		log.Println(err)
	}
	components := views.InvestigatorsList(investigators)
	components.Render(r.Context(), w)

}

func handleDeleteInvestigator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	key := strings.TrimPrefix(r.URL.Path, "/api/delete-investigator/")
	cm := storage.NewInvestigatorCookieConfig()

	cm.DeleteInvestigatorCookie(w, key)

	w.WriteHeader(http.StatusOK)

}

func handleGetInvestigator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	key := strings.TrimPrefix(r.URL.Path, "/api/get-investigator/")
	cm := storage.NewInvestigatorCookieConfig()
	investigator, err := cm.GetInvestigatorCookie(r, key)

	if err != nil {
		log.Println(err)
	}
	components := views.CharacterSheet(investigator)
	err = components.Render(r.Context(), w)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	fileServer := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	// routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/generate", handleGenerate)
	http.HandleFunc("/api/export-pdf/", handleExportPDF)
	http.HandleFunc("/api/list-investigators", handleListInvestigators)
	http.HandleFunc("/api/delete-investigator/", handleDeleteInvestigator)
	http.HandleFunc("/api/get-investigator/", handleGetInvestigator)
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
