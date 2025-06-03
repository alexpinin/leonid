package handler

import (
	"context"

	"github.com/go-telegram/bot"
)

type AuthGuard struct {
	basicHandler
}

func NewAuthGuard() *AuthGuard {
	return &AuthGuard{}
}

func (h *AuthGuard) Handle(c context.Context, b *bot.Bot, u *UpdateContext) {
	if !u.IsChatActive && !u.IsPassActive {
		return
	}
	h.nextHandle(c, b, u)
}
