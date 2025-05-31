package handler

import (
	"context"
	"time"

	"github.com/go-telegram/bot"
)

type ChatActivator struct {
	basicHandler
	passStorage chatActivatorPassStorage
	chatStorage chatActivatorChatStorage
}

func NewChatActivator(ps chatActivatorPassStorage, cs chatActivatorChatStorage) ChatActivator {
	return ChatActivator{
		passStorage: ps,
		chatStorage: cs,
	}
}

type chatActivatorPassStorage interface {
	PassExists(pass string, byDate time.Time) (bool, error)
	DeletePass(pass string) error
}

type chatActivatorChatStorage interface {
	ActivateChat(chatID int64) error
}

func (h *ChatActivator) Handle(c context.Context, b *bot.Bot, u *UpdateContext) {
	if u.IsChatActive {
		h.handleNext(c, b, u)
		return
	}

	ok, err := h.passStorage.PassExists(u.Message.Text, time.Now())
	if err != nil {
		return
	}
	if !ok {
		return
	}

	err = h.chatStorage.ActivateChat(u.Message.Chat.ID)
	if err != nil {
		return
	}
	err = h.passStorage.DeletePass(u.Message.Text)
	if err != nil {
		return
	}
	u.IsChatActive = true

	h.handleNext(c, b, u)
}
