package repository

import (
	"context"
	"database/sql"
	"fmt"
	"leonid/src/internal/db"
	"strings"
)

const (
	upsertChatQuery = `
	INSERT INTO chat (id, nicknames, system_prompt) 
	VALUES ($1, $2, $3)
	ON CONFLICT DO UPDATE
	SET nicknames = $2, 
	    system_prompt = $3
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

func (r *ChatRepository) UpsertChat(ctx context.Context, tx *sql.Tx, chat Chat) error {
	nicknames := strings.Join(chat.Nicknames, ",")
	_, err := r.db.ExecTx(ctx, tx, upsertChatQuery, chat.ID, nicknames, chat.SystemPrompt)
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

type Chat struct {
	ID           int64
	Nicknames    []string
	SystemPrompt string
}
