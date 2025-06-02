package listener

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"leonid/src/internal/bot/listener/handler"
	"leonid/src/internal/bot/repository"
	"leonid/src/internal/bot/service"
	"leonid/src/internal/db"
)

type Handler struct {
	handlerHead handler.UpdateHandler
}

func NewHandler() *Handler {
	messageSender := handler.NewMessageSender()

	messageCleaner := handler.NewMessageCleaner()
	messageCleaner.SetNext(messageSender)

	nicknameChecker := handler.NewNicknameChecker()
	nicknameChecker.SetNext(messageCleaner)

	quotaService := service.NewQuotaService()
	quotaGuard := handler.NewQuotaGuard(quotaService)
	quotaGuard.SetNext(nicknameChecker)

	authGuard := handler.NewAuthGuard()
	authGuard.SetNext(quotaGuard)

	database := db.NewDB()
	passRepo := repository.NewPassRepository(database)
	chatRepo := repository.NewChatRepository(database)
	chatService := service.NewChatService(database, passRepo, chatRepo)

	chatActivator := handler.NewChatActivator(chatService)
	chatActivator.SetNext(authGuard)

	chatChecker := handler.NewChatChecker(chatService)
	chatChecker.SetNext(chatActivator)

	inputValidator := handler.NewInputValidator()
	inputValidator.SetNext(chatChecker)

	return &Handler{
		handlerHead: inputValidator,
	}
}

func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	uc := handler.UpdateContext{Update: update}
	h.handlerHead.Handle(ctx, b, &uc)
}
