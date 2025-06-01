package listener

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"leonid/bot/listener/handler"
)

type Handler struct {
	handlerHead handler.UpdateHandler
}

func NewHandler() Handler {
	messageCleaner := handler.NewMessageCleaner()

	nicknameChecker := handler.NewNicknameChecker()
	nicknameChecker.SetNext(&messageCleaner)

	quotaGuard := handler.NewQuotaGuard(nil)
	quotaGuard.SetNext(&nicknameChecker)

	authGuard := handler.NewAuthGuard()
	authGuard.SetNext(&quotaGuard)

	chatActivator := handler.NewChatActivator(nil, nil)
	chatActivator.SetNext(&authGuard)

	chatChecker := handler.NewChatChecker(nil)
	chatChecker.SetNext(&chatActivator)

	inputValidator := handler.NewInputValidator()
	inputValidator.SetNext(&chatChecker)

	return Handler{
		handlerHead: &inputValidator,
	}
}

func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	uc := handler.UpdateContext{Update: update}
	h.handlerHead.Handle(ctx, b, &uc)
}
