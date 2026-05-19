package handler

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
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

func (h *llmInquirer) handle(ctx context.Context, b *bot.Bot, u *models.Update, s *UpdateState) error {
	reply, err := h.InquireLLM(ctx, u.Message.Chat.ID, u.Message.Text)
	if err != nil {
		return fmt.Errorf("llmInquirer.handle: %w", err)
	}
	s.LLMReply = reply
	return h.nextHandle(ctx, b, u, s)
}
