package handler

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
)

type MessageCleaner struct {
	basicHandler
	nicknameProvider nicknameProvider
}

func NewMessageCleaner(np nicknameProvider) *MessageCleaner {
	return &MessageCleaner{nicknameProvider: np}
}

func (h *MessageCleaner) Handle(ctx context.Context, b *bot.Bot, u *UpdateContext) {
	nicknames := h.nicknameProvider.ListNicknames(ctx, u.Message.Chat.ID)
	message := strings.TrimSpace(strings.ToLower(u.Message.Text))
	for _, nickname := range nicknames {
		if strings.Contains(message, nickname) {
			u.Message.Text = strings.ReplaceAll(message, nickname, "")
		}
	}
	h.nextHandle(ctx, b, u)
}
