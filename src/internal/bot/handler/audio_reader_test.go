package handler

import (
	"context"
	"testing"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"leonid/src/internal/testutil"
)

func TestAudioReaderHandle(t *testing.T) {
	testCases := []struct {
		description        string
		transcriptor       audioTranscriptor
		enabled            bool
		givenUpdate        *models.Update
		expectedUpdate     *models.Update
		expectedErr        error
		expectedNextCalled int
	}{
		{
			description:        "should call audioTranscriptor and next handler",
			transcriptor:       &mockAudioTranscriptor{transcribeAudioRes: "transcription"},
			enabled:            true,
			givenUpdate:        &models.Update{Message: &models.Message{Voice: &models.Voice{}}},
			expectedUpdate:     &models.Update{Message: &models.Message{Voice: &models.Voice{}, Text: "transcription"}},
			expectedErr:        nil,
			expectedNextCalled: 1,
		},
		{
			description:        "should not call audioTranscriptor and call next handler if disabled",
			transcriptor:       &mockAudioTranscriptor{transcribeAudioRes: "transcription"},
			enabled:            false,
			givenUpdate:        &models.Update{Message: &models.Message{Voice: &models.Voice{}}},
			expectedUpdate:     &models.Update{Message: &models.Message{Voice: &models.Voice{}}},
			expectedErr:        nil,
			expectedNextCalled: 1,
		},
		{
			description:        "should not call audioTranscriptor and call next handler if there is no voice",
			transcriptor:       &mockAudioTranscriptor{transcribeAudioRes: "transcription"},
			enabled:            true,
			givenUpdate:        &models.Update{Message: &models.Message{}},
			expectedUpdate:     &models.Update{Message: &models.Message{}},
			expectedErr:        nil,
			expectedNextCalled: 1,
		},
		{
			description:        "should call audioTranscriptor and exit if it returns error",
			transcriptor:       &mockAudioTranscriptor{transcribeAudioErr: testutil.TestError},
			enabled:            true,
			givenUpdate:        &models.Update{Message: &models.Message{Voice: &models.Voice{}}},
			expectedUpdate:     &models.Update{Message: &models.Message{Voice: &models.Voice{}}},
			expectedErr:        testutil.TestError,
			expectedNextCalled: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			next := &mockHandler{}
			sut := newAudioReader(tc.enabled, tc.transcriptor)
			sut.setNext(next)

			err := sut.handle(nil, nil, tc.givenUpdate, nil)

			testutil.ErrorIs(t, tc.expectedErr, err)
			testutil.Equal(t, tc.expectedUpdate, tc.givenUpdate)
			testutil.Equal(t, tc.expectedNextCalled, next.handleCount)
		})
	}
}

type mockAudioTranscriptor struct {
	transcribeAudioRes string
	transcribeAudioErr error
}

func (t *mockAudioTranscriptor) TranscribeAudio(context.Context, *bot.Bot, *models.Voice) (string, error) {
	return t.transcribeAudioRes, t.transcribeAudioErr
}
