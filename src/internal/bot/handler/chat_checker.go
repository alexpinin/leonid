package handler

import (
	"context"

	"github.com/go-telegram/bot"
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
	IsChatActive(ctx context.Context, chatID int64) bool
}

func (h *chatChecker) handle(ctx context.Context, b *bot.Bot, u *UpdateContext) {
	u.IsChatActive = h.IsChatActive(ctx, u.Message.Chat.ID)
	h.nextHandle(ctx, b, u)
}
