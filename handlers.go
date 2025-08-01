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
	"slices"
	"strconv"
	"strings"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	component := views.Home()
	err := component.Render(r.Context(), w)
	if err != nil {
		log.Println(err)
	}
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	modeQParam := r.URL.Query().Get("mode")
	mode := models.Pulp
	if modeQParam == "classic" {
		mode = models.Classic
	}

	investigator := models.RandomInvestigator(mode)
	cm := storage.NewInvestigatorCookieConfig()
	cm.SaveInvestigatorCookie(w, investigator)
	components := views.CharacterSheet(investigator)
	err := components.Render(r.Context(), w)
	if err != nil {
		log.Println(err)
	}
}

func handleExportPDF(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)
	key := strings.TrimPrefix(r.URL.Path, "/api/investigator/PDF/")
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
	cm := storage.NewInvestigatorCookieConfig()
	investigators, err := cm.ListInvestigators(r)
	if err != nil {
		log.Println(err)
	}
	components := views.InvestigatorsList(investigators)
	components.Render(r.Context(), w)

}

func handleListInvestigatorsExport(w http.ResponseWriter, r *http.Request) {
	cm := storage.NewInvestigatorCookieConfig()
	ExportCode, _ := cm.ExportInvestigatorsList(r)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(ExportCode))

}

func handleListInvestigatorsImport(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	var payload map[string]string
	err = json.Unmarshal(body, &payload)
	if err != nil {
		log.Println(err)
	}
	cm := storage.NewInvestigatorCookieConfig()
	cm.ImportInvestigatorsList(w, payload["ImportCode"])
	w.Header().Set("HX-Trigger", "import")
	w.WriteHeader(http.StatusCreated)
}

