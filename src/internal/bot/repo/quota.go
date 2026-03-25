package repo

import (
	"context"
	"fmt"

	"leonid/src/internal/bot/dto"
	"leonid/src/internal/db"
)

type QuotaRepo struct {
}

func NewQuotaRepo() *QuotaRepo {
	return &QuotaRepo{}
}

func (*QuotaRepo) FindQuotaByChatID(ex db.Executor, ctx context.Context, chatID int64) (dto.Quota, error) {
	q := `
		SELECT last_reset_date, remaining       
		FROM quota WHERE chat_id = $1
	`
	r := ex.QueryRowContext(ctx, q, chatID)
	var quota dto.Quota
	var lastResetDateUnix int64
	err := r.Scan(
		&lastResetDateUnix,
		&quota.Remaining,
	)
	if err != nil {
		return dto.Quota{}, fmt.Errorf("QuotaRepo.FindQuotaByChatID: %w", err)
	}

	return quota, nil
}

func (*QuotaRepo) UpsertQuota(ex db.Executor, ctx context.Context, quota dto.Quota) error {
	q := `
		INSERT INTO quota (chat_id, last_reset_date, remaining)                                                                              
	  	VALUES ($1, $2, $3)                                                                                                                     
	  	ON CONFLICT (chat_id) DO UPDATE SET                                                                                                  
			last_reset_date = excluded.last_reset_date,                                                                                      
		  	remaining = excluded.remaining;  
	`
	_, err := ex.ExecContext(ctx, q, quota.ChatID, quota.LastResetDate.Unix(), quota.Remaining)
	if err != nil {
		return fmt.Errorf("QuotaRepo.UpsertQuota: %w", err)
	}

	return nil
}
