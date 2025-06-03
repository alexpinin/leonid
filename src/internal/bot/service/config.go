package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"leonid/src/internal/bot/repository"
	"leonid/src/internal/common/db"
	"leonid/src/internal/common/logger"
	"time"
)

type ConfigService struct {
	db         *db.DB
	configRepo *repository.ConfigRepository
}

func NewConfigService(
	db *db.DB,
	cr *repository.ConfigRepository,
) *ConfigService {
	return &ConfigService{
		db:         db,
		configRepo: cr,
	}
}

func (s *ConfigService) Activate(ctx context.Context, pass string, chatID int64) bool {
	err := s.db.ExecInTx(ctx, func(tx *sql.Tx) error {
		config, err := s.configRepo.FindConfigByPass(ctx, tx, pass)
		if err != nil {
			return err
		}
		if config.PassValidBy.Before(time.Now()) {
			return errors.New("activation pass expired")
		}
		config.ChatID = chatID
		config.ChatActivatedAt = time.Now()
		err = s.configRepo.UpdateConfig(ctx, tx, config.ID, config)
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
	_, err := s.configRepo.FindConfigByChatID(ctx, nil, chatID)
	if err != nil {
		logger.Error(fmt.Sprintf("ConfigService.IsChatActive: %v", err))
		return false
	}
	return true
}

func (s *ConfigService) ListNicknames(ctx context.Context, chatID int64) []string {
	config, err := s.configRepo.FindConfigByChatID(ctx, nil, chatID)
	if err != nil {
		logger.Error(fmt.Sprintf("ConfigService.ListNicknames: %v", err))
		return nil
	}
	return config.Nicknames
}

func (s *ConfigService) FindSystemPrompt(ctx context.Context, chatID int64) string {
	config, err := s.configRepo.FindConfigByChatID(ctx, nil, chatID)
	if err != nil {
		logger.Error(fmt.Sprintf("ConfigService.FindSystemPrompt: %v", err))
		return ""
	}
	return config.SystemPrompt
}
