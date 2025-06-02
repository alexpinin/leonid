package db

import (
	"context"
	"fmt"
	"os"

	"database/sql"

	_ "modernc.org/sqlite"
)

type DB struct {
	db *sql.DB
}

var dbOnlyInstance *DB

func NewDB() *DB {
	if dbOnlyInstance != nil {
		return dbOnlyInstance
	}
	dbFile := os.Getenv("DB_FILE")
	if dbFile == "" {
		panic("DB_FILE environment variable not set")
	}
	var err error
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		panic(err)
	}
	dbOnlyInstance = &DB{db: db}
	return dbOnlyInstance
}

func (db *DB) ExecTx(ctx context.Context, tx *sql.Tx, query string, args ...any) (sql.Result, error) {
	if tx == nil {
		return db.db.ExecContext(ctx, query, args...)
	}
	return tx.ExecContext(ctx, query, args...)
}

func (db *DB) QueryTx(ctx context.Context, tx *sql.Tx, query string, args ...any) (*sql.Rows, error) {
	if tx == nil {
		return db.db.QueryContext(ctx, query, args...)
	}
	return tx.QueryContext(ctx, query, args...)
}

func (db *DB) QueryRowTx(ctx context.Context, tx *sql.Tx, query string, args ...any) *sql.Row {
	if tx == nil {
		return db.db.QueryRowContext(ctx, query, args...)
	}
	return tx.QueryRowContext(ctx, query, args...)
}

func (db *DB) ExecInTx(ctx context.Context, f func(tx *sql.Tx) error) error {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	err = f(tx)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return fmt.Errorf("transaction rollback caused by error: %w failed: %w", err, txErr)
		}
		return err
	}
	return tx.Commit()
}
