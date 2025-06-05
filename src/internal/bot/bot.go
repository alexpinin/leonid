package bot

import (
	"context"

	"leonid/src/internal/bot/handler"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
)

var (
	botToken = os.Getenv("BOT_TOKEN")
)

func Start() {
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
		panic(err)
	}

	b.Start(ctx)
}
