package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"leonid/src/internal/db"
)

type Config struct {
	ID                  string
	ChatID              int64
	ChatActivatedAt     time.Time
	Pass                string
	PassValidBy         time.Time
	Nicknames           []string
	SystemPrompt        string
	MessagePrompt       string
	ConversationContext string
}

type ConfigRepo struct {
	db *db.DB
}

func NewConfigRepo(db *db.DB) *ConfigRepo {
	return &ConfigRepo{
		db: db,
	}
}

func (r *ConfigRepo) FindConfigByPass(ctx context.Context, tx *sql.Tx, pass string) (Config, error) {
	query := `
		SELECT
			id,
			pass,
			pass_valid_by,
			chat_id,
			chat_activated_at,
			nicknames,
			system_prompt,
			message_prompt,
			conversation_context
		FROM config
		WHERE pass = $1
	`
	row := r.db.QueryRowTx(ctx, tx, query, pass)
	config, err := scanConfig(row)
	if err != nil {
		return Config{}, fmt.Errorf("ConfigRepo.FindConfigByPass: %w", err)
	}
	return config, nil
}

func (r *ConfigRepo) FindConfigByChatID(ctx context.Context, tx *sql.Tx, chatID int64) (Config, error) {
	query := `
		SELECT
			id,
			pass,
			pass_valid_by,
			chat_id,
			chat_activated_at,
			nicknames,
			system_prompt,
			message_prompt,
			conversation_context
		FROM config
		WHERE chat_id = $1
	`
	row := r.db.QueryRowTx(ctx, tx, query, chatID)
	config, err := scanConfig(row)
	if err != nil {
		return Config{}, fmt.Errorf("ConfigRepo.FindConfigByChatID: %w", err)
	}
	return config, nil
}

func (r *ConfigRepo) UpdateConfig(ctx context.Context, tx *sql.Tx, configID string, c Config) error {
	const query = `
		UPDATE config
		SET chat_id = $2,
			chat_activated_at = $3,
			conversation_context = $4
		WHERE id = $1
	`
	_, err := r.db.ExecTx(ctx, tx, query,
		configID, c.ChatID, c.ChatActivatedAt.Unix(), c.ConversationContext)
	if err != nil {
		return fmt.Errorf("ConfigRepo.UpdateConfig: %w", err)
	}
	return nil
}

func scanConfig(r *sql.Row) (Config, error) {
	var config Config
	var passValidByUnix, chatActivatedAt int64
	var nicknamesStr string
	err := r.Scan(
		&config.ID,
		&config.Pass,
		&passValidByUnix,
		&config.ChatID,
		&chatActivatedAt,
		&nicknamesStr,
		&config.SystemPrompt,
		&config.MessagePrompt,
		&config.ConversationContext,
	)
	if err != nil {
		return Config{}, err
	}
	config.PassValidBy = time.Unix(passValidByUnix, 0)
	config.ChatActivatedAt = time.Unix(chatActivatedAt, 0)
	config.Nicknames = strings.Split(nicknamesStr, ",")
	return config, nil
}
