package handler

import (
	"context"
	"testing"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"leonid/src/internal/testutil"
)

func TestCallGuardHandle(t *testing.T) {
	tBot := &bot.Bot{}
	tBot.SetToken("456")
	update := models.Update{
		Message: &models.Message{
			Text: "Hello, Bot",
			Chat: models.Chat{
				ID: 123,
			},
		},
	}
	testCases := []struct {
		description        string
		nicknameProvider   mockNicknameProvider
		givenUpdate        *models.Update
		expectedErr        error
		expectedNextCalled int
	}{
		{
			description:        "should call nicknameProvider and next handler if bot is called by a nickname ignoring case",
			nicknameProvider:   mockNicknameProvider{listNicknamesRes: []string{"Bot"}},
			givenUpdate:        &update,
			expectedNextCalled: 1,
			expectedErr:        nil,
		},
		{
			description:        "should call nicknameProvider and and exit if bot is not called by a nickname",
			nicknameProvider:   mockNicknameProvider{listNicknamesRes: []string{"bot2"}},
			givenUpdate:        &update,
			expectedNextCalled: 0,
			expectedErr:        nil,
		},
		{
			description:        "should call nicknameProvider and and exit if there are no nicknames present",
			nicknameProvider:   mockNicknameProvider{listNicknamesRes: nil},
			givenUpdate:        &update,
			expectedNextCalled: 0,
			expectedErr:        nil,
		},
		{
			description:        "should ignore empty nicknames",
			nicknameProvider:   mockNicknameProvider{listNicknamesRes: []string{""}},
			givenUpdate:        &update,
			expectedNextCalled: 0,
			expectedErr:        nil,
		},
		{
			description:      "should call nicknameProvider and next handler if bot is called by in reply",
			nicknameProvider: mockNicknameProvider{listNicknamesRes: []string{"bot"}},
			givenUpdate: &models.Update{
				Message: &models.Message{
					Text: "Hello",
					Chat: models.Chat{
						ID: 123,
					},
					ReplyToMessage: &models.Message{
						From: &models.User{
							ID: 456,
						},
					},
				},
			},
			expectedNextCalled: 1,
			expectedErr:        nil,
		},
		{
			description:      "should call nicknameProvider and and exit if From is nil",
			nicknameProvider: mockNicknameProvider{listNicknamesRes: []string{"bot"}},
			givenUpdate: &models.Update{
				Message: &models.Message{
					Text: "Hello",
					Chat: models.Chat{
						ID: 123,
					},
					ReplyToMessage: &models.Message{},
				},
			},
			expectedNextCalled: 0,
			expectedErr:        nil,
		},
		{
			description:      "should call nicknameProvider and and exit if ReplyToMessage is nil",
			nicknameProvider: mockNicknameProvider{listNicknamesRes: []string{"bot"}},
			givenUpdate: &models.Update{
				Message: &models.Message{
					Text: "Hello",
					Chat: models.Chat{
						ID: 123,
					},
				},
			},
			expectedNextCalled: 0,
			expectedErr:        nil,
		},
		{
			description:      "should call nicknameProvider and exit if it returns error",
			nicknameProvider: mockNicknameProvider{listNicknamesErr: testutil.TestError},
			givenUpdate: &models.Update{
				Message: &models.Message{
					Text: "Hello",
					Chat: models.Chat{
						ID: 123,
					},
				},
			},
			expectedNextCalled: 0,
			expectedErr:        testutil.TestError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			next := &mockHandler{}
			sut := newCallGuard(&tc.nicknameProvider)
			sut.setNext(next)

			err := sut.handle(nil, tBot, tc.givenUpdate, nil)

			testutil.ErrorIs(t, tc.expectedErr, err)
			testutil.Equal(t, tc.expectedNextCalled, next.handleCount)
		})
	}
}

type mockNicknameProvider struct {
	listNicknamesRes []string
	listNicknamesErr error
}

func (m *mockNicknameProvider) ListNicknames(context.Context, int64) ([]string, error) {
	return m.listNicknamesRes, m.listNicknamesErr
}
