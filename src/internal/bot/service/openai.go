package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"

	"leonid/src/internal/bot/dto"
	"leonid/src/internal/db"
)

const (
	maxMessageHistoryLen = 10
	llmRequestTimeout    = time.Duration(60) * time.Second
)

type OpenAIService struct {
	executor   db.QueryExecutor
	configRepo ConfigRepo
	client     llmClient
	chatLocks  sync.Map
}

func NewOpenAIService(
	qe db.QueryExecutor,
	cr ConfigRepo,
	lc llmClient,
) *OpenAIService {
	return &OpenAIService{
		executor:   qe,
		configRepo: cr,
		client:     lc,
	}
}

type llmClient interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionNewParams) (*openai.ChatCompletion, error)
	Model() string
}

func (s *OpenAIService) InquireLLM(ctx context.Context, chatID int64, message string) (string, error) {
	// Per-chat mutex serializes access instead of a DB transaction
	// to avoid holding SQLite's write lock during the LLM call.
	s.chatMutex(chatID).Lock()
	defer s.chatMutex(chatID).Unlock()

	config, err := s.configRepo.FindConfigByChatID(s.executor.Executor(), ctx, chatID)
	if err != nil {
		return "", fmt.Errorf("OpenAIService.InquireLLM: cannot find config: %w", err)
	}

	history, err := s.conversationHistory(config, message)
	if err != nil {
		return "", fmt.Errorf("OpenAIService.InquireLLM: cannot build LLM context: %w", err)
	}

	llmParams := openai.ChatCompletionNewParams{
		Messages: s.buildPrompt(config, history),
		Model:    s.client.Model(),
	}

	reqCtx, cancel := context.WithTimeout(ctx, llmRequestTimeout)
	defer cancel()

	completion, err := s.client.CreateChatCompletion(reqCtx, llmParams)
	if err != nil {
		return "", fmt.Errorf("OpenAIService.InquireLLM: cannot get LLM response: %w", err)
	}
	if len(completion.Choices) == 0 {
		return "", errors.New("OpenAIService.InquireLLM: no ai choices")
	}

	reply := completion.Choices[0].Message.Content
	if reply == "" {
		return "", errors.New("OpenAIService.InquireLLM: empty reply from LLM")
	}

	config.ConversationHistory, err = s.historyToPersist(history, reply)
	if err != nil {
		return "", fmt.Errorf("OpenAIService.InquireLLM: cannot convert history to persist: %w", err)
	}

	err = s.configRepo.UpdateConfig(s.executor.Executor(), ctx, config)
	if err != nil {
		return "", fmt.Errorf("OpenAIService.InquireLLM: cannot update config: %w", err)
	}

	return reply, nil
}

func (_ *OpenAIService) conversationHistory(config dto.Config, message string) (dto.OpenAIConversationHistory, error) {
	history := dto.OpenAIConversationHistory{}
	err := json.Unmarshal([]byte(config.ConversationHistory), &history)
	if err != nil {
		return dto.OpenAIConversationHistory{}, fmt.Errorf("failed to unmarshal conversation history: %w", err)
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
	history.Messages = append(history.Messages, dto.OpenAIConversationMessage{
		OfAssistant: &openai.ChatCompletionAssistantMessageParam{
			Content: openai.ChatCompletionAssistantMessageParamContentUnion{
				OfString: param.Opt[string]{Value: message},
			},
		},
	})

	if len(history.Messages) > maxMessageHistoryLen {
		history.Messages = history.Messages[len(history.Messages)-maxMessageHistoryLen:]
	}

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
