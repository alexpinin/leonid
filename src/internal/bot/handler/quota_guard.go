package handler

import (
	"context"
	"fmt"

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
	UseChatQuota(c context.Context, chatID int64) error
}

func (h *quotaGuard) handle(c context.Context, b *bot.Bot, u *UpdateContext) error {
	err := h.UseChatQuota(c, u.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("quotaManager.handle: %w", err)
	}

	return h.nextHandle(c, b, u)
}
