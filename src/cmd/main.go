package main

import (
	"fmt"
	"os"

	"leonid/src/internal/bot"
	"leonid/src/internal/db"
	"leonid/src/internal/logger"
)

func main() {
	dbFile := mustLoad("DB_FILE")
	database, err := db.OpenDB(db.Config{DBFile: dbFile})
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to open database: %s", err.Error()))
		os.Exit(1)
	}
	defer database.Close()

	cfg := bot.Config{
		BotToken:    mustLoad("BOT_TOKEN"),
		LLMProvider: mustLoad("LLM_PROVIDER"),
		LLMToken:    mustLoad("LLM_TOKEN"),
		LLMModel:    mustLoad("LLM_MODEL"),
	}

	logger.Info(fmt.Sprintf("Starting bot"))

	err = bot.Start(database, cfg)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to start bot: %s", err.Error()))
		os.Exit(1)
	}
}

func mustLoad(key string) string {
	val := os.Getenv(key)
	if val == "" {
		logger.Error(fmt.Sprintf("environment variable %s not set", key))
		os.Exit(1)
	}
	return val
}
