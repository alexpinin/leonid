package handler

import (
	"context"

	"github.com/go-telegram/bot"
)

type basicHandler struct {
	next UpdateHandler
}

func (h *basicHandler) SetNext(n UpdateHandler) {
	h.next = n
}

func (h *basicHandler) nextHandle(c context.Context, b *bot.Bot, u *UpdateContext) {
	if h.next == nil {
		return
	}
	h.next.Handle(c, b, u)
}
