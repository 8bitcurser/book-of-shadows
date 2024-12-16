package main

import (
	"book-of-shadows/models"
	"book-of-shadows/serializers"
	"book-of-shadows/storage"
	"book-of-shadows/views"
	"encoding/json"
	"io"
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
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	investigatorPDF := "./static/" + data["Investigators_Name"] + ".pdf"

	err := PDFExport(
		"./static/modernSheet.pdf",
		investigatorPDF,
		data)
	if err != nil {
		http.Error(w, "Error generating PDF: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename="+data["Investigator_Name"]+".pdf")
	http.ServeFile(w, r, investigatorPDF)
	err = os.Remove(investigatorPDF)
	if err != nil {
		log.Println(err)
	}
}

func handleImportJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyIO, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}
	inv, err := serializers.FromJSON(bodyIO)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	components := views.CharacterSheet(inv)
	err = components.Render(r.Context(), w)
	if err != nil {
		log.Println(err)
	}

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
	// Log existing cookies
	log.Println("Current cookies:")
	for _, cookie := range r.Cookies() {
		log.Printf("- %s: %s", cookie.Name, cookie.Value)
	}
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
	http.HandleFunc("/api/export-pdf", handleExportPDF)
	http.HandleFunc("/api/import-json", handleImportJSON)
	http.HandleFunc("/api/list-investigators", handleListInvestigators)
	http.HandleFunc("/api/delete-investigator/", handleDeleteInvestigator)
	http.HandleFunc("/api/get-investigator/", handleGetInvestigator)
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
