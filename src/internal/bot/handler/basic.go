package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type basicHandler struct {
	next updateHandler
}

func (h *basicHandler) setNext(n updateHandler) {
	h.next = n
}

func (h *basicHandler) getNext() updateHandler {
	return h.next
}

func (h *basicHandler) nextHandle(c context.Context, b *bot.Bot, u *models.Update, s *UpdateState) error {
	if h.next == nil {
		return nil
	}
	return h.next.handle(c, b, u, s)
}
