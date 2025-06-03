package repository

import (
	"context"
	"database/sql"
	"fmt"
	"leonid/src/internal/db"
	"time"
)

type Config struct {
	ID              string
	ChatID          int64
	ChatActivatedAt time.Time
	Pass            string
	PassValidBy     time.Time
	Nicknames       []string
	SystemPrompt    string
}

type ConfigRepository struct {
	db *db.DB
}

func NewConfigRepository(db *db.DB) *ConfigRepository {
	return &ConfigRepository{
		db: db,
	}
}

const findValidPassConfigIDQuery = `
	SELECT id
	FROM config
	WHERE pass = $1
	  AND pass_valid_by >= $2
`

func (r *ConfigRepository) FindValidPassConfigID(ctx context.Context, tx *sql.Tx, pass string, validBy time.Time) (string, error) {
	validByUnix := validBy.UTC().Unix()
	row := r.db.QueryRowTx(ctx, tx, findValidPassConfigIDQuery, pass, validByUnix)
	var configID string
	err := row.Scan(&configID)
	if err != nil {
		return "", fmt.Errorf("ConfigRepository.FindValidPassConfigID: %w", err)
	}
	return configID, nil
}

const activateChatQuery = `
	UPDATE config
	SET chat_id = $2, chat_activated_at = $3
	WHERE id = $1
`

func (r *ConfigRepository) ActivateChat(ctx context.Context, tx *sql.Tx, configID string, chatID int64, activatedAt time.Time) error {
	activatedAtUnix := activatedAt.UTC().Unix()
	_, err := r.db.ExecTx(ctx, tx, activateChatQuery, configID, chatID, activatedAtUnix)
	if err != nil {
		return fmt.Errorf("ConfigRepository.ActivateChat: %w", err)
	}
	return nil
}

const findValidChatConfigIDQuery = `
	SELECT id
	FROM config
	WHERE chat_id = $1
`

func (r *ConfigRepository) FindValidChatConfigID(ctx context.Context, tx *sql.Tx, chatID int64) (string, error) {
	row := r.db.QueryRowTx(ctx, tx, findValidChatConfigIDQuery, chatID)
	var configID string
	err := row.Scan(&configID)
	if err != nil {
		return "", fmt.Errorf("ConfigRepository.FindValidChatConfigID: %w", err)
	}
	return configID, nil
}
