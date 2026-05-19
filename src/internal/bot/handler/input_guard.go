package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type inputGuard struct {
	basicHandler
}

func newInputGuard() *inputGuard {
	return &inputGuard{}
}

func (h *inputGuard) handle(c context.Context, b *bot.Bot, u *models.Update, uc *UpdateContext) error {
	if u == nil || u.Message == nil || u.Message.Chat.ID == 0 {
		return nil
	}
	return h.nextHandle(c, b, u, uc)
}
