package handler

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
)

type messageCleaner struct {
	basicHandler
	nicknameProvider
}

func newMessageCleaner(np nicknameProvider) *messageCleaner {
	return &messageCleaner{nicknameProvider: np}
}

func (h *messageCleaner) handle(ctx context.Context, b *bot.Bot, u *UpdateContext) {
	nicknames := h.ListNicknames(ctx, u.Message.Chat.ID)
	message := strings.TrimSpace(strings.ToLower(u.Message.Text))
	for _, nickname := range nicknames {
		if strings.Contains(message, nickname) {
			u.Message.Text = strings.ReplaceAll(message, nickname, "")
		}
	}
	h.nextHandle(ctx, b, u)
}
