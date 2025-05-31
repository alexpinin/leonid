package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-telegram/bot"
)

type mockHandler struct {
	runLog *[]string
}

func (m *mockHandler) Handle(_ context.Context, _ *bot.Bot, uc *UpdateContext) {
	*m.runLog = append(*m.runLog, fmt.Sprintf("Handle: %s", testUpdateToStr(uc)))
}

func (m *mockHandler) SetNext(UpdateHandler) {
}

func testUpdateToStr(uc *UpdateContext) string {
	s, _ := json.Marshal(uc)
	return string(s)
}

var testError = errors.New("test")
