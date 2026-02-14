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

func (r *ConfigRepo) FindConfigByPass(ex db.Executor, ctx context.Context, pass string) (dto.Config, error) {
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
	row := ex.QueryRowContext(ctx, query, pass)
	config, err := scanConfig(row)
	if err != nil {
		return dto.Config{}, fmt.Errorf("ConfigRepo.FindConfigByPass: %w", err)
	}
	return config, nil
}

func (r *ConfigRepo) FindConfigByChatID(ex db.Executor, ctx context.Context, chatID int64) (dto.Config, error) {
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
	row := ex.QueryRowContext(ctx, query, chatID)
	config, err := scanConfig(row)
	if err != nil {
		return dto.Config{}, fmt.Errorf("ConfigRepo.FindConfigByChatID: %w", err)
	}
	return config, nil
}

func (r *ConfigRepo) UpdateConfig(ex db.Executor, ctx context.Context, configID string, c dto.Config) error {
	const query = `
		UPDATE config
		SET chat_id = $2,
			chat_activated_at = $3,
			conversation_context = $4
		WHERE id = $1
	`
	_, err := ex.ExecContext(ctx, query,
		configID, c.ChatID, c.ChatActivatedAt.Unix(), c.ConversationContext)
	if err != nil {
		return fmt.Errorf("ConfigRepo.UpdateConfig: %w", err)
	}
	return nil
}

func scanConfig(r *sql.Row) (dto.Config, error) {
	var config dto.Config
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
		return dto.Config{}, err
	}
	config.PassValidBy = time.Unix(passValidByUnix, 0)
	config.ChatActivatedAt = time.Unix(chatActivatedAt, 0)
	config.Nicknames = strings.Split(nicknamesStr, ",")
	return config, nil
}
