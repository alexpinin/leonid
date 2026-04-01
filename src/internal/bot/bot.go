package bot

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"

	"leonid/src/internal/bot/handler"
	"leonid/src/internal/bot/repo"
	"leonid/src/internal/bot/service"
	"leonid/src/internal/db"

	"github.com/go-telegram/bot"
)

type Config struct {
	BotToken    string
	LLMProvider string
	LLMToken    string
	LLMModel    string
}

func Start(database *sql.DB, cfg Config) error {
	url, err := baseURL(cfg.LLMProvider)
	if err != nil {
		return err
	}

	executor := db.NewAppQueryExecutor(database)

	configRepo := repo.NewConfigRepo()
	quotaRepo := repo.NewQuotaRepo()

	configService := service.NewConfigService(executor, configRepo)
	quotaService := service.NewQuotaService(executor, quotaRepo)
	audioService := service.NewAudioService()

	llmClient := service.NewOpenAIClient(service.OpenAIConfig{
		BaseURL: url,
		Token:   cfg.LLMToken,
		Model:   cfg.LLMModel,
	})
	messageService := service.NewOpenAIService(executor, configRepo, llmClient)

	botHandler := handler.NewBotHandler(
		configService,
		quotaService,
		audioService,
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

	b, err := bot.New(cfg.BotToken, opts...)
	if err != nil {
		return err
	}

	b.Start(ctx)
	return nil
}

func baseURL(provider string) (string, error) {
	switch provider {
	case "openai":
		return "https://api.openai.com/v1/", nil
	case "deepseek":
		return "https://api.deepseek.com", nil
	default:
		return "", fmt.Errorf("unknown LLM provider")
	}
}
