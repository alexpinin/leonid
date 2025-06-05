package handler

import (
	"context"
	"leonid/src/internal/bot/repository"
	"leonid/src/internal/bot/service"
	"leonid/src/internal/common/db"

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

func NewBotHandler() *BotHandler {
	database := db.NewDB()
	configRepo := repository.NewConfigRepository(database)
	configService := service.NewConfigService(database, configRepo)
	quotaService := service.NewQuotaService()
	messageService := service.NewOpenAIMessageService(database, configRepo)

	handlers := []updateHandler{
		newInputGuard(),
		NewChatChecker(configService),
		newChatActivator(configService),
		newAuthGuard(),
		newCallGuard(configService),
		newQuotaGuard(quotaService),
		newMessageCleaner(configService),
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
