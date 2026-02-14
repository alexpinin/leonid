package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"leonid/src/internal/db"
	"leonid/src/internal/logger"

	"leonid/src/internal/bot/handler"

	"github.com/go-telegram/bot"
)

var (
	botToken = os.Getenv("BOT_TOKEN")
)

func Start(database *db.DB) error {
	logger.Info(fmt.Sprintf("Starting bot"))

	h := handler.NewBotHandler(database)

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
		return err
	}

	b.Start(ctx)
	return nil
}
