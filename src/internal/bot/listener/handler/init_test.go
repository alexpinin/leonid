package handler

import (
	"context"
	"encoding/json"
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

func (m *mockHandler) GetNext() UpdateHandler {
	return nil
}

func testUpdateToStr(uc *UpdateContext) string {
	s, _ := json.Marshal(uc)
	return string(s)
}
