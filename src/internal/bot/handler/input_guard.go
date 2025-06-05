package handler

import (
	"context"

	"github.com/go-telegram/bot"
)

type inputGuard struct {
	basicHandler
}

func newInputGuard() *inputGuard {
	return &inputGuard{}
}

func (h *inputGuard) handle(c context.Context, b *bot.Bot, u *UpdateContext) {
	if u == nil || u.Message == nil || u.Message.Chat.ID == 0 {
		return
	}
	h.nextHandle(c, b, u)
}
