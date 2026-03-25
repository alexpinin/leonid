package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"leonid/src/internal/bot/dto"
	"leonid/src/internal/db"
)

type QuotaService struct {
	executor  db.QueryExecutor
	quotaRepo QuotaRepo
}

type QuotaRepo interface {
	FindQuotaByChatID(db.Executor, context.Context, int64) (dto.Quota, error)
	UpsertQuota(db.Executor, context.Context, dto.Quota) error
}

func NewQuotaService(
	ex db.QueryExecutor,
	qr QuotaRepo,
) *QuotaService {
	return &QuotaService{
		executor:  ex,
		quotaRepo: qr,
	}
}

const maxQuotaPerDay = 100

func (s *QuotaService) UseChatQuota(c context.Context, chatID int64) error {
	err := s.executor.ExecuteInTx(func(tx *sql.Tx) error {
		quota, err := s.quotaRepo.FindQuotaByChatID(tx, c, chatID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to find quota: %w", err)
		}

		if errors.Is(err, sql.ErrNoRows) {
			quota = dto.Quota{
				ChatID:        chatID,
				LastResetDate: time.Now().UTC().Truncate(24 * time.Hour),
				Remaining:     maxQuotaPerDay,
			}
		}

		now := time.Now().UTC().Truncate(24 * time.Hour)
		if quota.LastResetDate.Truncate(24 * time.Hour).Before(now) {
			quota.LastResetDate = now
			quota.Remaining = maxQuotaPerDay
		}

		if quota.Remaining <= 0 {
			return errors.New("quota exceeded")
		}

		quota.Remaining--

		err = s.quotaRepo.UpsertQuota(tx, c, quota)
		if err != nil {
			return fmt.Errorf("failed to update quota: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("QuotaService.UseChatQuota: chat ID %d: %w", chatID, err)
	}

	return nil
}
