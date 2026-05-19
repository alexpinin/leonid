package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type mockHandler struct {
	handleCount int
}

func (m *mockHandler) handle(context.Context, *bot.Bot, *models.Update, *UpdateState) error {
	m.handleCount++
	return nil
}

func (m *mockHandler) setNext(updateHandler) {
}

func (m *mockHandler) getNext() updateHandler {
	return nil
}
