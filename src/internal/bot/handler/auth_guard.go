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

func (h *authGuard) handle(c context.Context, b *bot.Bot, u *models.Update, s *UpdateState) error {
	if !s.IsChatActive && !s.IsPassActive {
		return nil
	}
	return h.nextHandle(c, b, u, s)
}
