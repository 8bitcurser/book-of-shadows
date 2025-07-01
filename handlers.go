package main

import (
	"book-of-shadows/models"
	"book-of-shadows/serializers"
	"book-of-shadows/storage"
	"book-of-shadows/views"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
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
	if err := r.ParseForm(); err != nil {
		log.Println(err)
	}
	formToInt := func(val string) int {
		value, err := strconv.Atoi(val)
		if err != nil {
			return 0 // or some default value
		}
		return value
	}
	payload := make(map[string]any)
	keysToConvert := []string{"age"}
	for key, val := range r.PostForm {
		val = r.PostForm[key]
		if slices.Contains(keysToConvert, key) {
			payload[key] = formToInt(val[0])
		} else {
			payload[key] = val[0]
		}
	}
	investigator := models.InvestigatorBaseCreate(payload)
	cm := storage.NewInvestigatorCookieConfig()
	cm.SaveInvestigatorCookie(w, investigator)
	components := views.AssignAttrForm(investigator)
	err := components.Render(r.Context(), w)
	if err != nil {
		log.Println(err)
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

func handleCreateStepInvestigator(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/api/generate-step/")
	cm := storage.NewInvestigatorCookieConfig()
	investigator := &models.Investigator{}
	err := error(nil)
	if key != "/api/generate-step" {
		investigator, err = cm.GetInvestigatorCookie(r, key)
		if err != nil {
			http.Error(w, "Investigator cookie missing", http.StatusNotFound)
		}
	}
	component := views.BaseInvForm(investigator)
	err = component.Render(r.Context(), w)
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

func handleReportIssue(w http.ResponseWriter, r *http.Request) {
	// Get Telegram credentials from environment variables
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")

	if telegramToken == "" || telegramChatID == "" {
		log.Println("Telegram credentials not configured")
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}

	// Parse the request body
	var report models.IssueReport
	err := json.NewDecoder(r.Body).Decode(&report)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if report.IssueType == "" || report.Description == "" {
		http.Error(w, "Issue type and description are required", http.StatusBadRequest)
		return
	}

	// Format message for Telegram
	parsedTime, _ := time.Parse(time.RFC3339, report.Timestamp)
	localTime := parsedTime.Format("2006-01-02 15:04:05")

	messageText := fmt.Sprintf(
		"🐞 *CorbittFiles Issue Report*\n\n"+
			"📋 *Type*: %s\n\n"+
			"💬 *Description*: %s\n\n"+
			"📧 *User Email*: %s\n\n"+
			"⏰ *Date*: %s",
		report.IssueType,
		report.Description,
		report.Email,
		localTime,
	)

	// Create Telegram message payload
	telegramMsg := models.TelegramMessage{
		ChatID:    telegramChatID,
		Text:      messageText,
		ParseMode: "Markdown",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(telegramMsg)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Send to Telegram API
	telegramURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegramToken)
	resp, err := http.Post(telegramURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending to Telegram: %v", err)
		http.Error(w, "Failed to send to Telegram", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check response from Telegram
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Telegram API error. Status: %d, Response: %s", resp.StatusCode, string(body))
		http.Error(w, "Failed to send to Telegram", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func handleArchetypeOccupations(w http.ResponseWriter, r *http.Request) {
	// Extract archetype name from URL path
	// Expected format: /api/archetype/{archetypeName}/occupations
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/archetype/"), "/")
	if len(pathParts) < 2 || pathParts[1] != "occupations" {
		http.Error(w, "Invalid URL format. Expected: /api/archetype/{name}/occupations", http.StatusBadRequest)
		return
	}

	archetypeName := pathParts[0]
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
