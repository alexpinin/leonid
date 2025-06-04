package handler

import (
	"context"

	"github.com/go-telegram/bot"
)

type QuotaGuard struct {
	basicHandler
	quota quotaManager
}

func NewQuotaGuard(q quotaManager) *QuotaGuard {
	return &QuotaGuard{
		quota: q,
	}
}

type quotaManager interface {
	UseChatQuota(chatID int64) bool
}

func (h *QuotaGuard) Handle(c context.Context, b *bot.Bot, u *UpdateContext) {
	if !h.quota.UseChatQuota(u.Message.Chat.ID) {
		return
	}
	h.nextHandle(c, b, u)
}
