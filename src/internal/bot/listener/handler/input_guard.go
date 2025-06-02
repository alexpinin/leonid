package handler

import (
	"context"

	"github.com/go-telegram/bot"
)

type InputGuard struct {
	basicHandler
}

func NewInputGuard() *InputGuard {
	return &InputGuard{}
}

func (h *InputGuard) Handle(c context.Context, b *bot.Bot, u *UpdateContext) {
	if u == nil || u.Message == nil {
		return
	}
	h.nextHandle(c, b, u)
}
