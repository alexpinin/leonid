package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-telegram/bot/models"

	"leonid/src/internal/testutil"
)

type chatActivatorMock struct {
	runLog      *[]string
	activateRes bool
	activateErr error
}

func (m *chatActivatorMock) Activate(_ context.Context, pass string, chatID int64) (bool, error) {
	*m.runLog = append(*m.runLog, fmt.Sprintf("Activate: %s, %d", pass, chatID))
	return m.activateRes, m.activateErr
}

func Test_chatActivator_handle(t *testing.T) {
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
			description:   "it should do nothing and call next handler if chat is already active",
			chatActivator: chatActivatorMock{},
			given:         &UpdateContext{Update: update, IsChatActive: true},
			expectedRunLog: []string{
				"handle: " + testUpdateToStr(&UpdateContext{Update: update, IsChatActive: true}),
			},
		},
		{
			description:   "it should set IsPassActive from Activate function result and call next handler",
			chatActivator: chatActivatorMock{activateRes: true},
			given:         &UpdateContext{Update: update},
			expectedRunLog: []string{
				"Activate: pass, 123",
				"handle: " + testUpdateToStr(&UpdateContext{Update: update, IsPassActive: true}),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			runLog := make([]string, 0)
			tc.chatActivator.runLog = &runLog
			moc := newChatActivator(&tc.chatActivator)
			moc.setNext(&mockHandler{runLog: &runLog})

			_ = moc.handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}
