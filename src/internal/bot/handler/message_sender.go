package handler

import (
	"context"

	"github.com/go-telegram/bot"

	"leonid/src/internal/logger"
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
	SendMessage(ctx context.Context, b *bot.Bot, chatID int64, message string) error
}

func (h *messageSender) handle(ctx context.Context, b *bot.Bot, u *UpdateContext) {
	err := h.SendMessage(ctx, b, u.Message.Chat.ID, u.Message.Text)
	if err != nil {
		logger.Error(err.Error())
	}
	h.nextHandle(ctx, b, u)
}
