package handler

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type chatChecker struct {
	basicHandler
	chChecker
}

func newChatChecker(ch chChecker) *chatChecker {
	return &chatChecker{
		chChecker: ch,
	}
}

type chChecker interface {
	IsChatActive(ctx context.Context, chatID int64) (bool, error)
}

func (h *chatChecker) handle(ctx context.Context, b *bot.Bot, u *models.Update, uc *UpdateContext) error {
	isChatActive, err := h.IsChatActive(ctx, u.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("chatChecker.handle: %w", err)
	}
	uc.IsChatActive = isChatActive

	return h.nextHandle(ctx, b, u, uc)
}
