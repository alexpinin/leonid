package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-telegram/bot"
)

type mockHandler struct {
	updateContext *UpdateContext
	runLog        *[]string
}

func (m *mockHandler) Handle(_ context.Context, _ *bot.Bot, uc *UpdateContext) {
	if m.runLog != nil {
		con := ""
		if uc == nil || uc.Message == nil {
			con = "nil"
		} else {
			con = fmt.Sprintf("{message: %s, chat: %d}", uc.Message.Text, uc.Message.Chat.ID)
		}
		*m.runLog = append(*m.runLog, fmt.Sprintf("Handle: %s", con))
	}

	m.updateContext = uc
}

func (m *mockHandler) SetNext(UpdateHandler) {
}

var testError = errors.New("test")
