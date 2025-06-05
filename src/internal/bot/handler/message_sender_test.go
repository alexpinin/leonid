package handler

import (
	"context"
	"fmt"
	"leonid/testutil"
	"testing"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type messageSenderMock struct {
	runLog *[]string
}

func (m *messageSenderMock) SendMessage(_ context.Context, _ *bot.Bot, chatID int64, message string) {
	*m.runLog = append(*m.runLog, fmt.Sprintf("SendMessage: %d, %s", chatID, message))
}

func Test_messageSender_handle(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Text: "message",
			Chat: models.Chat{ID: 123},
		},
	}
	testCases := []struct {
		description    string
		messageSender  messageSenderMock
		given          *UpdateContext
		expectedRunLog []string
	}{
		{
			description:   "should send message and call next handler",
			messageSender: messageSenderMock{},
			given:         &UpdateContext{Update: update},
			expectedRunLog: []string{
				"SendMessage: 123, message",
				"handle: " + testUpdateToStr(&UpdateContext{Update: update}),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			runLog := make([]string, 0)
			tc.messageSender.runLog = &runLog
			sut := newMessageSender(&tc.messageSender)
			sut.setNext(&mockHandler{runLog: &runLog})

			sut.handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}
