package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-telegram/bot/models"

	"leonid/src/internal/testutil"
)

func TestCallGuardHandle(t *testing.T) {
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
		nicknameProvider mockNicknameProvider
		given            *UpdateContext
		expectedRunLog   []string
	}{
		{
			description:      "should call next handler if bot is called by a nickname ignoring case",
			nicknameProvider: mockNicknameProvider{listNicknamesRes: []string{"Bot"}},
			given:            &UpdateContext{Update: update},
			expectedRunLog: []string{
				"ListNicknames: 123",
				"handle: " + testUpdateToStr(&UpdateContext{Update: update}),
			},
		},
		{
			description:      "should not call next handler and exit if bot is not called by a nickname",
			nicknameProvider: mockNicknameProvider{listNicknamesRes: []string{"bot2"}},
			given:            &UpdateContext{Update: update},
			expectedRunLog: []string{
				"ListNicknames: 123",
			},
		},
		{
			description:      "should not call next handler and exit if there are no nicknames present",
			nicknameProvider: mockNicknameProvider{listNicknamesRes: nil},
			given:            &UpdateContext{Update: update},
			expectedRunLog: []string{
				"ListNicknames: 123",
			},
		},
		{
			description:      "should ignore empty nicknames",
			nicknameProvider: mockNicknameProvider{listNicknamesRes: []string{""}},
			given:            &UpdateContext{Update: update},
			expectedRunLog: []string{
				"ListNicknames: 123",
			},
		},
		{
			description:      "should call next handler if bot is called by in reply",
			nicknameProvider: mockNicknameProvider{listNicknamesRes: []string{"bot"}},
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
			nicknameProvider: mockNicknameProvider{listNicknamesRes: []string{"bot"}},
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
			nicknameProvider: mockNicknameProvider{listNicknamesRes: []string{"bot"}},
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
			nicknameProvider: mockNicknameProvider{listNicknamesRes: []string{"bot"}},
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

			_ = sut.handle(nil, nil, tc.given)

			testutil.Equal(t, tc.expectedRunLog, runLog)
		})
	}
}

type mockNicknameProvider struct {
	runLog           *[]string
	listNicknamesRes []string
	listNicknamesErr error
}

func (m *mockNicknameProvider) ListNicknames(_ context.Context, chatID int64) ([]string, error) {
	*m.runLog = append(*m.runLog, fmt.Sprintf("ListNicknames: %d", chatID))
	return m.listNicknamesRes, m.listNicknamesErr
}
