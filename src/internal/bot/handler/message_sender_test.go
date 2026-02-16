package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"leonid/src/internal/testutil"
)

func TestMessageSenderHandle(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Text: "message",
			Chat: models.Chat{ID: 123},
		},
	}
	testCases := []struct {
		description    string
		messageSender  mockMessageSender
		given          *UpdateContext
		expectedRunLog []string
	}{
		{
			description:   "should send message and call next handler",
			messageSender: mockMessageSender{},
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

			_ = sut.handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}

type mockMessageSender struct {
	runLog *[]string
}

func (m *mockMessageSender) SendMessage(_ context.Context, _ *bot.Bot, chatID int64, message string) error {
	*m.runLog = append(*m.runLog, fmt.Sprintf("SendMessage: %d, %s", chatID, message))
	return nil
}
