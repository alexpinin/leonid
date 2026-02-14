package handler

import (
	"context"

	"leonid/src/internal/bot/repo"
	"leonid/src/internal/bot/service"
	"leonid/src/internal/db"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type UpdateContext struct {
	*models.Update
	IsChatActive bool
	IsPassActive bool
}

type updateHandler interface {
	handle(context.Context, *bot.Bot, *UpdateContext)
	setNext(updateHandler)
	getNext() updateHandler
}

type BotHandler struct {
	handlerHead updateHandler
}

func NewBotHandler(db *db.DB) *BotHandler {
	configRepo := repo.NewConfigRepo(db)
	configService := service.NewConfigService(db, configRepo)
	quotaService := service.NewQuotaService()
	messageService := service.NewDeepSeekMessageService(db, configRepo)

	handlers := []updateHandler{
		newInputGuard(),
		NewChatChecker(configService),
		newChatActivator(configService),
		newAuthGuard(),
		newCallGuard(configService),
		newQuotaGuard(quotaService),
		newMessageSender(messageService),
		nil,
	}

	for curr, next := 0, 1; next < len(handlers); curr, next = curr+1, next+1 {
		handlers[curr].setNext(handlers[next])
	}

	return &BotHandler{
		handlerHead: handlers[0],
	}
}

func (h *BotHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	uc := UpdateContext{Update: update}
	h.handlerHead.handle(ctx, b, &uc)
}
