package handler

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
)

type CallGuard struct {
	basicHandler
	nicknameProvider nicknameProvider
}

func NewCallGuard(np nicknameProvider) *CallGuard {
	return &CallGuard{nicknameProvider: np}
}

func (h *CallGuard) Handle(ctx context.Context, b *bot.Bot, u *UpdateContext) {
	nicknames := h.nicknameProvider.ListNicknames(ctx, u.Message.Chat.ID)
	message := strings.ToLower(u.Message.Text)
	for _, nickname := range nicknames {
		if strings.Contains(message, nickname) {
			h.nextHandle(ctx, b, u)
			return
		}
	}
}
