package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"leonid/src/internal/bot/service"
	"leonid/src/internal/logger"
)

type UpdateContext struct {
	IsChatActive bool
	IsPassActive bool
	LLMReply     string
}

type updateHandler interface {
	handle(context.Context, *bot.Bot, *models.Update, *UpdateContext) error
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
	audioEnabled bool,
) *BotHandler {
	handlers := []updateHandler{
		newInputGuard(),
		newChatChecker(configService),
		newChatActivator(configService),
		newAuthGuard(),
		newCallGuard(configService),
		newQuotaGuard(quotaService),
		newAudioReader(audioEnabled, audioService),
		newLLMInquirer(messageService),
		newReplySender(),
		nil,
	}

	for curr, next := 0, 1; next < len(handlers); curr, next = curr+1, next+1 {
		handlers[curr].setNext(handlers[next])
	}

	return &BotHandler{
		handlerHead: handlers[0],
	}
}

func (h *BotHandler) Handle(ctx context.Context, b *bot.Bot, u *models.Update) {
	uc := UpdateContext{}
	err := h.handlerHead.handle(ctx, b, u, &uc)
	if err != nil {
		logger.Error(err.Error())
	}
}
