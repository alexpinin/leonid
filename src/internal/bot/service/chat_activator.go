package service

import (
	"context"
	"database/sql"
	"fmt"
	"leonid/src/internal/bot/repository"
	"leonid/src/internal/db"
	"time"
)

type ChatActivator struct {
	db       *db.DB
	passRepo repository.PassRepository
	chatRepo repository.ChatRepository
}

func NewChatActivator() ChatActivator { // todo pass dependencies
	return ChatActivator{
		db:       db.NewDB(),
		passRepo: repository.NewPassRepository(db.NewDB()),
		chatRepo: repository.NewChatRepository(db.NewDB()),
	}
}

func (p *ChatActivator) Activate(ctx context.Context, pass string, chatID int64) error {
	err := p.db.ExecInTx(ctx, func(tx *sql.Tx) error {
		exists, err := p.passRepo.PassExists(ctx, tx, pass, time.Now())
		if err != nil {
			return err
		}
		if !exists {
			return errUnknownPass
		}
		err = p.passRepo.DeletePass(ctx, tx, pass)
		if err != nil {
			return err
		}
		err = p.chatRepo.UpsertChat(ctx, tx, chatID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("ChatActivator.Activate: %w", err)
	}
	return nil
}
