package handler

import (
	"fmt"
	"github.com/go-telegram/bot/models"
	"leonid/testutil"
	"testing"
)

type quotaManagerMock struct {
	runLog          *[]string
	getChatQuotaRes int
	getChatQuotaErr error
	useChatQuotaErr error
}

func (m *quotaManagerMock) GetChatQuota(chatID int64) (int, error) {
	*m.runLog = append(*m.runLog, fmt.Sprintf("GetChatQuota: %d", chatID))
	return m.getChatQuotaRes, m.getChatQuotaErr
}

func (m *quotaManagerMock) UseChatQuota(chatID int64) error {
	*m.runLog = append(*m.runLog, fmt.Sprintf("UseChatQuota: %d", chatID))
	return m.useChatQuotaErr
}

func Test_QuotaGuard_Handle(t *testing.T) {
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
			description:  "should call next handler if update quote",
			quotaManager: &quotaManagerMock{getChatQuotaRes: 1},
			given:        &UpdateContext{Update: update},
			expectedRunLog: []string{
				"GetChatQuota: 123",
				"UseChatQuota: 123",
				"Handle: " + testUpdateToStr(&UpdateContext{Update: update}),
			},
		},
		{
			description:  "should exit if GetChatQuota returns error",
			quotaManager: &quotaManagerMock{getChatQuotaErr: testError},
			given:        &UpdateContext{Update: update},
			expectedRunLog: []string{
				"GetChatQuota: 123",
			},
		},
		{
			description:  "should exit if quota is expired",
			quotaManager: &quotaManagerMock{getChatQuotaRes: 0},
			given:        &UpdateContext{Update: update},
			expectedRunLog: []string{
				"GetChatQuota: 123",
			},
		},
		{
			description:  "should exit if UseChatQuota returns error",
			quotaManager: &quotaManagerMock{getChatQuotaRes: 1, useChatQuotaErr: testError},
			given:        &UpdateContext{Update: update},
			expectedRunLog: []string{
				"GetChatQuota: 123",
				"UseChatQuota: 123",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			runLog := make([]string, 0)
			tc.quotaManager.runLog = &runLog
			h := NewQuotaGuard(tc.quotaManager)
			h.SetNext(&mockHandler{runLog: &runLog})

			h.Handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}
