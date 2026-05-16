package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	audioMaxDurationSec    = 10
	audioMimeType          = "audio/ogg"
	audioLanguage          = "ru"
	audioProcessingTimeout = time.Duration(10) * time.Second
)

type AudioService struct {
	transcribeURL string
}

func NewAudioService(transcribeURL string) *AudioService {
	return &AudioService{
		transcribeURL: transcribeURL,
	}
}

type audioReq struct {
	URL  string `json:"url"`
	Lang string `json:"lang"`
}

type audioRes struct {
	Text string `json:"text"`
}

func (s *AudioService) TranscribeAudio(ctx context.Context, b *bot.Bot, voice *models.Voice) (string, error) {
	if voice == nil {
		return "", fmt.Errorf("AudioService.TranscribeAudio: voice is nil")
	}

	if voice.Duration > audioMaxDurationSec {
		return "", fmt.Errorf("AudioService.TranscribeAudio: voice duration is too long: %d", voice.Duration)
	}

	if voice.MimeType != audioMimeType {
		return "", fmt.Errorf("AudioService.TranscribeAudio: unsupported audio format: %s", voice.MimeType)
	}

	params := &bot.GetFileParams{FileID: voice.FileID}
	file, err := b.GetFile(ctx, params)
	if err != nil {
		return "", fmt.Errorf("AudioService.TranscribeAudio: cannot get audio file: %w", err)
	}

	downloadLink := b.FileDownloadLink(file)

	body, err := json.Marshal(audioReq{
		URL:  downloadLink,
		Lang: audioLanguage,
	})
	if err != nil {
		return "", fmt.Errorf("AudioService.TranscribeAudio: cannot marshal request body: %w", err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, audioProcessingTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, s.transcribeURL, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("AudioService.TranscribeAudio: cannot create a context: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("AudioService.TranscribeAudio: transcribing service error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AudioService.TranscribeAudio: unexpected status code: %d", resp.StatusCode)
	}

	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("AudioService.TranscribeAudio: cannot read audio file: %w", err)
	}

	var res audioRes
	err = json.Unmarshal(bb, &res)
	if err != nil {
		return "", fmt.Errorf("AudioService.TranscribeAudio: cannot unmarshal response body: %w", err)
	}

	if res.Text == "" {
		return "", errors.New("AudioService.TranscribeAudio: transcribing service failed to process")
	}

	return res.Text, nil
}
