package bot

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"

	"leonid/src/internal/bot/repo"
	"leonid/src/internal/bot/service"
	"leonid/src/internal/db"
	"leonid/src/internal/logger"

	"leonid/src/internal/bot/handler"

	"github.com/go-telegram/bot"
)

var (
	botToken = os.Getenv("BOT_TOKEN")
)

func Start(database *sql.DB) error {
	logger.Info(fmt.Sprintf("Starting bot"))

	executor := db.NewAppQueryExecutor(database)

	configRepo := repo.NewConfigRepo()

	configService := service.NewConfigService(executor, configRepo)
	quotaService := service.NewQuotaService()
	messageService := service.NewOpenAIMessageService(executor, configRepo)

	botHandler := handler.NewBotHandler(
		configService,
		quotaService,
		messageService,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(botHandler.Handle),
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
