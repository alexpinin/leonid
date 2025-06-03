package handler

import (
	"context"
	"fmt"
	"leonid/testutil"
	"testing"

	"github.com/go-telegram/bot/models"
)

type chatCheckerStorageMock struct {
	runLog *[]string
	result bool
}

func (m *chatCheckerStorageMock) IsChatActive(_ context.Context, chatID int64) bool {
	*m.runLog = append(*m.runLog, fmt.Sprintf("IsChatActive: %d", chatID))
	return m.result
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
			description: "should set IsChatActive from the IsChatActive function result and call next handler",
			storage:     chatCheckerStorageMock{result: true},
			given:       &UpdateContext{Update: update},
			expectedRunLog: []string{
				"IsChatActive: 123",
				"Handle: " + testUpdateToStr(&UpdateContext{Update: update, IsChatActive: true}),
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
