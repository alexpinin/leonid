package handler

import (
	"context"

	"github.com/go-telegram/bot"
)

type MessageSender struct {
	basicHandler
	messageSender messageSender
}

func NewMessageSender(ms messageSender) *MessageSender {
	return &MessageSender{
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
