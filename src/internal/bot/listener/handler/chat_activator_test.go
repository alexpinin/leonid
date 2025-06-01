package handler

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot/models"
	"leonid/testutil"
	"testing"
)

type chatActivatorMock struct {
	runLog *[]string
	err    error
}

func (m *chatActivatorMock) Activate(_ context.Context, pass string, chatID int64) error {
	*m.runLog = append(*m.runLog, fmt.Sprintf("Activate: %s, %d", pass, chatID))
	return m.err
}

func Test_ChatActivator_Handle(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Text: "pass",
			Chat: models.Chat{
				ID: 123,
			},
		},
	}
	testCases := []struct {
		description    string
		chatActivator  chatActivatorMock
		given          *UpdateContext
		expectedRunLog []string
	}{
		{
			description:   "it should activate chat if it's not active and call next handler",
			chatActivator: chatActivatorMock{},
			given:         &UpdateContext{Update: update},
			expectedRunLog: []string{
				"Activate: pass, 123",
				"Handle: " + testUpdateToStr(&UpdateContext{Update: update, IsPassActive: true}),
			},
		},
		{
			description:   "it should do nothing if chat is already active",
			chatActivator: chatActivatorMock{},
			given:         &UpdateContext{Update: update, IsChatActive: true},
			expectedRunLog: []string{
				"Handle: " + testUpdateToStr(&UpdateContext{Update: update, IsChatActive: true}),
			},
		},
		{
			description:   "it should mark chat as not active if Activate returns error",
			chatActivator: chatActivatorMock{err: testError},
			given:         &UpdateContext{Update: update},
			expectedRunLog: []string{
				"Activate: pass, 123",
				"Handle: " + testUpdateToStr(&UpdateContext{Update: update}),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			runLog := make([]string, 0)
			tc.chatActivator.runLog = &runLog
			h := NewChatActivator(&tc.chatActivator)
			h.SetNext(&mockHandler{runLog: &runLog})

			h.Handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}
