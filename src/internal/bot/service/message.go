package service

import (
	"context"
	"fmt"
	"leonid/src/internal/bot/repository"
	"leonid/src/internal/common/logger"
	"os"

	"github.com/go-telegram/bot"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
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

type OpenAIMessageService struct {
	configRepo *repository.ConfigRepository
	llmClient  openai.Client
}

func NewOpenAIMessageService(cr *repository.ConfigRepository) *OpenAIMessageService {
	return &OpenAIMessageService{
		llmClient: openai.NewClient(
			option.WithAPIKey(os.Getenv("OPENAI_LLM_TOKEN")),
		),
		configRepo: cr,
	}
}

func (s *OpenAIMessageService) SendMessage(ctx context.Context, b *bot.Bot, chatID int64, message string) {
	config, err := s.configRepo.FindConfigByChatID(ctx, nil, chatID)
	if err != nil {
		logger.Error(fmt.Sprintf("OpenAIMessageService.SendMessage: %v", err))
		return
	}

	resp, err := s.llmClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfSystem: &openai.ChatCompletionSystemMessageParam{
					Content: openai.ChatCompletionSystemMessageParamContentUnion{
						OfString: param.Opt[string]{Value: config.SystemPrompt},
					},
				},
			},
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: param.Opt[string]{Value: message},
					},
				},
			},
		},
		Model: openai.ChatModelGPT4oMini,
	})
	if err != nil {
		logger.Error(fmt.Sprintf("OpenAIMessageService.SendMessage: %v", err))
		return
	}

	if len(resp.Choices) == 0 {
		return
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   resp.Choices[0].Message.Content,
	})
}
