package handler

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type audioReader struct {
	basicHandler
	enabled      bool
	transcriptor audioTranscriptor
}

func newAudioReader(
	enabled bool,
	t audioTranscriptor,
) *audioReader {
	return &audioReader{
		enabled:      enabled,
		transcriptor: t,
	}
}

type audioTranscriptor interface {
	TranscribeAudio(ctx context.Context, b *bot.Bot, voice *models.Voice) (string, error)
}

func (h *audioReader) handle(c context.Context, b *bot.Bot, u *UpdateContext) error {
	if h.enabled && u.Message.Voice != nil {
		message, err := h.transcriptor.TranscribeAudio(c, b, u.Message.Voice)
		if err != nil {
			return fmt.Errorf("audioReader.handle: %w", err)
		}
		u.Message.Text = message
	}
	return h.nextHandle(c, b, u)
}
