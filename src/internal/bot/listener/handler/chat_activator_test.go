package handler

import (
	"context"
	"fmt"
	"leonid/testutil"
	"testing"

	"github.com/go-telegram/bot/models"
)

type chatActivatorMock struct {
	runLog *[]string
	result bool
}

func (m *chatActivatorMock) Activate(_ context.Context, pass string, chatID int64) bool {
	*m.runLog = append(*m.runLog, fmt.Sprintf("Activate: %s, %d", pass, chatID))
	return m.result
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
			description:   "it should do nothing and call next handler if chat is already active",
			chatActivator: chatActivatorMock{},
			given:         &UpdateContext{Update: update, IsChatActive: true},
			expectedRunLog: []string{
				"Handle: " + testUpdateToStr(&UpdateContext{Update: update, IsChatActive: true}),
			},
		},
		{
			description:   "it should set IsPassActive from Activate function result and call next handler",
			chatActivator: chatActivatorMock{result: true},
			given:         &UpdateContext{Update: update},
			expectedRunLog: []string{
				"Activate: pass, 123",
				"Handle: " + testUpdateToStr(&UpdateContext{Update: update, IsPassActive: true}),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			runLog := make([]string, 0)
			tc.chatActivator.runLog = &runLog
			moc := NewChatActivator(&tc.chatActivator)
			moc.SetNext(&mockHandler{runLog: &runLog})

			moc.Handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}
