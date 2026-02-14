package db

import (
	"context"
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Config struct {
	DBFile string
}

func OpenDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite", cfg.DBFile)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
