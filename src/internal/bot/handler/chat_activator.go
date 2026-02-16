package handler

import (
	"context"
	"fmt"

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
	Activate(context.Context, string, int64) (bool, error)
}

func (h *chatActivator) handle(ctx context.Context, b *bot.Bot, u *UpdateContext) error {
	if u.IsChatActive {
		return h.nextHandle(ctx, b, u)
	}

	isPassActive, err := h.Activate(ctx, u.Message.Text, u.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("chatActivator.handle: %w", err)
	}
	u.IsPassActive = isPassActive

	return h.nextHandle(ctx, b, u)
}
