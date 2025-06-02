package handler

import (
	"context"

	"github.com/go-telegram/bot"
)

type InputValidator struct {
	basicHandler
}

func NewInputValidator() *InputValidator {
	return &InputValidator{}
}

func (h *InputValidator) Handle(c context.Context, b *bot.Bot, u *UpdateContext) {
	if u == nil || u.Message == nil {
		return
	}
	u.IsInputValid = true
	h.nextHandle(c, b, u)
}
