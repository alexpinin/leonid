package handler

import (
	"testing"

	"github.com/go-telegram/bot/models"

	"leonid/src/internal/testutil"
)

func TestInputGuardHandle(t *testing.T) {
	testCases := []struct {
		description        string
		givenUpdate        *models.Update
		expectedErr        error
		expectedNextCalled int
	}{
		{
			description: "should call next handler if update is valid",
			givenUpdate: &models.Update{
				Message: &models.Message{
					Chat: models.Chat{ID: 123},
				},
			},
			expectedErr:        nil,
			expectedNextCalled: 1,
		},
		{
			description:        "should not call next handler and exit if update is nil",
			givenUpdate:        nil,
			expectedErr:        nil,
			expectedNextCalled: 0,
		},
		{
			description:        "should not call next handler and exit if update message is nil",
			givenUpdate:        &models.Update{},
			expectedErr:        nil,
			expectedNextCalled: 0,
		},
		{
			description:        "should not call next handler and exit if chat ID is invalid",
			givenUpdate:        &models.Update{Message: &models.Message{}},
			expectedErr:        nil,
			expectedNextCalled: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			next := &mockHandler{}
			sut := &inputGuard{}
			sut.setNext(next)

			err := sut.handle(nil, nil, tc.givenUpdate, nil)

			testutil.Equal(t, tc.expectedErr, err)
			testutil.Equal(t, tc.expectedNextCalled, next.handleCount)
		})
	}
}
