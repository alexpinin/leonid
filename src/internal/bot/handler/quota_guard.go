package handler

import (
	"context"

	"github.com/go-telegram/bot"
)

type quotaGuard struct {
	basicHandler
	quotaManager
}

func newQuotaGuard(qm quotaManager) *quotaGuard {
	return &quotaGuard{
		quotaManager: qm,
	}
}

type quotaManager interface {
	UseChatQuota(chatID int64) bool
}

func (h *quotaGuard) handle(c context.Context, b *bot.Bot, u *UpdateContext) error {
	if !h.UseChatQuota(u.Message.Chat.ID) {
		return nil
	}
	return h.nextHandle(c, b, u)
}
