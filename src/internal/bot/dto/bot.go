package dto

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// TelegramBot abstraction allows openai.go OpenAIService.SendMessage to implement message_sender.go mSender interface
// The mSender is dictated by the go-telegram/bot library where bot is injected during the runtime
type TelegramBot interface {
	SendMessage(context.Context, *bot.SendMessageParams) (*models.Message, error)
}
