package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

type Config struct {
	DBFile string
}

func OpenDB(cfg Config) (*DB, error) {
	db, err := sql.Open("sqlite", cfg.DBFile)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
}

func (db *DB) ExecTx(ctx context.Context, tx *sql.Tx, query string, args ...any) (sql.Result, error) {
	if tx == nil {
		return db.ExecContext(ctx, query, args...)
	}
	return tx.ExecContext(ctx, query, args...)
}

func (db *DB) QueryTx(ctx context.Context, tx *sql.Tx, query string, args ...any) (*sql.Rows, error) {
	if tx == nil {
		return db.QueryContext(ctx, query, args...)
	}
	return tx.QueryContext(ctx, query, args...)
}

func (db *DB) QueryRowTx(ctx context.Context, tx *sql.Tx, query string, args ...any) *sql.Row {
	if tx == nil {
		return db.QueryRowContext(ctx, query, args...)
	}
	return tx.QueryRowContext(ctx, query, args...)
}

func (db *DB) ExecInTx(ctx context.Context, f func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
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
