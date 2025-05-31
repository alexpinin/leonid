package handler

import (
	"context"

	"github.com/go-telegram/bot"
)

type InputValidator struct {
	basicHandler
}

func (h *InputValidator) Handle(c context.Context, b *bot.Bot, u *UpdateContext) {
	if u == nil || u.Message == nil {
		return
	}
	u.IsInputValid = true
	h.handleNext(c, b, u)
}
