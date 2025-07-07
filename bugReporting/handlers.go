package bugreporting

import (
	"book-of-shadows/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func HandleReportIssue(w http.ResponseWriter, r *http.Request) {
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
