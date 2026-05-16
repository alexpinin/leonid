package handler

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
)

type llmInquirer struct {
	basicHandler
	inquirer
}

func newLLMInquirer(i inquirer) *llmInquirer {
	return &llmInquirer{
		inquirer: i,
	}
}

type inquirer interface {
	InquireLLM(ctx context.Context, chatID int64, message string) (string, error)
}

func (h *llmInquirer) handle(ctx context.Context, b *bot.Bot, u *UpdateContext) error {
	reply, err := h.InquireLLM(ctx, u.Message.Chat.ID, u.Message.Text)
	if err != nil {
		return fmt.Errorf("llmInquirer.handle: %w", err)
	}
	u.LLMReply = reply
	return h.nextHandle(ctx, b, u)
}
