package handler

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"

	"leonid/src/internal/logger"
)

type chatChecker struct {
	basicHandler
	chChecker
}

func NewChatChecker(ch chChecker) *chatChecker {
	return &chatChecker{
		chChecker: ch,
	}
}

type chChecker interface {
	IsChatActive(ctx context.Context, chatID int64) (bool, error)
}

func (h *chatChecker) handle(ctx context.Context, b *bot.Bot, u *UpdateContext) {
	isChatActive, err := h.IsChatActive(ctx, u.Message.Chat.ID)
	if err != nil {
		logger.Error(fmt.Sprintf("ChatChecker.handle: %v", err))
	}
	u.IsChatActive = isChatActive

	h.nextHandle(ctx, b, u)
}
