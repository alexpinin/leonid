package handler

import (
	"fmt"
	"leonid/testutil"
	"testing"

	"github.com/go-telegram/bot/models"
)

type quotaManagerMock struct {
	runLog *[]string
	result bool
}

func (m *quotaManagerMock) UseChatQuota(chatID int64) bool {
	*m.runLog = append(*m.runLog, fmt.Sprintf("UseChatQuota: %d", chatID))
	return m.result
}

func Test_quotaGuard_handle(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Chat: models.Chat{
				ID: 123,
			},
		},
	}
	testCases := []struct {
		description    string
		quotaManager   *quotaManagerMock
		given          *UpdateContext
		expectedRunLog []string
	}{
		{
			description:  "should call next handler if UseChatQuota returns true",
			quotaManager: &quotaManagerMock{result: true},
			given:        &UpdateContext{Update: update},
			expectedRunLog: []string{
				"UseChatQuota: 123",
				"handle: " + testUpdateToStr(&UpdateContext{Update: update}),
			},
		},
		{
			description:  "should not call next handler and exit if UseChatQuota returns false",
			quotaManager: &quotaManagerMock{result: false},
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

			moc.handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}
