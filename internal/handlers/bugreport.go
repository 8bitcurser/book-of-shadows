package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"book-of-shadows/internal/errors"
)

// IssueReport represents a bug report from a user
type IssueReport struct {
	IssueType   string `json:"issueType"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Timestamp   string `json:"timestamp"`
}

// TelegramMessage represents a message to send to Telegram
type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// ReportIssue handles bug reports and sends them to Telegram
func (h *Handler) ReportIssue(w http.ResponseWriter, r *http.Request) {
	// Get Telegram credentials from environment variables
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")

	if telegramToken == "" || telegramChatID == "" {
		h.logger.Println("Telegram credentials not configured")
		h.respondError(w, errors.NewHTTPError(http.StatusInternalServerError, "Server configuration error", nil))
		return
	}

	// Parse the request body
	var report IssueReport
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		h.logger.Printf("Error parsing request body: %v", err)
		h.respondError(w, errors.NewHTTPError(http.StatusBadRequest, "Invalid request body", err))
		return
	}
	defer r.Body.Close()

	// Validate required fields
	if report.IssueType == "" {
		h.respondError(w, errors.NewValidationError("issueType", "Issue type is required"))
		return
	}
	if report.Description == "" {
		h.respondError(w, errors.NewValidationError("description", "Description is required"))
		return
	}

	// Format message for Telegram
	localTime := time.Now().Format("2006-01-02 15:04:05")
	if report.Timestamp != "" {
		if parsedTime, err := time.Parse(time.RFC3339, report.Timestamp); err == nil {
			localTime = parsedTime.Format("2006-01-02 15:04:05")
		}
	}

	email := report.Email
	if email == "" {
		email = "Not provided"
	}

	messageText := fmt.Sprintf(
		"🐞 *Book of Shadows Issue Report*\n\n"+
			"📋 *Type*: %s\n\n"+
			"💬 *Description*: %s\n\n"+
			"📧 *User Email*: %s\n\n"+
			"⏰ *Date*: %s",
		escapeMarkdown(report.IssueType),
		escapeMarkdown(report.Description),
		escapeMarkdown(email),
		localTime,
	)

	// Send to Telegram
	if err := h.sendToTelegram(telegramToken, telegramChatID, messageText); err != nil {
		h.logger.Printf("Failed to send to Telegram: %v", err)
		h.respondError(w, errors.NewHTTPError(http.StatusInternalServerError, "Failed to send report", err))
		return
	}

	// Return success response
	h.respondJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// sendToTelegram sends a message to a Telegram chat
func (h *Handler) sendToTelegram(token, chatID, message string) error {
	telegramMsg := TelegramMessage{
		ChatID:    chatID,
		Text:      message,
		ParseMode: "Markdown",
	}

	jsonData, err := json.Marshal(telegramMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	telegramURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	resp, err := http.Post(telegramURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram API error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// escapeMarkdown escapes special characters for Telegram Markdown
func escapeMarkdown(text string) string {
	replacer := map[string]string{
		"_": "\\_",
		"*": "\\*",
		"[": "\\[",
		"`": "\\`",
	}
	result := text
	for old, new := range replacer {
		result = replaceAll(result, old, new)
	}
	return result
}

// replaceAll is a simple string replace function
func replaceAll(s, old, new string) string {
	var result []byte
	for i := 0; i < len(s); i++ {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result = append(result, new...)
			i += len(old) - 1
		} else {
			result = append(result, s[i])
		}
	}
	return string(result)
}