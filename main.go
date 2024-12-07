package main

import (
	"book-of-shadows/models"
	"book-of-shadows/serializers"
	"book-of-shadows/views"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}
	defer r.Body.Close()
	inv, err := serializers.FromJSON(body)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(inv)
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
	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
