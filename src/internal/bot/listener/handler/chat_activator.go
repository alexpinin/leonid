package handler

import (
	"context"
	"github.com/go-telegram/bot"
)

type ChatActivator struct {
	basicHandler
	activatorService chatActivator
}

func NewChatActivator(ca chatActivator) ChatActivator {
	return ChatActivator{
		activatorService: ca,
	}
}

type chatActivator interface {
	Activate(context.Context, string, int64) error
}

func (h *ChatActivator) Handle(ctx context.Context, b *bot.Bot, u *UpdateContext) {
	if u.IsChatActive {
		h.nextHandle(ctx, b, u)
		return
	}
	err := h.activatorService.Activate(ctx, u.Message.Text, u.Message.Chat.ID)
	if err != nil {
		u.IsPassActive = false
		// todo log
	} else {
		u.IsPassActive = true
	}
	h.nextHandle(ctx, b, u)
}
