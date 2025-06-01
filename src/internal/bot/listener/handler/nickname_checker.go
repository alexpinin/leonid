package handler

import (
	"context"
	"github.com/go-telegram/bot"
	"strings"
)

type NicknameChecker struct {
	basicHandler
}

func NewNicknameChecker() NicknameChecker {
	return NicknameChecker{}
}

func (h *NicknameChecker) Handle(c context.Context, b *bot.Bot, u *UpdateContext) {
	message := strings.ToLower(u.Message.Text)
	for _, nickname := range nicknames {
		if strings.Contains(message, nickname) {
			u.IsNicknameCall = true
			break
		}
	}
	h.nextHandle(c, b, u)
}
