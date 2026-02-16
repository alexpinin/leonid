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

type ConfigService struct {
	executor   db.QueryExecutor
	configRepo ConfigRepo
}

func NewConfigService(
	qe db.QueryExecutor,
	cr ConfigRepo,
) *ConfigService {
	return &ConfigService{
		executor:   qe,
		configRepo: cr,
	}
}

type ConfigRepo interface {
	FindConfigByPass(db.Executor, context.Context, string) (dto.Config, error)
	FindConfigByChatID(db.Executor, context.Context, int64) (dto.Config, error)
	UpdateConfig(db.Executor, context.Context, string, dto.Config) error
}

func (s *ConfigService) Activate(ctx context.Context, pass string, chatID int64) (bool, error) {
	err := s.executor.ExecuteInTx(func(tx *sql.Tx) error {
		config, err := s.configRepo.FindConfigByPass(tx, ctx, pass)
		if err != nil {
			return fmt.Errorf("cannot find config: %w", err)
		}

		if config.PassValidBy.Before(time.Now()) {
			return errors.New("activation pass expired")
		}

		config.ChatID = chatID
		config.ChatActivatedAt = time.Now()

		err = s.configRepo.UpdateConfig(tx, ctx, config.ID, config)
		if err != nil {
			return fmt.Errorf("cannot activate config: %w", err)
		}

		return nil
	})
	if err != nil {
		return false, fmt.Errorf("ConfigService.Activate: %w", err)
	}
	return true, nil
}

func (s *ConfigService) IsChatActive(ctx context.Context, chatID int64) (bool, error) {
	config, err := s.configRepo.FindConfigByChatID(s.executor.Executor(), ctx, chatID)
	if err != nil {
		return false, fmt.Errorf("ConfigService.IsChatActive: %w", err)
	}
	return config.ChatActivatedAt.After(time.Time{}), nil
}

func (s *ConfigService) ListNicknames(ctx context.Context, chatID int64) ([]string, error) {
	config, err := s.configRepo.FindConfigByChatID(s.executor.Executor(), ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("ConfigService.ListNicknames: %w", err)
	}
	return config.Nicknames, nil
}
