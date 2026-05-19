package handler

import (
	"context"
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
		description        string
		quotaManager       *mockQuotaManager
		givenUpdate        *models.Update
		expectedErr        error
		expectedNextCalled int
	}{
		{
			description:        "should call quotaManager and next handler",
			quotaManager:       &mockQuotaManager{},
			givenUpdate:        update,
			expectedErr:        nil,
			expectedNextCalled: 1,
		},
		{
			description:        "should call quotaManager and exit if it returns error",
			quotaManager:       &mockQuotaManager{useChatQuotaErr: testutil.TestError},
			givenUpdate:        update,
			expectedErr:        testutil.TestError,
			expectedNextCalled: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			next := &mockHandler{}
			sut := newQuotaGuard(tc.quotaManager)
			sut.setNext(next)

			err := sut.handle(nil, nil, tc.givenUpdate, nil)

			testutil.ErrorIs(t, tc.expectedErr, err)
			testutil.Equal(t, tc.expectedNextCalled, next.handleCount)
		})
	}
}

type mockQuotaManager struct {
	useChatQuotaErr error
}

func (m *mockQuotaManager) UseChatQuota(context.Context, int64) error {
	return m.useChatQuotaErr
}
