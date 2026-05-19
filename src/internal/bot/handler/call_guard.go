package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type callGuard struct {
	basicHandler
	nicknameProvider
}

func newCallGuard(np nicknameProvider) *callGuard {
	return &callGuard{nicknameProvider: np}
}

type nicknameProvider interface {
	ListNicknames(ctx context.Context, chatID int64) ([]string, error)
}

func (h *callGuard) handle(ctx context.Context, b *bot.Bot, u *models.Update, s *UpdateState) error {
	nicknames, err := h.ListNicknames(ctx, u.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("callGuard.handle: %w", err)
	}
	for i, nickname := range nicknames {
		nicknames[i] = strings.ToLower(nickname)
	}

	message := strings.ToLower(u.Message.Text)
	for _, nickname := range nicknames {
		if nickname != "" && strings.Contains(message, nickname) {
			return h.nextHandle(ctx, b, u, s)
		}
	}

	if u.Message.ReplyToMessage == nil || u.Message.ReplyToMessage.From == nil {
		return nil
	}

	if u.Message.ReplyToMessage.From.ID == b.ID() {
		return h.nextHandle(ctx, b, u, s)
	}

	return nil
}
