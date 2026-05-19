package handler

import (
	"testing"

	"leonid/src/internal/testutil"
)

func TestAuthGuardHandle(t *testing.T) {
	testCases := []struct {
		description        string
		givenState         *UpdateState
		expectedNextCalled int
	}{
		{
			description:        "should call next handler if chat is active",
			givenState:         &UpdateState{IsChatActive: true},
			expectedNextCalled: 1,
		},
		{
			description:        "should call next handler if pass phrase is active",
			givenState:         &UpdateState{IsPassActive: true},
			expectedNextCalled: 1,
		},
		{
			description:        "should exit and not call next handler if neither chat nor pass is active",
			givenState:         &UpdateState{},
			expectedNextCalled: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			next := &mockHandler{}
			sut := newAuthGuard()
			sut.setNext(next)

			_ = sut.handle(nil, nil, nil, tc.givenState)

			testutil.Equal(t, tc.expectedNextCalled, next.handleCount)
		})
	}
}
