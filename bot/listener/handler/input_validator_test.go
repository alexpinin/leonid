package handler

import (
	"leonid/testutil"
	"testing"

	"github.com/go-telegram/bot/models"
)

func Test_InputValidator_Handle(t *testing.T) {
	testCases := []struct {
		description    string
		given          *UpdateContext
		expectedRunLog []string
	}{
		{
			description: "should call next handler if update is valid",
			given:       &UpdateContext{Update: &models.Update{Message: &models.Message{}}},
			expectedRunLog: []string{
				"Handle: " + testUpdateToStr(&UpdateContext{Update: &models.Update{Message: &models.Message{}}, IsInputValid: true}),
			},
		},
		{
			description:    "should not call next handler if update is nil",
			given:          nil,
			expectedRunLog: []string{},
		},
		{
			description:    "should not call next handler if update message is nil",
			given:          &UpdateContext{Update: &models.Update{}},
			expectedRunLog: []string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			runLog := make([]string, 0)
			h := &InputValidator{}
			h.SetNext(&mockHandler{runLog: &runLog})

			h.Handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}
