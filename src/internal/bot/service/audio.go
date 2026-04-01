package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-telegram/bot"
)

type AudioService struct {
}

func NewAudioService() *AudioService {
	return &AudioService{}
}

//FileID = {string} "AwACAgIAAxkBAAOzac2O8wnx6ud0zjCvxnzBhCQfuaUAAkadAAIeQ3FK669RI1w1kHA6BA"
//FileUniqueID = {string} "AgADRp0AAh5DcUo"
//Duration = {int} 2
//MimeType = {string} "audio/ogg"
//FileSize = {int64} 9859

func (s *AudioService) TranscribeAudio(ctx context.Context, b *bot.Bot, chatID int64, fileID string) (string, error) {
	params := &bot.GetFileParams{FileID: fileID}
	file, err := b.GetFile(ctx, params)
	if err != nil {
		return "", fmt.Errorf("AudioService.TranscribeAudio: cannot get audio file: %w", err)
	}

	downloadLink := b.FileDownloadLink(file)

	resp, err := http.Get(downloadLink)
	if err != nil {
		return "", fmt.Errorf("AudioService.TranscribeAudio: cannot download audio file: %w", err)
	}
	defer resp.Body.Close()

	var bb []byte
	n, err := resp.Body.Read(bb)
	if err != nil {
		return "", fmt.Errorf("AudioService.TranscribeAudio: cannot read audio file: %w", err)
	}

	println(downloadLink)
	return "", err
}
