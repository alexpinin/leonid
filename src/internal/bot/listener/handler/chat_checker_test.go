package handler

import (
	"context"
	"fmt"
	"leonid/testutil"
	"testing"

	"github.com/go-telegram/bot/models"
)

type chatCheckerStorageMock struct {
	runLog        *[]string
	chatExistsRes bool
	chatExistsErr error
}

func (m *chatCheckerStorageMock) ChatExists(_ context.Context, chatID int64) (bool, error) {
	*m.runLog = append(*m.runLog, fmt.Sprintf("ChatExists: %d", chatID))
	return m.chatExistsRes, m.chatExistsErr
}

func Test_ChatChecker_Handle(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Chat: models.Chat{ID: 123},
		},
	}
	testCases := []struct {
		description    string
		storage        chatCheckerStorageMock
		given          *UpdateContext
		expectedRunLog []string
	}{
		{
			description: "should set chat as active and call next handler if chat exists",
			storage:     chatCheckerStorageMock{chatExistsRes: true},
			given:       &UpdateContext{Update: update},
			expectedRunLog: []string{
				"ChatExists: 123",
				"Handle: " + testUpdateToStr(&UpdateContext{Update: update, IsChatActive: true}),
			},
		},
		{
			description: "should stop and exit if ChatExists returns error",
			storage:     chatCheckerStorageMock{chatExistsErr: testError},
			given:       &UpdateContext{Update: update},
			expectedRunLog: []string{
				"ChatExists: 123",
			},
		},
		{
			description: "should set chat as not active and call next handler if chat doesn't exist",
			storage:     chatCheckerStorageMock{},
			given:       &UpdateContext{Update: update},
			expectedRunLog: []string{
				"ChatExists: 123",
				"Handle: " + testUpdateToStr(&UpdateContext{Update: update, IsChatActive: false}),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			runLog := make([]string, 0)
			tc.storage.runLog = &runLog
			h := NewChatChecker(&tc.storage)
			h.SetNext(&mockHandler{runLog: &runLog})

			h.Handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}
