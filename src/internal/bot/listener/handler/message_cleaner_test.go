package handler

import (
	"leonid/testutil"
	"testing"

	"github.com/go-telegram/bot/models"
)

func Test_MessageCleaner_Handle(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Text: "Hello, Bot",
			Chat: models.Chat{ID: 123},
		},
	}
	testCases := []struct {
		description      string
		nicknameProvider nicknameProviderMock
		given            *UpdateContext
		expectedRunLog   []string
	}{
		{
			description:      "should delete nicknames from messages, lowercase and call next handler",
			nicknameProvider: nicknameProviderMock{result: []string{"bot"}},
			given:            &UpdateContext{Update: update},
			expectedRunLog: []string{
				"ListNicknames: 123",
				"Handle: " + testUpdateToStr(&UpdateContext{Update: &models.Update{
					Message: &models.Message{
						Text: "hello, ",
						Chat: models.Chat{ID: 123},
					},
				}}),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			runLog := make([]string, 0)
			tc.nicknameProvider.runLog = &runLog
			sut := NewMessageCleaner(&tc.nicknameProvider)
			sut.SetNext(&mockHandler{runLog: &runLog})

			sut.Handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}
