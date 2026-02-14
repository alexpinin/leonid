package dto

import "time"

type Config struct {
	ID                  string
	ChatID              int64
	ChatActivatedAt     time.Time
	Pass                string
	PassValidBy         time.Time
	Nicknames           []string
	SystemPrompt        string
	MessagePrompt       string
	ConversationContext string
}
