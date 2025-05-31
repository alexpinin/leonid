package handler

import (
	"fmt"
	"leonid/testutil"
	"testing"
	"time"

	"github.com/go-telegram/bot/models"
)

type chatActivatorPassStorageMock struct {
	runLog        *[]string
	passExistsRes bool
	passExistsErr error
	deletePassErr error
}

func (m *chatActivatorPassStorageMock) PassExists(pass string, _ time.Time) (bool, error) {
	*m.runLog = append(*m.runLog, fmt.Sprintf("PassExists: %s", pass))
	return m.passExistsRes, m.passExistsErr
}

func (m *chatActivatorPassStorageMock) DeletePass(pass string) error {
	*m.runLog = append(*m.runLog, fmt.Sprintf("DeletePass: %s", pass))
	return m.deletePassErr
}

type chatActivatorChatStorageMock struct {
	runLog          *[]string
	activateChatErr error
}

func (m *chatActivatorChatStorageMock) ActivateChat(chatID int64) error {
	*m.runLog = append(*m.runLog, fmt.Sprintf("ActivateChat: %d", chatID))
	return m.activateChatErr
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
		passStorage    chatActivatorPassStorageMock
		chatStorage    chatActivatorChatStorageMock
		given          *UpdateContext
		expectedRunLog []string
	}{
		{
			description: "it should activate chat if it's not active and call next handler",
			passStorage: chatActivatorPassStorageMock{passExistsRes: true},
			chatStorage: chatActivatorChatStorageMock{},
			given:       &UpdateContext{Update: update},
			expectedRunLog: []string{
				"PassExists: pass",
				"ActivateChat: 123",
				"DeletePass: pass",
				"Handle: " + testUpdateToStr(&UpdateContext{Update: update, IsChatActive: true}),
			},
		},
		{
			description: "it should do nothing if chat is already active",
			passStorage: chatActivatorPassStorageMock{passExistsRes: true},
			chatStorage: chatActivatorChatStorageMock{},
			given:       &UpdateContext{Update: update, IsChatActive: true},
			expectedRunLog: []string{
				"Handle: " + testUpdateToStr(&UpdateContext{Update: update, IsChatActive: true}),
			},
		},
		{
			description: "it should exit and change nothing if PassExists returns error",
			passStorage: chatActivatorPassStorageMock{passExistsErr: testError},
			chatStorage: chatActivatorChatStorageMock{},
			given:       &UpdateContext{Update: update},
			expectedRunLog: []string{
				"PassExists: pass",
			},
		},
		{
			description: "it should exit and change nothing if pass doesn't exist",
			passStorage: chatActivatorPassStorageMock{},
			chatStorage: chatActivatorChatStorageMock{},
			given:       &UpdateContext{Update: update},
			expectedRunLog: []string{
				"PassExists: pass",
			},
		},
		{
			description: "it should exit and change nothing if ActivateChat returns error",
			passStorage: chatActivatorPassStorageMock{passExistsRes: true},
			chatStorage: chatActivatorChatStorageMock{activateChatErr: testError},
			given:       &UpdateContext{Update: update},
			expectedRunLog: []string{
				"PassExists: pass",
				"ActivateChat: 123",
			},
		},
		{
			description: "it should exit and change nothing if DeletePass returns error",
			passStorage: chatActivatorPassStorageMock{passExistsRes: true, deletePassErr: testError},
			chatStorage: chatActivatorChatStorageMock{},
			given:       &UpdateContext{Update: update},
			expectedRunLog: []string{
				"PassExists: pass",
				"ActivateChat: 123",
				"DeletePass: pass",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			runLog := make([]string, 0)
			tc.passStorage.runLog = &runLog
			tc.chatStorage.runLog = &runLog
			h := NewChatActivator(&tc.passStorage, &tc.chatStorage)
			h.SetNext(&mockHandler{runLog: &runLog})

			h.Handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}
