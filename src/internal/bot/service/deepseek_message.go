package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"

	"leonid/src/internal/bot/dto"
	"leonid/src/internal/bot/repo"
	"leonid/src/internal/db"
)

type OpenAIMessageService struct {
	executor   db.QueryExecutor
	configRepo *repo.ConfigRepo
	llmClient  openai.Client
}

type openAIMessage struct {
	OfUser      *openai.ChatCompletionUserMessageParam      `json:"ofUser"`
	OfAssistant *openai.ChatCompletionAssistantMessageParam `json:"ofAssistant"`
}

type openAIContext struct {
	Messages []openAIMessage `json:"messages"`
}

func NewOpenAIMessageService(
	qe db.QueryExecutor,
	cr *repo.ConfigRepo,
) *OpenAIMessageService {
	return &OpenAIMessageService{
		executor:   qe,
		configRepo: cr,
		llmClient: openai.NewClient(
			option.WithBaseURL("https://api.deepseek.com"),
			option.WithAPIKey(os.Getenv("DEEPSEEK_LLM_TOKEN")),
		),
	}
}

func (s *OpenAIMessageService) SendMessage(ctx context.Context, b *bot.Bot, chatID int64, message string) error {
	err := s.executor.ExecuteInTx(func(tx *sql.Tx) error {
		config, err := s.configRepo.FindConfigByChatID(tx, ctx, chatID)
		if err != nil {
			return err
		}
		aiContext, err := s.buildOpenAIContext(config, message)
		if err != nil {
			return err
		}

		messages := []openai.ChatCompletionMessageParamUnion{
			{
				OfSystem: &openai.ChatCompletionSystemMessageParam{
					Content: openai.ChatCompletionSystemMessageParamContentUnion{
						OfString: param.Opt[string]{Value: config.SystemPrompt +
							fmt.Sprintf(". Your nicknames are: %s", strings.Join(config.Nicknames, ","))},
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
			Model:    "deepseek-reasoner",
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
		err = s.configRepo.UpdateConfig(tx, ctx, config.ID, config)
		if err != nil {
			return err
		}

		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   answer,
		})
		return nil
	})
	if err != nil {
		return fmt.Errorf("OpenAIMessageService.SendMessage: %v", err)
	}
	return nil
}

func (_ *OpenAIMessageService) buildOpenAIContext(config dto.Config, message string) (openAIContext, error) {
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
