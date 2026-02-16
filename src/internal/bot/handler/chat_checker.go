package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-telegram/bot"
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

func (h *chatChecker) handle(ctx context.Context, b *bot.Bot, u *UpdateContext) error {
	isChatActive, err := h.IsChatActive(ctx, u.Message.Chat.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("chatChecker.handle: %w", err)
	}
	u.IsChatActive = isChatActive

	return h.nextHandle(ctx, b, u)
}
