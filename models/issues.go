package models

// IssueReport represents the data structure for bug reports
type IssueReport struct {
	IssueType   string `json:"issueType"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Timestamp   string `json:"timestamp"`
}

type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}
