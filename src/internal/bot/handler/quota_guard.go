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
	UseChatQuota(chatID int64) (bool, error)
}

func (h *quotaGuard) handle(c context.Context, b *bot.Bot, u *UpdateContext) error {
	quota, err := h.UseChatQuota(u.Message.Chat.ID)
	if err != nil {
		return fmt.Errorf("quotaManager.handle: %w", err)
	}
	if !quota {
		return nil
	}

	return h.nextHandle(c, b, u)
}
