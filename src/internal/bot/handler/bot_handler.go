package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"leonid/src/internal/bot/service"
	"leonid/src/internal/logger"
)

type UpdateContext struct {
	*models.Update
	IsChatActive bool
	IsPassActive bool
}

type updateHandler interface {
	handle(context.Context, *bot.Bot, *UpdateContext) error
	setNext(updateHandler)
	getNext() updateHandler
}

type BotHandler struct {
	handlerHead updateHandler
}

func NewBotHandler(
	configService *service.ConfigService,
	quotaService *service.QuotaService,
	audioService *service.AudioService,
	messageService *service.OpenAIService,
) *BotHandler {
	handlers := []updateHandler{
		newInputGuard(),
		newChatChecker(configService),
		newChatActivator(configService),
		newAuthGuard(),
		newCallGuard(configService),
		newQuotaGuard(quotaService),
		newAudioReader(true, audioService),
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
	err := h.handlerHead.handle(ctx, b, &uc)
	if err != nil {
		logger.Error(err.Error())
	}
}
