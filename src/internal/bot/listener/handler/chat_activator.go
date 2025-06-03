package handler

import (
	"context"

	"github.com/go-telegram/bot"
)

type ChatActivator struct {
	basicHandler
	activator chatActivator
}

func NewChatActivator(cha chatActivator) *ChatActivator {
	return &ChatActivator{
		activator: cha,
	}
}

type chatActivator interface {
	Activate(context.Context, string, int64) bool
}

func (h *ChatActivator) Handle(ctx context.Context, b *bot.Bot, u *UpdateContext) {
	if u.IsChatActive {
		h.nextHandle(ctx, b, u)
		return
	}
	u.IsPassActive = h.activator.Activate(ctx, u.Message.Text, u.Message.Chat.ID)
	h.nextHandle(ctx, b, u)
}
