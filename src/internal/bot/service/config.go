package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"leonid/src/internal/bot/repository"
	"leonid/src/internal/db"
	"leonid/src/internal/logger"
)

type ConfigService struct {
	db   *db.DB
	repo configRepo
}

func NewConfigService(
	db *db.DB,
	cr configRepo,
) *ConfigService {
	return &ConfigService{
		db:   db,
		repo: cr,
	}
}

type configRepo interface {
	FindConfigByPass(ctx context.Context, tx *sql.Tx, pass string) (repository.Config, error)
	FindConfigByChatID(ctx context.Context, tx *sql.Tx, chatID int64) (repository.Config, error)
	UpdateConfig(ctx context.Context, tx *sql.Tx, configID string, c repository.Config) error
}

func (s *ConfigService) Activate(ctx context.Context, pass string, chatID int64) bool {
	err := s.db.ExecInTx(ctx, func(tx *sql.Tx) error {
		config, err := s.repo.FindConfigByPass(ctx, tx, pass)
		if err != nil {
			return err
		}
		if config.PassValidBy.Before(time.Now()) {
			return errors.New("activation pass expired")
		}
		config.ChatID = chatID
		config.ChatActivatedAt = time.Now()
		err = s.repo.UpdateConfig(ctx, tx, config.ID, config)
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

func (s *ConfigService) IsChatActive(ctx context.Context, chatID int64) bool {
	_, err := s.repo.FindConfigByChatID(ctx, nil, chatID)
	if err != nil {
		logger.Error(fmt.Sprintf("ConfigService.IsChatActive: %v", err))
		return false
	}
	return true
}

func (s *ConfigService) ListNicknames(ctx context.Context, chatID int64) []string {
	config, err := s.repo.FindConfigByChatID(ctx, nil, chatID)
	if err != nil {
		logger.Error(fmt.Sprintf("ConfigService.ListNicknames: %v", err))
		return nil
	}
	return config.Nicknames
}

func (s *ConfigService) FindSystemPrompt(ctx context.Context, chatID int64) string {
	config, err := s.repo.FindConfigByChatID(ctx, nil, chatID)
	if err != nil {
		logger.Error(fmt.Sprintf("ConfigService.FindSystemPrompt: %v", err))
		return ""
	}
	return config.SystemPrompt
}
