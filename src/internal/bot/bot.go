package bot

import (
	"context"
	"fmt"
	"leonid/src/internal/common/logger"

	"leonid/src/internal/bot/handler"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
)

var (
	botToken = os.Getenv("BOT_TOKEN")
)

func Start() {
	logger.Info(fmt.Sprintf("Starting bot"))

	h := handler.NewBotHandler()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(h.Handle),
		bot.WithAllowedUpdates(bot.AllowedUpdates{
			"message",
		}),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		logger.Panic(err.Error())
		return
	}

	b.Start(ctx)
}
