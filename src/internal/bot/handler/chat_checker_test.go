package handler

import (
	"context"
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
		description        string
		storage            mockChatCheckerStorage
		givenUpdate        *models.Update
		givenContext       *UpdateContext
		expectedContext    *UpdateContext
		expectedErr        error
		expectedNextCalled int
	}{
		{
			description:        "should call chChecker and next handler",
			storage:            mockChatCheckerStorage{isChatActiveRes: true},
			givenUpdate:        update,
			givenContext:       &UpdateContext{},
			expectedContext:    &UpdateContext{IsChatActive: true},
			expectedErr:        nil,
			expectedNextCalled: 1,
		},
		{
			description:        "should call chChecker and exit if it returns error",
			storage:            mockChatCheckerStorage{isChatActiveErr: testutil.TestError},
			givenUpdate:        update,
			givenContext:       &UpdateContext{},
			expectedContext:    &UpdateContext{},
			expectedErr:        testutil.TestError,
			expectedNextCalled: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			next := mockHandler{}
			sut := newChatChecker(&tc.storage)
			sut.setNext(&next)

			err := sut.handle(nil, nil, tc.givenUpdate, tc.givenContext)

			testutil.ErrorIs(t, tc.expectedErr, err)
			testutil.Equal(t, tc.expectedNextCalled, next.handleCount)
		})
	}
}

type mockChatCheckerStorage struct {
	isChatActiveRes bool
	isChatActiveErr error
}

func (m *mockChatCheckerStorage) IsChatActive(context.Context, int64) (bool, error) {
	return m.isChatActiveRes, m.isChatActiveErr
}
