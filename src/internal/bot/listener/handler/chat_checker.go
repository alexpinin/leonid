package handler

import (
	"context"

	"github.com/go-telegram/bot"
)

type ChatChecker struct {
	basicHandler
	storage chatCheckerStorage
}

func NewChatChecker(s chatCheckerStorage) ChatChecker {
	return ChatChecker{
		storage: s,
	}
}

type chatCheckerStorage interface {
	ChatExists(chatID int64) (bool, error)
}

func (h *ChatChecker) Handle(c context.Context, b *bot.Bot, u *UpdateContext) {
	exists, err := h.storage.ChatExists(u.Message.Chat.ID)
	if err != nil {
		return
	}
	if exists {
		u.IsChatActive = true
	}
	h.nextHandle(c, b, u)
}
