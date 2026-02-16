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

func (h *chatActivator) handle(ctx context.Context, b *bot.Bot, u *UpdateContext) error {
	if u.IsChatActive {
		return h.nextHandle(ctx, b, u)
	}
	u.IsPassActive = h.Activate(ctx, u.Message.Text, u.Message.Chat.ID)
	return h.nextHandle(ctx, b, u)
}
