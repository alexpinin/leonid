package service

import (
	"context"
	"database/sql"
	"fmt"
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

func (s *ChatService) Activate(ctx context.Context, pass string, chatID int64) error {
	err := s.db.ExecInTx(ctx, func(tx *sql.Tx) error {
		configID, err := s.configRepo.FindValidPassConfigID(ctx, tx, pass, time.Now())
		if err != nil {
			return err
		}
		err = s.configRepo.ActivateChat(ctx, tx, configID, chatID, time.Now())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("ChatService.Activate: %w", err)
	}
	return nil
}

func (s *ChatService) ChatExists(ctx context.Context, chatID int64) (bool, error) {
	_, err := s.configRepo.FindValidChatConfigID(ctx, nil, chatID)
	if err != nil {
		return false, fmt.Errorf("ChatService.ChatExists: %w", err)
	}
	return true, nil
}
