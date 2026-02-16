package handler

import (
	"testing"

	"leonid/src/internal/testutil"
)

func TestAuthGuardHandle(t *testing.T) {
	testCases := []struct {
		description    string
		given          *UpdateContext
		expectedRunLog []string
	}{
		{
			description: "should call next handler if chat is active",
			given:       &UpdateContext{IsChatActive: true},
			expectedRunLog: []string{
				"handle: " + testUpdateToStr(&UpdateContext{IsChatActive: true}),
			},
		},
		{
			description: "should call next handler if pass phrase is active",
			given:       &UpdateContext{IsPassActive: true},
			expectedRunLog: []string{
				"handle: " + testUpdateToStr(&UpdateContext{IsPassActive: true}),
			},
		},
		{
			description:    "should exit and not call next handler if neither chat nor pass is active",
			given:          &UpdateContext{},
			expectedRunLog: []string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			runLog := make([]string, 0)
			sut := newAuthGuard()
			sut.setNext(&mockHandler{runLog: &runLog})

			_ = sut.handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}
