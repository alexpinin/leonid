package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"

	"leonid/src/internal/bot/dto"
	"leonid/src/internal/bot/repo"
	"leonid/src/internal/db"
)

type OpenAIService struct {
	config     OpenAIConfig
	executor   db.QueryExecutor
	configRepo *repo.ConfigRepo
	llmClient  openai.Client
}

type OpenAIConfig struct {
	BaseURL string
	Token   string
	Model   string
}

type openAIMessage struct {
	OfUser      *openai.ChatCompletionUserMessageParam      `json:"ofUser"`
	OfAssistant *openai.ChatCompletionAssistantMessageParam `json:"ofAssistant"`
}

type openAIContext struct {
	Messages []openAIMessage `json:"messages"`
}

func NewOpenAIService(
	cfg OpenAIConfig,
	qe db.QueryExecutor,
	cr *repo.ConfigRepo,
) *OpenAIService {
	return &OpenAIService{
		config:     cfg,
		executor:   qe,
		configRepo: cr,
		llmClient: openai.NewClient(
			option.WithBaseURL(cfg.BaseURL),
			option.WithAPIKey(cfg.Token),
		),
	}
}

func (s *OpenAIService) SendMessage(ctx context.Context, b *bot.Bot, chatID int64, message string) error {
	err := s.executor.ExecuteInTx(func(tx *sql.Tx) error {

		config, err := s.configRepo.FindConfigByChatID(tx, ctx, chatID)
		if err != nil {
			return fmt.Errorf("cannot find config: %w", err)
		}

		aiContext, err := s.buildAIContext(config, message)
		if err != nil {
			return fmt.Errorf("cannot build openai context: %w", err)
		}

		prompt := s.buildPrompt(config, aiContext)

		completion, err := s.llmClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Messages: prompt,
			Model:    s.config.Model,
		})
		if err != nil {
			return fmt.Errorf("cannot get LLM response: %w", err)
		}

		if len(completion.Choices) == 0 {
			return errors.New(fmt.Sprintf("no ai choices"))
		}

		response := completion.Choices[0].Message.Content

		config.ConversationContext, err = s.buildConversationContext(aiContext, response)
		if err != nil {
			return err
		}
		err = s.configRepo.UpdateConfig(tx, ctx, config.ID, config)
		if err != nil {
			return err
		}

		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   response,
		})
		if err != nil {
			return fmt.Errorf("cannot send Telegram message: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("OpenAIService.SendMessage: %w", err)
	}
	return nil
}

func (_ *OpenAIService) buildAIContext(config dto.Config, message string) (openAIContext, error) {
	aiContext := openAIContext{}
	err := json.Unmarshal([]byte(config.ConversationContext), &aiContext)
	if err != nil {
		return openAIContext{}, err
	}

	if len(aiContext.Messages) >= 10 {
		aiContext.Messages = aiContext.Messages[1:]
	}

	aiContext.Messages = append(aiContext.Messages, openAIMessage{
		OfUser: &openai.ChatCompletionUserMessageParam{
			Content: openai.ChatCompletionUserMessageParamContentUnion{
				OfString: param.Opt[string]{Value: message + fmt.Sprintf(". %s", config.MessagePrompt)},
			},
		},
	})

	return aiContext, nil
}

func (_ *OpenAIService) buildPrompt(cfg dto.Config, ctx openAIContext) []openai.ChatCompletionMessageParamUnion {
	prompt := fmt.Sprintf("%s. Your nicknames are: %s", cfg.SystemPrompt, strings.Join(cfg.Nicknames, ","))

	messages := []openai.ChatCompletionMessageParamUnion{
		{
			OfSystem: &openai.ChatCompletionSystemMessageParam{
				Content: openai.ChatCompletionSystemMessageParamContentUnion{
					OfString: param.Opt[string]{
						Value: prompt,
					},
				},
			},
		},
	}

	for _, m := range ctx.Messages {
		messages = append(messages, openai.ChatCompletionMessageParamUnion{
			OfUser:      m.OfUser,
			OfAssistant: m.OfAssistant,
		})
	}

	return messages
}

func (_ *OpenAIService) buildConversationContext(aiContext openAIContext, message string) (string, error) {
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
