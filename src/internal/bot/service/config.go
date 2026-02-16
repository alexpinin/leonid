package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"leonid/src/internal/bot/dto"
	"leonid/src/internal/db"
	"leonid/src/internal/logger"
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

func (s *ConfigService) Activate(ctx context.Context, pass string, chatID int64) bool {
	err := s.executor.ExecuteInTx(func(tx *sql.Tx) error {
		config, err := s.configRepo.FindConfigByPass(tx, ctx, pass)
		if err != nil {
			return err
		}
		if config.PassValidBy.Before(time.Now()) {
			return errors.New("activation pass expired")
		}
		config.ChatID = chatID
		config.ChatActivatedAt = time.Now()
		err = s.configRepo.UpdateConfig(tx, ctx, config.ID, config)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Error(fmt.Sprintf("ConfigService.Activate: %v", err))
		return false
	}
	return true
}

func (s *ConfigService) IsChatActive(ctx context.Context, chatID int64) (bool, error) {
	config, err := s.configRepo.FindConfigByChatID(s.executor.Executor(), ctx, chatID)
	if err != nil {
		logger.Error(fmt.Sprintf("ConfigService.IsChatActive: %v", err))
		return false, err
	}
	return config.ChatActivatedAt.After(time.Time{}), nil
}

func (s *ConfigService) ListNicknames(ctx context.Context, chatID int64) []string {
	config, err := s.configRepo.FindConfigByChatID(s.executor.Executor(), ctx, chatID)
	if err != nil {
		logger.Error(fmt.Sprintf("ConfigService.ListNicknames: %v", err))
		return nil
	}
	return config.Nicknames
}

func (s *ConfigService) FindSystemPrompt(ctx context.Context, chatID int64) string {
	config, err := s.configRepo.FindConfigByChatID(s.executor.Executor(), ctx, chatID)
	if err != nil {
		logger.Error(fmt.Sprintf("ConfigService.FindSystemPrompt: %v", err))
		return ""
	}
	return config.SystemPrompt
}
