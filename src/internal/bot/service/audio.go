package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-telegram/bot"
)

type AudioService struct {
}

func NewAudioService() *AudioService {
	return &AudioService{}
}

type audioReq struct {
	URL  string `json:"url"`
	Lang string `json:"lang"`
}

type audioRes struct {
	Text string `json:"text"`
}

func (s *AudioService) TranscribeAudio(ctx context.Context, b *bot.Bot, chatID int64, fileID string) (string, error) {
	params := &bot.GetFileParams{FileID: fileID}
	file, err := b.GetFile(ctx, params)
	if err != nil {
		return "", fmt.Errorf("AudioService.TranscribeAudio: cannot get audio file: %w", err)
	}

	downloadLink := b.FileDownloadLink(file)

	req := audioReq{
		URL:  downloadLink,
		Lang: "ru",
	}
	marshal, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("AudioService.TranscribeAudio: cannot marshal request body: %w", err)
	}

	resp, err := http.Post("http://localhost:8005/transcribe", "application/json", bytes.NewReader(marshal))
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
