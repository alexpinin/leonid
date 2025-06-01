package listener

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	handler2 "leonid/src/internal/bot/listener/handler"
)

type Handler struct {
	handlerHead handler2.UpdateHandler
}

func NewHandler() Handler {
	messageCleaner := handler2.NewMessageCleaner()

	nicknameChecker := handler2.NewNicknameChecker()
	nicknameChecker.SetNext(&messageCleaner)

	quotaGuard := handler2.NewQuotaGuard(nil)
	quotaGuard.SetNext(&nicknameChecker)

	authGuard := handler2.NewAuthGuard()
	authGuard.SetNext(&quotaGuard)

	chatActivator := handler2.NewChatActivator(nil, nil)
	chatActivator.SetNext(&authGuard)

	chatChecker := handler2.NewChatChecker(nil)
	chatChecker.SetNext(&chatActivator)

	inputValidator := handler2.NewInputValidator()
	inputValidator.SetNext(&chatChecker)

	return Handler{
		handlerHead: &inputValidator,
	}
}

func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	uc := handler2.UpdateContext{Update: update}
	h.handlerHead.Handle(ctx, b, &uc)
}
