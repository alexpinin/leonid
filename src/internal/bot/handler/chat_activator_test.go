package handler

import (
	"context"
	"testing"

	"github.com/go-telegram/bot/models"

	"leonid/src/internal/testutil"
)

func TestChatActivatorHandle(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Text: "pass",
			Chat: models.Chat{
				ID: 123,
			},
		},
	}
	testCases := []struct {
		description            string
		chatActivator          mockChatActivator
		givenUpdate            *models.Update
		givenState             *UpdateState
		expectedContext        *UpdateState
		expectedErr            error
		expectedActivateCalled int
		expectedNextCalled     int
	}{
		{
			description:            "should do nothing and call next handler if chat is active",
			chatActivator:          mockChatActivator{},
			givenUpdate:            update,
			givenState:             &UpdateState{IsChatActive: true},
			expectedContext:        &UpdateState{IsChatActive: true},
			expectedErr:            nil,
			expectedActivateCalled: 0,
			expectedNextCalled:     1,
		},
		{
			description:            "should call chActivator and next handler",
			chatActivator:          mockChatActivator{activateRes: true},
			givenUpdate:            update,
			givenState:             &UpdateState{},
			expectedContext:        &UpdateState{IsPassActive: true},
			expectedErr:            nil,
			expectedActivateCalled: 1,
			expectedNextCalled:     1,
		},
		{
			description:            "should call chActivator and exit if it returns error",
			chatActivator:          mockChatActivator{activateErr: testutil.TestError},
			givenUpdate:            update,
			givenState:             &UpdateState{},
			expectedContext:        &UpdateState{},
			expectedErr:            testutil.TestError,
			expectedActivateCalled: 1,
			expectedNextCalled:     0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			next := &mockHandler{}
			sut := newChatActivator(&tc.chatActivator)
			sut.setNext(next)

			err := sut.handle(nil, nil, tc.givenUpdate, tc.givenState)

			testutil.ErrorIs(t, tc.expectedErr, err)
			testutil.Equal(t, tc.expectedContext, tc.givenState)
			testutil.Equal(t, tc.expectedActivateCalled, tc.chatActivator.activateCalled)
			testutil.Equal(t, tc.expectedNextCalled, next.handleCount)
		})
	}
}

type mockChatActivator struct {
	activateCalled int
	activateRes    bool
	activateErr    error
}

func (m *mockChatActivator) Activate(context.Context, string, int64) (bool, error) {
	m.activateCalled++
	return m.activateRes, m.activateErr
}
