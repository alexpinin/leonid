package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	audioMaxDurationSec = 10
	audioMimeType       = "audio/ogg"
	audioLanguage       = "ru"
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

	req := audioReq{
		URL:  downloadLink,
		Lang: audioLanguage,
	}
	marshal, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("AudioService.TranscribeAudio: cannot marshal request body: %w", err)
	}

	resp, err := http.Post(s.transcribeURL, "application/json", bytes.NewReader(marshal))
	if err != nil {
		return "", fmt.Errorf("AudioService.TranscribeAudio: cannot download audio file: %w", err)
	}
	defer resp.Body.Close()

	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("AudioService.TranscribeAudio: cannot read audio file: %w", err)
	}

	var res audioRes
	err = json.Unmarshal(bb, &res)
	if err != nil {
		return "", fmt.Errorf("AudioService.TranscribeAudio: cannot unmarshal response body: %w", err)
	}

	return res.Text, err
}
