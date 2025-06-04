package handler

import (
	"context"
	"leonid/testutil"
	"testing"

	"github.com/go-telegram/bot/models"
)

type nicknameProviderMock struct {
	runLog *[]string
	result []string
}

func (p *nicknameProviderMock) ListNicknames(ctx context.Context, chatID int64) []string {
	return p.result
}

func Test_CallGuard_Handle(t *testing.T) {
	update := &models.Update{Message: &models.Message{Text: "Hello, Bot"}}
	testCases := []struct {
		description      string
		nicknameProvider nicknameProviderMock
		given            *UpdateContext
		expectedRunLog   []string
	}{
		{
			description:      "should call next handler if bot is called by a nickname ignoring case",
			nicknameProvider: nicknameProviderMock{result: []string{"bot"}},
			given:            &UpdateContext{Update: update},
			expectedRunLog: []string{
				"Handle: " + testUpdateToStr(&UpdateContext{Update: update}),
			},
		},
		{
			description:      "should not call next handler and exit if bot is not called by a nickname",
			nicknameProvider: nicknameProviderMock{result: []string{"bot2"}},
			given:            &UpdateContext{Update: update},
			expectedRunLog:   []string{},
		},
		{
			description:      "should not call next handler and exit if there are no nicknames present",
			nicknameProvider: nicknameProviderMock{result: nil},
			given:            &UpdateContext{Update: update},
			expectedRunLog:   []string{},
		},
		{
			description:      "should ignore empty nicknames",
			nicknameProvider: nicknameProviderMock{result: []string{""}},
			given:            &UpdateContext{Update: update},
			expectedRunLog:   []string{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			runLog := make([]string, 0)
			tc.nicknameProvider.runLog = &runLog
			sut := NewCallGuard(&tc.nicknameProvider)
			sut.SetNext(&mockHandler{runLog: &runLog})

			sut.Handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}
