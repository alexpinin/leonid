package repository

import (
	"context"
	"database/sql"
	"fmt"
	"leonid/src/internal/db"
)

const upsertChatQuery = `
	INSERT INTO chat (chat_id) VALUES ($1)
	ON CONFLICT DO NOTHING;
`

type ChatRepository struct {
	db *db.DB
}

func NewChatRepository(db *db.DB) ChatRepository {
	return ChatRepository{
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
