package db

import (
	"context"
	"database/sql"
	"fmt"
)

type QueryExecutor interface {
	ExecuteInTx(f func(tx *sql.Tx) error) error
	Executor() Executor
}

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type AppQueryExecutor struct {
	db *sql.DB
}

func NewAppQueryExecutor(db *sql.DB) *AppQueryExecutor {
	return &AppQueryExecutor{db: db}
}

func (t *AppQueryExecutor) ExecuteInTx(f func(tx *sql.Tx) error) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	err = f(tx)
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return fmt.Errorf("AppQueryExecutor.ExecuteInTx: transaction rollback caused by error: %w failed: %w", err, txErr)
		}
		return fmt.Errorf("AppQueryExecutor.ExecuteInTx: %w", err)
	}
	return tx.Commit()
}

func (t *AppQueryExecutor) Executor() Executor {
	return t.db
}

type MockQueryExecutor struct {
}

func (t *MockQueryExecutor) ExecuteInTx(f func(tx *sql.Tx) error) error {
	return f(nil)
}

func (t *MockQueryExecutor) Executor() Executor {
	return nil
}
