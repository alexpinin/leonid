package handler

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
)

type replySender struct {
	basicHandler
}

func newReplySender() *replySender {
	return &replySender{}
}

func (h *replySender) handle(ctx context.Context, b *bot.Bot, u *UpdateContext) error {
	telegramParams := bot.SendMessageParams{
		ChatID: u.Message.Chat.ID,
		Text:   u.LLMReply,
	}
	_, err := b.SendMessage(ctx, &telegramParams)
	if err != nil {
		return fmt.Errorf("replySender.handle: %w", err)
	}

	return h.nextHandle(ctx, b, u)
}
