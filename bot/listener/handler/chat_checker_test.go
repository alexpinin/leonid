package handler

import (
	"errors"
	"leonid/testutil"
	"testing"

	"github.com/go-telegram/bot/models"
)

type chatCheckerStorageMock struct {
	chatExistsRes bool
	chatExistsErr error
}

func (m *chatCheckerStorageMock) ChatExists(int64) (bool, error) {
	return m.chatExistsRes, m.chatExistsErr
}

func Test_ChatChecker_Handle(t *testing.T) {
	testCases := []struct {
		description string
		storage     chatCheckerStorage
		given       *UpdateContext
		expected    *UpdateContext
	}{
		{
			description: "should set chat as active and call next handler if chat exists",
			storage: &chatCheckerStorageMock{
				chatExistsRes: true,
				chatExistsErr: nil,
			},
			given: &UpdateContext{
				Update: &models.Update{
					Message: &models.Message{},
				},
			},
			expected: &UpdateContext{
				Update: &models.Update{
					Message: &models.Message{},
				},
				isChatActive: true,
			},
		},
		{
			description: "should stop and exit if ChatExists returns error",
			storage: &chatCheckerStorageMock{
				chatExistsErr: errors.New("test"),
			},
			given: &UpdateContext{
				Update: &models.Update{
					Message: &models.Message{},
				},
			},
			expected: nil,
		},
		{
			description: "should set chat as not active and call next handler if chat doesn't exist",
			storage: &chatCheckerStorageMock{
				chatExistsRes: false,
				chatExistsErr: nil,
			},
			given: &UpdateContext{
				Update: &models.Update{
					Message: &models.Message{},
				},
			},
			expected: &UpdateContext{
				Update: &models.Update{
					Message: &models.Message{},
				},
				isChatActive: false,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			m := &mockHandler{}
			h := NewChatChecker(tc.storage)
			h.SetNext(m)

			h.Handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expected, m.updateContext)
		})
	}
}
