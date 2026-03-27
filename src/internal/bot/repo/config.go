package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"leonid/src/internal/bot/dto"
	"leonid/src/internal/db"
)

type ConfigRepo struct {
}

func NewConfigRepo() *ConfigRepo {
	return &ConfigRepo{}
}

func (*ConfigRepo) FindConfigByPass(ex db.Executor, ctx context.Context, pass string) (dto.Config, error) {
	q := `
		SELECT
		    id,
			chat_id,
			pass,
			pass_valid_by,
			chat_activated_at,
			nicknames,
			system_prompt,
			conversation_context
		FROM config
		WHERE pass = $1
	`
	row := ex.QueryRowContext(ctx, q, pass)
	config, err := scanConfig(row)
	if err != nil {
		return dto.Config{}, fmt.Errorf("ConfigRepo.FindConfigByPass: %w", err)
	}
	return config, nil
}

func (*ConfigRepo) FindConfigByChatID(ex db.Executor, ctx context.Context, chatID int64) (dto.Config, error) {
	q := `
		SELECT
		    id,
			chat_id,
			pass,
			pass_valid_by,
			chat_activated_at,
			nicknames,
			system_prompt,
			conversation_context
		FROM config
		WHERE chat_id = $1
	`
	row := ex.QueryRowContext(ctx, q, chatID)
	config, err := scanConfig(row)
	if err != nil {
		return dto.Config{}, fmt.Errorf("ConfigRepo.FindConfigByChatID: %w", err)
	}
	return config, nil
}

func (*ConfigRepo) UpdateConfig(ex db.Executor, ctx context.Context, c dto.Config) error {
	q := `
		UPDATE config
		SET chat_id = $2,
		    chat_activated_at = $3,
			conversation_context = $4
		WHERE id = $1
	`
	_, err := ex.ExecContext(ctx, q, c.ID, c.ChatID, c.ChatActivatedAt.Unix(), c.ConversationHistory)
	if err != nil {
		return fmt.Errorf("ConfigRepo.UpdateConfig: %w", err)
	}
	return nil
}

func scanConfig(r *sql.Row) (dto.Config, error) {
	var config dto.Config
	var chatID *int64
	var passValidByUnix, chatActivatedAtUnix int64
	var nicknamesStr string
	err := r.Scan(
		&config.ID,
		&chatID,
		&config.Pass,
		&passValidByUnix,
		&chatActivatedAtUnix,
		&nicknamesStr,
		&config.SystemPrompt,
		&config.ConversationHistory,
	)
	if err != nil {
		return dto.Config{}, err
	}
	if chatID != nil {
		config.ChatID = *chatID
	}
	config.PassValidBy = time.Unix(passValidByUnix, 0)
	config.ChatActivatedAt = time.Unix(chatActivatedAtUnix, 0)
	config.Nicknames = strings.Split(nicknamesStr, ",")
	return config, nil
}
