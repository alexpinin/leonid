package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-telegram/bot/models"

	"leonid/src/internal/testutil"
)

type nicknameProviderMock struct {
	runLog *[]string
	result []string
}

func (m *nicknameProviderMock) ListNicknames(_ context.Context, chatID int64) []string {
	*m.runLog = append(*m.runLog, fmt.Sprintf("ListNicknames: %d", chatID))
	return m.result
}

func Test_callGuard_handle(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Text: "Hello, Bot",
			Chat: models.Chat{
				ID: 123,
			},
		},
	}
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
				"ListNicknames: 123",
				"handle: " + testUpdateToStr(&UpdateContext{Update: update}),
			},
		},
		{
			description:      "should not call next handler and exit if bot is not called by a nickname",
			nicknameProvider: nicknameProviderMock{result: []string{"bot2"}},
			given:            &UpdateContext{Update: update},
			expectedRunLog: []string{
				"ListNicknames: 123",
			},
		},
		{
			description:      "should not call next handler and exit if there are no nicknames present",
			nicknameProvider: nicknameProviderMock{result: nil},
			given:            &UpdateContext{Update: update},
			expectedRunLog: []string{
				"ListNicknames: 123",
			},
		},
		{
			description:      "should ignore empty nicknames",
			nicknameProvider: nicknameProviderMock{result: []string{""}},
			given:            &UpdateContext{Update: update},
			expectedRunLog: []string{
				"ListNicknames: 123",
			},
		},
		{
			description:      "should call next handler if bot is called by in reply",
			nicknameProvider: nicknameProviderMock{result: []string{"bot"}},
			given: &UpdateContext{Update: &models.Update{
				Message: &models.Message{
					Text: "Hello",
					Chat: models.Chat{
						ID: 123,
					},
					ReplyToMessage: &models.Message{
						From: &models.User{
							FirstName: "Bot",
						},
					},
				},
			}},
			expectedRunLog: []string{
				"ListNicknames: 123",
				"handle: " + testUpdateToStr(&UpdateContext{Update: &models.Update{
					Message: &models.Message{
						Text: "Hello",
						Chat: models.Chat{
							ID: 123,
						},
						ReplyToMessage: &models.Message{
							From: &models.User{
								FirstName: "Bot",
							},
						},
					},
				}}),
			},
		},
		{
			description:      "should not call next handler and exit if replay name doesn't match",
			nicknameProvider: nicknameProviderMock{result: []string{"bot"}},
			given: &UpdateContext{Update: &models.Update{
				Message: &models.Message{
					Text: "Hello",
					Chat: models.Chat{
						ID: 123,
					},
					ReplyToMessage: &models.Message{
						From: &models.User{
							FirstName: "Bot1",
						},
					},
				},
			}},
			expectedRunLog: []string{
				"ListNicknames: 123",
			},
		},
		{
			description:      "should not call next handler and exit if From is nil",
			nicknameProvider: nicknameProviderMock{result: []string{"bot"}},
			given: &UpdateContext{Update: &models.Update{
				Message: &models.Message{
					Text: "Hello",
					Chat: models.Chat{
						ID: 123,
					},
					ReplyToMessage: &models.Message{},
				},
			}},
			expectedRunLog: []string{
				"ListNicknames: 123",
			},
		},
		{
			description:      "should not call next handler and exit if ReplyToMessage is nil",
			nicknameProvider: nicknameProviderMock{result: []string{"bot"}},
			given: &UpdateContext{Update: &models.Update{
				Message: &models.Message{
					Text: "Hello",
					Chat: models.Chat{
						ID: 123,
					},
				},
			}},
			expectedRunLog: []string{
				"ListNicknames: 123",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			runLog := make([]string, 0)
			tc.nicknameProvider.runLog = &runLog
			sut := newCallGuard(&tc.nicknameProvider)
			sut.setNext(&mockHandler{runLog: &runLog})

			sut.handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}
