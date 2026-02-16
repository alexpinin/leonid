package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
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

func (h *callGuard) handle(ctx context.Context, b *bot.Bot, u *UpdateContext) error {
	nicknames, err := h.ListNicknames(ctx, u.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("callGuard.handle: %w", err)
	}

	message := strings.ToLower(u.Message.Text)
	for _, nickname := range nicknames {
		if nickname != "" && strings.Contains(message, nickname) {
			return h.nextHandle(ctx, b, u)
		}
	}

	if u.Message.ReplyToMessage == nil || u.Message.ReplyToMessage.From == nil {
		return nil
	}

	replyToNickname := strings.ToLower(u.Message.ReplyToMessage.From.FirstName)
	for _, nickname := range nicknames {
		if nickname == replyToNickname {
			return h.nextHandle(ctx, b, u)
		}
	}

	return nil
}
