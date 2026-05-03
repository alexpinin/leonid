package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"leonid/src/internal/bot"
	"leonid/src/internal/db"
	"leonid/src/internal/logger"

	"github.com/joho/godotenv"
)

func main() {
	if isDevMode() {
		logger.Info("development mode")
		mustLoadEnvFile(".env")
	}

	dbFile := mustLoad("DB_FILE")
	database, err := db.OpenDB(db.Config{DBFile: dbFile})
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to open database: %s", err.Error()))
		os.Exit(1)
	}
	defer database.Close()

	cfg := bot.Config{
		BotToken:      mustLoad("BOT_TOKEN"),
		LLMProvider:   mustLoad("LLM_PROVIDER"),
		LLMToken:      mustLoad("LLM_TOKEN"),
		LLMModel:      mustLoad("LLM_MODEL"),
		TranscribeURL: mustLoad("TRANSCRIBE_URL"),
	}

	logger.Info("Starting bot")

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

func mustLoadEnvFile(envFile string) {
	wd, err := os.Getwd()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	path := filepath.Join(wd, envFile)

	err = godotenv.Overload(path)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load environment variables: %s", err.Error()))
		os.Exit(1)
	}
}

func isDevMode() bool {
	var dev bool
	flag.BoolVar(&dev, "dev", false, "Development mode")
	flag.Parse()
	return dev
}
