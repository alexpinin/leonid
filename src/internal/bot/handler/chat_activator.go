package handler

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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

func (h *chatActivator) handle(ctx context.Context, b *bot.Bot, u *models.Update, uc *UpdateContext) error {
	if uc.IsChatActive {
		return h.nextHandle(ctx, b, u, uc)
	}

	isPassActive, err := h.Activate(ctx, u.Message.Text, u.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("chatActivator.handle: %w", err)
	}
	uc.IsPassActive = isPassActive

	return h.nextHandle(ctx, b, u, uc)
}
