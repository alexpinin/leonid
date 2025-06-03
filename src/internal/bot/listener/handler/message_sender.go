package handler

import (
	"context"
	"os"

	"github.com/go-telegram/bot"
	"github.com/revrost/go-openrouter"
)

type MessageSender struct {
	basicHandler
	llmClient *openrouter.Client
	messageSender
}

func NewMessageSender(ms messageSender) *MessageSender {
	return &MessageSender{
		llmClient: openrouter.NewClient(
			os.Getenv("LLM_TOKEN"),
		),
		messageSender: ms,
	}
}

type messageSender interface {
	SendMessage(ctx context.Context, b *bot.Bot, chatID int64, message string)
}

func (h *MessageSender) Handle(ctx context.Context, b *bot.Bot, u *UpdateContext) {
	h.messageSender.SendMessage(ctx, b, u.Message.Chat.ID, u.Message.Text)
	h.nextHandle(ctx, b, u)
}
