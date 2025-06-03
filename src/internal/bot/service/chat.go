package service

import (
	"context"
	"database/sql"
	"errors"
	"leonid/src/internal/bot/repository"
	"leonid/src/internal/db"
	"time"
)

type ChatService struct {
	db         *db.DB
	configRepo *repository.ConfigRepository
}

func NewChatService(
	db *db.DB,
	configRepo *repository.ConfigRepository,
) *ChatService {
	return &ChatService{
		db:         db,
		configRepo: configRepo,
	}
}

func (s *ChatService) Activate(ctx context.Context, pass string, chatID int64) bool {
	err := s.db.ExecInTx(ctx, func(tx *sql.Tx) error {
		config, err := s.configRepo.FindConfigByPass(ctx, tx, pass)
		if err != nil {
			return err
		}
		if config.ChatActivatedAt.Before(time.Now()) {
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
		return false
	}
	return true
}

func (s *ChatService) IsChatActive(ctx context.Context, chatID int64) bool {
	_, err := s.configRepo.FindConfigByChatID(ctx, nil, chatID)
	if err != nil {
		return false
	}
	return true
}

func (s *ChatService) ListNicknames(ctx context.Context, chatID int64) []string {
	config, err := s.configRepo.FindConfigByChatID(ctx, nil, chatID)
	if err != nil {
		return nil
	}
	return config.Nicknames
}
