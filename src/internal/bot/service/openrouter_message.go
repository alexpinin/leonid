package service

import (
	"context"
	"fmt"
	"os"

	"leonid/src/internal/bot/repository"
	"leonid/src/internal/logger"

	"github.com/go-telegram/bot"
	"github.com/revrost/go-openrouter"
)

type OpenrouterMessageService struct {
	llmClient  *openrouter.Client
	configRepo *repository.ConfigRepository
}

func NewOpenrouterMessageService(cr *repository.ConfigRepository) *OpenrouterMessageService {
	return &OpenrouterMessageService{
		llmClient: openrouter.NewClient(
			os.Getenv("LLM_TOKEN"),
		),
		configRepo: cr,
	}
}

func (s *OpenrouterMessageService) SendMessage(ctx context.Context, b *bot.Bot, chatID int64, message string) {
	config, err := s.configRepo.FindConfigByChatID(ctx, nil, chatID)
	if err != nil {
		logger.Error(fmt.Sprintf("OpenrouterMessageService.SendMessage: %v", err))
		return
	}

	llmRequest := fmt.Sprintf(`%s "%s"`, config.SystemPrompt, message)
	resp, err := s.llmClient.CreateChatCompletion(
		context.Background(),
		openrouter.ChatCompletionRequest{
			Model: "deepseek/deepseek-r1-0528-qwen3-8b:free",
			Messages: []openrouter.ChatCompletionMessage{
				{
					Role:    openrouter.ChatMessageRoleUser,
					Content: openrouter.Content{Text: llmRequest},
				},
			},
		},
	)
	if err != nil {
		logger.Error(fmt.Sprintf("OpenrouterMessageService.SendMessage: %v", err))
		return
	}

	if len(resp.Choices) == 0 {
		return
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   resp.Choices[0].Message.Content.Text,
	})
}
