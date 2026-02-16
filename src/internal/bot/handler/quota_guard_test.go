package handler

import (
	"fmt"
	"testing"

	"github.com/go-telegram/bot/models"

	"leonid/src/internal/testutil"
)

func TestQuotaGuardHandle(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Chat: models.Chat{
				ID: 123,
			},
		},
	}
	testCases := []struct {
		description    string
		quotaManager   *mockQuotaManager
		given          *UpdateContext
		expectedRunLog []string
	}{
		{
			description:  "should call next handler if UseChatQuota returns true",
			quotaManager: &mockQuotaManager{useChatQuotaRes: true},
			given:        &UpdateContext{Update: update},
			expectedRunLog: []string{
				"UseChatQuota: 123",
				"handle: " + testUpdateToStr(&UpdateContext{Update: update}),
			},
		},
		{
			description:  "should not call next handler and exit if UseChatQuota returns false",
			quotaManager: &mockQuotaManager{useChatQuotaRes: false},
			given:        &UpdateContext{Update: update},
			expectedRunLog: []string{
				"UseChatQuota: 123",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			runLog := make([]string, 0)
			tc.quotaManager.runLog = &runLog
			moc := newQuotaGuard(tc.quotaManager)
			moc.setNext(&mockHandler{runLog: &runLog})

			_ = moc.handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}

type mockQuotaManager struct {
	runLog          *[]string
	useChatQuotaRes bool
	useChatQuotaErr error
}

func (m *mockQuotaManager) UseChatQuota(chatID int64) (bool, error) {
	*m.runLog = append(*m.runLog, fmt.Sprintf("UseChatQuota: %d", chatID))
	return m.useChatQuotaRes, m.useChatQuotaErr
}
