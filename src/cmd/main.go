package main

import (
	"fmt"
	"os"

	"leonid/src/internal/bot"
	"leonid/src/internal/db"
	"leonid/src/internal/logger"
)

func main() {
	dbFile := os.Getenv("DB_FILE")
	if dbFile == "" {
		logger.Error("DB_FILE environment variable not set")
		os.Exit(1)
	}

	database, err := db.OpenDB(db.Config{DBFile: dbFile})
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to open database: %s", err.Error()))
		os.Exit(1)
	}
	defer database.Close()

	err = bot.Start(database)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to start bot: %s", err.Error()))
		os.Exit(1)
	}
}
