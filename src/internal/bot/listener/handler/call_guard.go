package handler

import (
	"context"
	"github.com/go-telegram/bot"
	"strings"
)

type CallGuard struct {
	basicHandler
}

func NewCallGuard() *CallGuard {
	return &CallGuard{}
}

func (h *CallGuard) Handle(ctx context.Context, b *bot.Bot, u *UpdateContext) {
	message := strings.ToLower(u.Message.Text)
	for _, nickname := range nicknames {
		if strings.Contains(message, nickname) {
			h.nextHandle(ctx, b, u)
			return
		}
	}
}
