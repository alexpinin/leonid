package handler

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"

	"leonid/src/internal/bot/dto"
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
	SendMessage(ctx context.Context, b dto.TelegramBot, chatID int64, message string) error
}

func (h *messageSender) handle(ctx context.Context, b *bot.Bot, u *UpdateContext) error {
	err := h.SendMessage(ctx, b, u.Message.Chat.ID, u.Message.Text)
	if err != nil {
		return fmt.Errorf("messageSender.handle: %w", err)
	}
	return h.nextHandle(ctx, b, u)
}
