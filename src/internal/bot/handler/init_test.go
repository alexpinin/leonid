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

func (m *mockHandler) handle(_ context.Context, _ *bot.Bot, uc *UpdateContext) error {
	*m.runLog = append(*m.runLog, fmt.Sprintf("handle: %s", testUpdateToStr(uc)))
	return nil
}

func (m *mockHandler) setNext(updateHandler) {
}

func (m *mockHandler) getNext() updateHandler {
	return nil
}

func testUpdateToStr(uc *UpdateContext) string {
	s, _ := json.Marshal(uc)
	return string(s)
}
