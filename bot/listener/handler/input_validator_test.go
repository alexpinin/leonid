package handler

import (
	"leonid/testutil"
	"testing"

	"github.com/go-telegram/bot/models"
)

func Test_InputValidator_Handle(t *testing.T) {
	testCases := []struct {
		description string
		given       *UpdateContext
		expected    *UpdateContext
	}{
		{
			description: "should call next handler if update is valid",
			given: &UpdateContext{
				Update: &models.Update{
					Message: &models.Message{},
				},
			},
			expected: &UpdateContext{
				Update: &models.Update{
					Message: &models.Message{},
				},
				IsInputValid: true,
			},
		},
		{
			description: "should not call next handler if update is nil",
			given:       nil,
			expected:    nil,
		},
		{
			description: "should not call next handler if update message is nil",
			given: &UpdateContext{
				Update: &models.Update{},
			},
			expected: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			m := &mockHandler{}
			h := &InputValidator{}
			h.SetNext(m)

			h.Handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expected, m.updateContext)
		})
	}
}
