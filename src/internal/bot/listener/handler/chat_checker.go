package handler

import (
	"context"

	"github.com/go-telegram/bot"
)

type ChatChecker struct {
	basicHandler
	checker chatChecker
}

func NewChatChecker(ch chatChecker) *ChatChecker {
	return &ChatChecker{
		checker: ch,
	}
}

type chatChecker interface {
	IsChatActive(ctx context.Context, chatID int64) bool
}

func (h *ChatChecker) Handle(ctx context.Context, b *bot.Bot, u *UpdateContext) {
	u.IsChatActive = h.checker.IsChatActive(ctx, u.Message.Chat.ID)
	h.nextHandle(ctx, b, u)
}
