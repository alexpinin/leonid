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
	GetChatQuota(chatID int64) (int, error)
	UseChatQuota(chatID int64) error
}

func (h *QuotaGuard) Handle(c context.Context, b *bot.Bot, u *UpdateContext) {
	quota, err := h.quota.GetChatQuota(u.Message.Chat.ID)
	if err != nil {
		return
	}
	if quota <= 0 {
		// todo ass security logs
		// todo pass through and return a response (if so don't decrease quota)
		return
	}
	err = h.quota.UseChatQuota(u.Message.Chat.ID)
	if err != nil {
		return
	}
	h.nextHandle(c, b, u)
}
