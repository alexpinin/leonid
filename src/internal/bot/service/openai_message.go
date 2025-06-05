package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"leonid/src/internal/bot/repository"
	"leonid/src/internal/common/db"
	"leonid/src/internal/common/logger"
	"os"

	"github.com/go-telegram/bot"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
)

type openAIMessage struct {
	OfUser      *openai.ChatCompletionUserMessageParam      `json:"ofUser"`
	OfAssistant *openai.ChatCompletionAssistantMessageParam `json:"ofAssistant"`
}

type openAIContext struct {
	Messages []openAIMessage `json:"messages"`
}

type OpenAIMessageService struct {
	db         *db.DB
	configRepo *repository.ConfigRepository
	llmClient  openai.Client
}

func NewOpenAIMessageService(db *db.DB, cr *repository.ConfigRepository) *OpenAIMessageService {
	return &OpenAIMessageService{
		db:         db,
		configRepo: cr,
		llmClient: openai.NewClient(
			option.WithAPIKey(os.Getenv("OPENAI_LLM_TOKEN")),
		),
	}
}

func (s *OpenAIMessageService) SendMessage(ctx context.Context, b *bot.Bot, chatID int64, message string) {
	err := s.db.ExecInTx(ctx, func(tx *sql.Tx) error {
		return s.sendMessage(ctx, b, chatID, message)
	})
	if err != nil {
		logger.Error(fmt.Sprintf("OpenAIMessageService.SendMessage: %v", err))
	}
}

func (s *OpenAIMessageService) sendMessage(ctx context.Context, b *bot.Bot, chatID int64, message string) error {
	config, err := s.configRepo.FindConfigByChatID(ctx, nil, chatID)
	if err != nil {
		return err
	}
	aiContext, err := s.buildOpenAIContext(config.ConversationContext, message)
	if err != nil {
		return err
	}

	messages := []openai.ChatCompletionMessageParamUnion{
		{
			OfSystem: &openai.ChatCompletionSystemMessageParam{
				Content: openai.ChatCompletionSystemMessageParamContentUnion{
					OfString: param.Opt[string]{Value: config.SystemPrompt},
				},
			},
		},
	}
	for _, m := range aiContext.Messages {
		messages = append(messages, openai.ChatCompletionMessageParamUnion{
			OfUser:      m.OfUser,
			OfAssistant: m.OfAssistant,
		})
	}

	resp, err := s.llmClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: messages,
		Model:    openai.ChatModelGPT4oMini,
	})
	if err != nil {
		return err
	}

	if len(resp.Choices) == 0 {
		return errors.New(fmt.Sprintf("no ai choices: %v", err))
	}
	answer := resp.Choices[0].Message.Content

	config.ConversationContext, err = s.buildConversationContext(aiContext, answer)
	if err != nil {
		return err
	}
	err = s.configRepo.UpdateConfig(ctx, nil, config.ID, config)
	if err != nil {
		return err
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   answer,
	})
	return nil
}

func (_ *OpenAIMessageService) buildOpenAIContext(conversationContext string, message string) (openAIContext, error) {
	aiContext := openAIContext{}
	err := json.Unmarshal([]byte(conversationContext), &aiContext)
	if err != nil {
		return openAIContext{}, err
	}

	if len(aiContext.Messages) >= 10 {
		aiContext.Messages = aiContext.Messages[1:]
	}

	aiContext.Messages = append(aiContext.Messages, openAIMessage{
		OfUser: &openai.ChatCompletionUserMessageParam{
			Content: openai.ChatCompletionUserMessageParamContentUnion{
				OfString: param.Opt[string]{Value: message},
			},
		},
	})

	return aiContext, nil
}

func (_ *OpenAIMessageService) buildConversationContext(aiContext openAIContext, message string) (string, error) {
	if len(aiContext.Messages) >= 10 {
		aiContext.Messages = aiContext.Messages[1:]
	}
	aiContext.Messages = append(aiContext.Messages, openAIMessage{
		OfAssistant: &openai.ChatCompletionAssistantMessageParam{
			Content: openai.ChatCompletionAssistantMessageParamContentUnion{
				OfString: param.Opt[string]{Value: message},
			},
		},
	})
	conversationContext, err := json.Marshal(aiContext)
	return string(conversationContext), err
}
