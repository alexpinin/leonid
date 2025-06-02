package bot

import (
	"context"
	"github.com/go-telegram/bot"
	"leonid/src/internal/bot/listener"
	"os"
	"os/signal"
)

var (
	botToken = os.Getenv("BOT_TOKEN")
)

func Start() {
	handler := listener.NewHandler()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler.Handle),
		bot.WithAllowedUpdates(bot.AllowedUpdates{
			"message",
		}),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}
