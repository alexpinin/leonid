package handler

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
)

type callGuard struct {
	basicHandler
	nicknameProvider
}

func newCallGuard(np nicknameProvider) *callGuard {
	return &callGuard{nicknameProvider: np}
}

type nicknameProvider interface {
	ListNicknames(ctx context.Context, chatID int64) []string
}

func (h *callGuard) handle(ctx context.Context, b *bot.Bot, u *UpdateContext) {
	nicknames := h.ListNicknames(ctx, u.Message.Chat.ID)
	message := strings.ToLower(u.Message.Text)
	for _, nickname := range nicknames {
		if nickname != "" && strings.Contains(message, nickname) {
			h.nextHandle(ctx, b, u)
			return
		}
	}
}
