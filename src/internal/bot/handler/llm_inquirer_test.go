package handler

import (
	"context"
	"testing"

	"github.com/go-telegram/bot/models"

	"leonid/src/internal/testutil"
)

func TestLLMInquirerHandle(t *testing.T) {
	update := &models.Update{
		Message: &models.Message{
			Text: "message",
			Chat: models.Chat{ID: 123},
		},
	}
	testCases := []struct {
		description        string
		messageSender      mockMessageSender
		givenUpdate        *models.Update
		givenState         *UpdateState
		expectedState      *UpdateState
		expectedErr        error
		expectedNextCalled int
	}{
		{
			description:        "should call inquirer and next handler",
			messageSender:      mockMessageSender{inquireLLMRes: "hello"},
			givenUpdate:        update,
			givenState:         &UpdateState{},
			expectedState:      &UpdateState{LLMReply: "hello"},
			expectedErr:        nil,
			expectedNextCalled: 1,
		},
		{
			description:        "should call inquirer and exit if it returns error",
			messageSender:      mockMessageSender{inquireLLMErr: testutil.TestError},
			givenUpdate:        update,
			givenState:         &UpdateState{},
			expectedState:      &UpdateState{},
			expectedErr:        testutil.TestError,
			expectedNextCalled: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			next := &mockHandler{}
			sut := newLLMInquirer(&tc.messageSender)
			sut.setNext(next)

			err := sut.handle(nil, nil, tc.givenUpdate, tc.givenState)

			testutil.ErrorIs(t, tc.expectedErr, err)
			testutil.Equal(t, tc.expectedState, tc.givenState)
			testutil.Equal(t, tc.expectedNextCalled, next.handleCount)
		})
	}
}

type mockMessageSender struct {
	inquireLLMRes string
	inquireLLMErr error
}

func (m *mockMessageSender) InquireLLM(context.Context, int64, string) (string, error) {
	return m.inquireLLMRes, m.inquireLLMErr
}
