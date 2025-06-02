package repository

import (
	"context"
	"database/sql"
	"fmt"
	"leonid/src/internal/db"
)

const (
	upsertChatQuery = `
	INSERT INTO chat (id) VALUES ($1)
	ON CONFLICT DO NOTHING
`
	chatExistsQuery = `
	SELECT COUNT(*) > 0
	FROM chat
	WHERE id = $1
`
)

type ChatRepository struct {
	db *db.DB
}

func NewChatRepository(db *db.DB) *ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (r *ChatRepository) UpsertChat(ctx context.Context, tx *sql.Tx, chatID int64) error {
	_, err := r.db.ExecTx(ctx, tx, upsertChatQuery, chatID)
	if err != nil {
		return fmt.Errorf("ChatRepository.UpsertChat: %w", err)
	}
	return nil
}

func (r *ChatRepository) ChatExists(ctx context.Context, tx *sql.Tx, chatID int64) (bool, error) {
	res := false
	row := r.db.QueryRowTx(ctx, tx, chatExistsQuery, chatID)
	err := row.Scan(&res)
	if err != nil {
		return false, fmt.Errorf("ChatRepository.ChatExists: %w", err)
	}
	return res, nil
}
