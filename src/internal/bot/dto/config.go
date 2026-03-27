package dto

import (
	"time"

	"github.com/openai/openai-go"
)

type Config struct {
	ID                  int
	ChatID              int64
	ChatActivatedAt     time.Time
	Pass                string
	PassValidBy         time.Time
	Nicknames           []string
	SystemPrompt        string
	ConversationHistory string // might vary for different LLM providers
}

type OpenAIConversationMessage struct {
	OfUser      *openai.ChatCompletionUserMessageParam      `json:"ofUser"`
	OfAssistant *openai.ChatCompletionAssistantMessageParam `json:"ofAssistant"`
}

type OpenAIConversationHistory struct {
	Messages []OpenAIConversationMessage `json:"messages"`
}
