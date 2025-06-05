package handler

import (
	"context"

	"github.com/go-telegram/bot"
)

type chatActivator struct {
	basicHandler
	chActivator
}

func newChatActivator(cha chActivator) *chatActivator {
	return &chatActivator{
		chActivator: cha,
	}
}

type chActivator interface {
	Activate(context.Context, string, int64) bool
}

func (h *chatActivator) handle(ctx context.Context, b *bot.Bot, u *UpdateContext) {
	if u.IsChatActive {
		h.nextHandle(ctx, b, u)
		return
	}
	u.IsPassActive = h.Activate(ctx, u.Message.Text, u.Message.Chat.ID)
	h.nextHandle(ctx, b, u)
}
