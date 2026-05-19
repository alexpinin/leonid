package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type authGuard struct {
	basicHandler
}

func newAuthGuard() *authGuard {
	return &authGuard{}
}

func (h *authGuard) handle(c context.Context, b *bot.Bot, u *models.Update, uc *UpdateContext) error {
	if !uc.IsChatActive && !uc.IsPassActive {
		return nil
	}
	return h.nextHandle(c, b, u, uc)
}
