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
	connStr := cfg.DBFile + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"
	db, err := sql.Open("sqlite", connStr)
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
