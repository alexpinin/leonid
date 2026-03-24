package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"

	"leonid/src/internal/bot/dto"
	"leonid/src/internal/db"
)

type OpenAIService struct {
	config     OpenAIConfig
	executor   db.QueryExecutor
	configRepo configRepo
	llmClient  openai.Client
	chatLocks  sync.Map
}

type OpenAIConfig struct {
	BaseURL string
	Token   string
	Model   string
}

type configRepo interface {
	FindConfigByChatID(db.Executor, context.Context, int64) (dto.Config, error)
	UpdateConfig(db.Executor, context.Context, string, dto.Config) error
}

func NewOpenAIService(
	cfg OpenAIConfig,
	qe db.QueryExecutor,
	cr configRepo,
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
	s.chatMutex(chatID).Lock()
	defer s.chatMutex(chatID).Unlock()

	config, err := s.configRepo.FindConfigByChatID(s.executor.Executor(), ctx, chatID)
	if err != nil {
		return fmt.Errorf("OpenAIService.SendMessage: cannot find config: %w", err)
	}

	history, err := s.conversationHistory(config, message)
	if err != nil {
		return fmt.Errorf("OpenAIService.SendMessage: cannot build openai context: %w", err)
	}

	llmParams := openai.ChatCompletionNewParams{
		Messages: s.buildPrompt(config, history),
		Model:    s.config.Model,
	}
	completion, err := s.llmClient.Chat.Completions.New(ctx, llmParams)
	if err != nil {
		return fmt.Errorf("OpenAIService.SendMessage: cannot get LLM response: %w", err)
	}
	if len(completion.Choices) == 0 {
		return errors.New("OpenAIService.SendMessage: no ai choices")
	}

	response := completion.Choices[0].Message.Content
	config.ConversationHistory, err = s.historyToPersist(history, response)
	if err != nil {
		return fmt.Errorf("OpenAIService.SendMessage: cannot convert history to persist: %w", err)
	}

	err = s.configRepo.UpdateConfig(s.executor.Executor(), ctx, config.ID, config)
	if err != nil {
		return fmt.Errorf("OpenAIService.SendMessage: cannot update config: %w", err)
	}

	telegramParams := bot.SendMessageParams{
		ChatID: chatID,
		Text:   response,
	}
	_, err = b.SendMessage(ctx, &telegramParams)
	if err != nil {
		return fmt.Errorf("OpenAIService.SendMessage: cannot send Telegram message: %w", err)
	}
	return nil
}

func (_ *OpenAIService) conversationHistory(config dto.Config, message string) (dto.OpenAIConversationHistory, error) {
	history := dto.OpenAIConversationHistory{}
	err := json.Unmarshal([]byte(config.ConversationHistory), &history)
	if err != nil {
		return dto.OpenAIConversationHistory{}, err
	}

	if len(history.Messages) >= 10 {
		history.Messages = history.Messages[1:]
	}

	history.Messages = append(history.Messages, dto.OpenAIConversationMessage{
		OfUser: &openai.ChatCompletionUserMessageParam{
			Content: openai.ChatCompletionUserMessageParamContentUnion{
				OfString: param.Opt[string]{Value: message},
			},
		},
	})

	return history, nil
}

func (_ *OpenAIService) buildPrompt(config dto.Config, history dto.OpenAIConversationHistory) []openai.ChatCompletionMessageParamUnion {
	prompt := fmt.Sprintf("%s. Your nicknames are: %s",
		config.SystemPrompt, strings.Join(config.Nicknames, ","))

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

	for _, m := range history.Messages {
		messages = append(messages, openai.ChatCompletionMessageParamUnion{
			OfUser:      m.OfUser,
			OfAssistant: m.OfAssistant,
		})
	}

	return messages
}

func (_ *OpenAIService) historyToPersist(history dto.OpenAIConversationHistory, message string) (string, error) {
	if len(history.Messages) >= 10 {
		history.Messages = history.Messages[1:]
	}

	history.Messages = append(history.Messages, dto.OpenAIConversationMessage{
		OfAssistant: &openai.ChatCompletionAssistantMessageParam{
			Content: openai.ChatCompletionAssistantMessageParamContentUnion{
				OfString: param.Opt[string]{Value: message},
			},
		},
	})

	conversationContext, err := json.Marshal(history)
	if err != nil {
		return "", fmt.Errorf("failed to marshal conversation history: %w", err)
	}

	return string(conversationContext), nil
}

func (s *OpenAIService) chatMutex(chatID int64) *sync.Mutex {
	mu, _ := s.chatLocks.LoadOrStore(chatID, &sync.Mutex{})
	return mu.(*sync.Mutex)
}
