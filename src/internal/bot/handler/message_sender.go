package handler

import (
	"context"

	"github.com/go-telegram/bot"
)

type messageSender struct {
	basicHandler
	mSender
}

func newMessageSender(ms mSender) *messageSender {
	return &messageSender{
		mSender: ms,
	}
}

type mSender interface {
	SendMessage(ctx context.Context, b *bot.Bot, chatID int64, message string)
}

func (h *messageSender) handle(ctx context.Context, b *bot.Bot, u *UpdateContext) {
	h.SendMessage(ctx, b, u.Message.Chat.ID, u.Message.Text)
	h.nextHandle(ctx, b, u)
}
