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
	db       *db.DB
	passRepo *repository.PassRepository
	chatRepo *repository.ChatRepository
}

func NewChatService(
	db *db.DB,
	passRepo *repository.PassRepository,
	chatRepo *repository.ChatRepository,
) *ChatService {
	return &ChatService{
		db:       db,
		passRepo: passRepo,
		chatRepo: chatRepo,
	}
}

func (s *ChatService) Activate(ctx context.Context, pass string, chatID int64) error {
	err := s.db.ExecInTx(ctx, func(tx *sql.Tx) error {
		p, err := s.passRepo.FindPass(ctx, tx, pass, time.Now())
		if err != nil {
			return err
		}
		err = s.passRepo.DeletePass(ctx, tx, pass)
		if err != nil {
			return err
		}
		chat := repository.Chat{
			ID:           chatID,
			Nicknames:    p.Nicknames,
			SystemPrompt: p.SystemPrompt,
		}
		err = s.chatRepo.UpsertChat(ctx, tx, chat)
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
	exists, err := s.chatRepo.ChatExists(ctx, nil, chatID)
	if err != nil {
		return false, fmt.Errorf("ChatService.ChatExists: %w", err)
	}
	return exists, nil
}
