package listener

import (
	"context"
	"leonid/src/internal/bot/listener/handler"
	"leonid/src/internal/bot/repository"
	"leonid/src/internal/bot/service"
	"leonid/src/internal/common/db"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct {
	handlerHead handler.UpdateHandler
}

func NewHandler() *Handler {
	database := db.NewDB()
	configRepo := repository.NewConfigRepository(database)
	configService := service.NewConfigService(database, configRepo)
	quotaService := service.NewQuotaService()
	messageService := service.NewOpenAIMessageService(configRepo)

	inputGuard := handler.NewInputGuard()
	chatChecker := handler.NewChatChecker(configService)
	chatActivator := handler.NewChatActivator(configService)
	authGuard := handler.NewAuthGuard()
	callGuard := handler.NewCallGuard(configService)
	quotaGuard := handler.NewQuotaGuard(quotaService)
	messageCleaner := handler.NewMessageCleaner(configService)
	messageSender := handler.NewMessageSender(messageService)

	inputGuard.SetNext(chatChecker)
	chatChecker.SetNext(chatActivator)
	chatActivator.SetNext(authGuard)
	authGuard.SetNext(callGuard)
	callGuard.SetNext(quotaGuard)
	quotaGuard.SetNext(messageCleaner)
	messageCleaner.SetNext(messageSender)
	messageSender.SetNext(nil)

	return &Handler{
		handlerHead: inputGuard,
	}
}

func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	uc := handler.UpdateContext{Update: update}
	h.handlerHead.Handle(ctx, b, &uc)
}
