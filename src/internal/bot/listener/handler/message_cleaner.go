package handler

import (
	"context"
	"github.com/go-telegram/bot"
	"strings"
)

type MessageCleaner struct {
	basicHandler
}

func NewMessageCleaner() *MessageCleaner {
	return &MessageCleaner{}
}

var nicknames = []string{
	"леня",
	"лёня",
	"ленька",
	"лёнька",
	"леонид",
}

func (h *MessageCleaner) Handle(c context.Context, b *bot.Bot, u *UpdateContext) {
	message := strings.TrimSpace(strings.ToLower(u.Message.Text))
	for _, nickname := range nicknames {
		if strings.Contains(message, nickname) {
			u.Message.Text = strings.ReplaceAll(message, nickname, "")
		}
	}
	h.nextHandle(c, b, u)
}
