package handler

import (
	"testing"

	"github.com/go-telegram/bot/models"

	"leonid/src/internal/testutil"
)

func Test_inputGuard_handle(t *testing.T) {
	testCases := []struct {
		description    string
		given          *UpdateContext
		expectedRunLog []string
	}{
		{
			description: "should call next handler if update is valid",
			given: &UpdateContext{
				Update: &models.Update{
					Message: &models.Message{
						Chat: models.Chat{ID: 123},
					},
				},
			},
			expectedRunLog: []string{
				"handle: " + testUpdateToStr(&UpdateContext{
					Update: &models.Update{
						Message: &models.Message{
							Chat: models.Chat{ID: 123},
						},
					},
				}),
			},
		},
		{
			description:    "should not call next handler and exit if update is nil",
			given:          nil,
			expectedRunLog: []string{},
		},
		{
			description:    "should not call next handler and exit if update message is nil",
			given:          &UpdateContext{Update: &models.Update{}},
			expectedRunLog: []string{},
		},
		{
			description:    "should not call next handler and exit if chat ID is invalid",
			given:          &UpdateContext{Update: &models.Update{Message: &models.Message{}}},
			expectedRunLog: []string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			runLog := make([]string, 0)
			sut := &inputGuard{}
			sut.setNext(&mockHandler{runLog: &runLog})

			sut.handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}
