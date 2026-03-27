package handler

import (
	"context"

	"github.com/go-telegram/bot"
)

type authGuard struct {
	basicHandler
}

func newAuthGuard() *authGuard {
	return &authGuard{}
}

func (h *authGuard) handle(c context.Context, b *bot.Bot, u *UpdateContext) error {
	if !u.IsChatActive && !u.IsPassActive {
		return nil
	}
	return h.nextHandle(c, b, u)
}
