package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-telegram/bot/models"

	"leonid/src/internal/testutil"
)

func TestChatCheckerHandle(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Chat: models.Chat{ID: 123},
		},
	}
	testCases := []struct {
		description    string
		storage        mockChatCheckerStorage
		given          *UpdateContext
		expectedRunLog []string
	}{
		{
			description: "should set IsChatActive from the IsChatActive function result and call next handler",
			storage:     mockChatCheckerStorage{isChatActiveRes: true},
			given:       &UpdateContext{Update: update},
			expectedRunLog: []string{
				"IsChatActive: 123",
				"handle: " + testUpdateToStr(&UpdateContext{Update: update, IsChatActive: true}),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			runLog := make([]string, 0)
			tc.storage.runLog = &runLog
			moc := NewChatChecker(&tc.storage)
			moc.setNext(&mockHandler{runLog: &runLog})

			_ = moc.handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}

type mockChatCheckerStorage struct {
	runLog          *[]string
	isChatActiveRes bool
	isChatActiveErr error
}

func (m *mockChatCheckerStorage) IsChatActive(_ context.Context, chatID int64) (bool, error) {
	*m.runLog = append(*m.runLog, fmt.Sprintf("IsChatActive: %d", chatID))
	return m.isChatActiveRes, m.isChatActiveErr
}