func handleCreateBaseInvestigator(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		log.Println(err)
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	payload := make(map[string]any)
	keysToConvert := []string{"age"}

	for key, val := range jsonData {
		if slices.Contains(keysToConvert, key) {
			// Handle age conversion from string to int
			if strVal, ok := val.(string); ok {
				if intVal, err := strconv.Atoi(strVal); err == nil {
					payload[key] = intVal
				} else {
					payload[key] = 0
				}
			} else {
				payload[key] = val
			}
		} else {
			payload[key] = val
		}
	}
	investigator := models.InvestigatorBaseCreate(payload)
	cm := storage.NewInvestigatorCookieConfig()
	key, err := cm.SaveInvestigatorCookie(w, investigator)
	if err != nil {
		log.Printf("Error encoding Cookie: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")

	type Response struct {
		Key string
	}

	response := Response{
		Key: key,
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

func handleDeleteInvestigator(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/investigator/")
	key = strings.Trim(key, "/")
	cm := storage.NewInvestigatorCookieConfig()

	cm.DeleteInvestigatorCookie(w, key)
	w.Header().Set("HX-Trigger", "deleted")
	w.WriteHeader(http.StatusOK)

}

func handleUpdateInvestigator(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/investigator/")
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
	case "stats":
		switch serializer.Field {
		case "TemporaryInsane":
			investigator.TemporaryInsane = !investigator.TemporaryInsane
		case "IndefiniteInsane":
			investigator.IndefiniteInsane = !investigator.IndefiniteInsane
		case "MajorWound":
			investigator.MajorWound = !investigator.MajorWound
		case "Unconscious":
			investigator.Unconscious = !investigator.Unconscious
		case "Dying":
			investigator.Dying = !investigator.Dying
		}
	case "personalInfo":
		switch serializer.Field {
		case "Name":
			investigator.Name = serializer.Value.(string)
		case "Age":
			age, _ := strconv.Atoi(serializer.Value.(string))
			investigator.Age = age
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
	case "skill_check":
		skill, ok := investigator.Skills[serializer.Field]
		if !ok {
			http.Error(w, "Skill not found", http.StatusNotFound)
		}
		skill.IsSelected = !skill.IsSelected
		investigator.Skills[serializer.Field] = skill
	case "skill_prio":
		skill, ok := investigator.Skills[serializer.Field]
		if !ok {
			http.Error(w, "Skill not found", http.StatusNotFound)
		}
		skill.IsPriority = !skill.IsPriority
		investigator.Skills[serializer.Field] = skill

	case "skill_name":
		skill, ok := investigator.Skills[serializer.Field]

		if !ok {
			http.Error(w, "Skill not found", http.StatusNotFound)
		}
		if serializer.Value.(string) == skill.Name {
			http.Error(w, "No change to Skill", http.StatusNotModified)
		}
		skill.Name = serializer.Value.(string)
		investigator.Skills[skill.Name] = skill
		delete(investigator.Skills, serializer.Field)

	default:
		http.Error(w, "Unknown section", http.StatusBadRequest)
	}
	cm.UpdateInvestigatorCookie(w, key, investigator)

}

func handleGetInvestigator(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/investigator/")
	key = strings.Trim(key, "/")
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

func handleConfirmAttrStepInvestigator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	if err := r.ParseForm(); err != nil {
		log.Println(err)
	}
	key := strings.TrimPrefix(r.URL.Path, "/api/investigator/confirm-attributes/")
	cm := storage.NewInvestigatorCookieConfig()
	investigator, err := cm.GetInvestigatorCookie(r, key)
	if err != nil {
		log.Println(err)
	}
	formToInt := func(val string) int {
		value, err := strconv.Atoi(val)
		if err != nil {
			return 0 // or some default value
		}
		return value
	}
	payload := make(map[string]int)
	keysToConvert := []string{"STR", "CON", "DEX", "INT", "POW", "APP", "EDU", "SIZ", "LCK"}
	for key, val := range r.PostForm {
		val = r.PostForm[key]
		if slices.Contains(keysToConvert, key) {
			payload[key] = formToInt(val[0])
		}
	}
	investigator.InvestigatorUpdateAttributes(payload)
	cm.UpdateInvestigatorCookie(w, key, investigator)

	components := views.SkillAssignmentForm(investigator)
	err = components.Render(r.Context(), w)
	if err != nil {
		log.Println(err)
	}
}

func handleArchetypeOccupations(w http.ResponseWriter, r *http.Request) {
	// Extract archetype name from URL path
	// Expected format: /api/archetype/{archetypeName}/occupations
	archetypeName := r.Context().Value("params").([]string)[0]
	if archetypeName == "" {
		http.Error(w, "Archetype name is required", http.StatusBadRequest)
		return
	}

	// Get the archetype from the models
	archetype, exists := models.Archetypes[archetypeName]
	if !exists {
		http.Error(w, "Archetype not found", http.StatusNotFound)
		return
	}

	// Create response structure
	type OccupationInfo struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	type Response struct {
		Suggested []OccupationInfo `json:"suggested"`
		Others    []OccupationInfo `json:"others"`
	}

	response := Response{
		Suggested: make([]OccupationInfo, 0),
		Others:    make([]OccupationInfo, 0),
	}

	// Helper function to check if occupation is suggested
	isSuggested := func(occupationName string) bool {
		for _, suggested := range archetype.SuggestedOccupations {
			if suggested == occupationName {
				return true
			}
		}
		return false
	}

	// Add suggested occupations first
	for _, suggestedOccName := range archetype.SuggestedOccupations {
		if occupation, exists := models.Occupations[suggestedOccName]; exists {
			response.Suggested = append(response.Suggested, OccupationInfo{
				Name:        occupation.Name,
				Description: occupation.GetDescription(),
			})
		}
	}

	// Add all other occupations
	for _, occupationName := range models.OccupationsList {
		if !isSuggested(occupationName) {
			if occupation, exists := models.Occupations[occupationName]; exists {
				response.Others = append(response.Others, OccupationInfo{
					Name:        occupation.Name,
					Description: occupation.GetDescription(),
				})
			}
		}
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
