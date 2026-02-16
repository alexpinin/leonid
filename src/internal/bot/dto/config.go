package dto

import (
	"time"

	"github.com/openai/openai-go"
)

type Config struct {
	ID                  string
	ChatID              int64
	ChatActivatedAt     time.Time
	Pass                string
	PassValidBy         time.Time
	Nicknames           []string
	SystemPrompt        string
	MessagePrompt       string
	ConversationHistory string
}

type OpenAIConversationMessage struct {
	OfUser      *openai.ChatCompletionUserMessageParam      `json:"ofUser"`
	OfAssistant *openai.ChatCompletionAssistantMessageParam `json:"ofAssistant"`
}

type OpenAIConversationHistory struct {
	Messages []OpenAIConversationMessage `json:"messages"`
}
