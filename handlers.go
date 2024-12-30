package main

import (
	"book-of-shadows/models"
	"book-of-shadows/serializers"
	"book-of-shadows/storage"
	"book-of-shadows/views"
	"encoding/json"
	"fmt"
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
	key := strings.TrimPrefix(r.URL.Path, "/api/investigator/export/")
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
	key := strings.TrimPrefix(r.URL.Path, "/api/investigator/delete/")
	cm := storage.NewInvestigatorCookieConfig()

	cm.DeleteInvestigatorCookie(w, key)
	w.Header().Set("HX-Trigger", "deleted")
	w.WriteHeader(http.StatusOK)

}

func handleUpdateInvestigator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	key := strings.TrimPrefix(r.URL.Path, "/api/investigator/update/")
	cm := storage.NewInvestigatorCookieConfig()
	investigator, err := cm.GetInvestigatorCookie(r, key)

	if err != nil || investigator == nil {
		http.Error(w, "Investigator cookie missing", http.StatusNotFound)
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	var serializer serializers.UpdateRequestSerializer
	if err := json.Unmarshal(body, &serializer); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	switch serializer.Section {
	case "skills":
		skill, ok := investigator.Skills[serializer.Field]
		if !ok {
			http.Error(w, "Skill not found", http.StatusNotFound)
		}
		skill.Value = int(serializer.Value.(float64))
		investigator.Skills[serializer.Field] = skill
	case "personalInfo":
		switch serializer.Field {
		case "Name":
			investigator.Name = serializer.Value.(string)
		case "Age":
			investigator.Age = serializer.Value.(int)
		case "Residence":
			investigator.Residence = serializer.Value.(string)
		case "Birthplace":
			investigator.Birthplace = serializer.Value.(string)
		default:
			http.Error(w, "Unsupported field", http.StatusNotFound)
		}
	case "combat":
		attr, ok := investigator.Attributes[serializer.Field]
		if !ok {
			http.Error(w, "Attribute not found", http.StatusNotFound)
		}
		attr.Value = int(serializer.Value.(float64))
		investigator.Attributes[serializer.Field] = attr
	default:
		http.Error(w, "Unknown section", http.StatusBadRequest)
	}
	cm.UpdateInvestigatorCookie(w, key, investigator)

}

func handleGetInvestigator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	key := strings.TrimPrefix(r.URL.Path, "/api/investigator/")
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
