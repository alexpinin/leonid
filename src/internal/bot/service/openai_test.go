package service

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"

	"leonid/src/internal/bot/dto"
	"leonid/src/internal/db"
	"leonid/src/internal/testutil"
)

func TestOpenAIService(t *testing.T) {
	b := &mockBot{}
	cr := &mockConfigRepo{
		findConfigByChatIDRes: dto.Config{
			ConversationHistory: "{}",
		},
	}
	llmMessage := "LLM message"
	lc := &mockLLMClient{
		createChatCompletionRes: &openai.ChatCompletion{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Content: llmMessage,
					},
				},
			},
		},
	}

	ctx := context.Background()
	chatID := int64(123)
	userMessage := "user message"

	t.Run("should send message and save conversation history", func(t *testing.T) {
		cr := &mockConfigRepo{
			findConfigByChatIDRes: dto.Config{
				ConversationHistory: "{}",
			},
		}
		sut := NewOpenAIService(mockQueryExecutor{}, cr, lc)

		err := sut.SendMessage(ctx, b, chatID, userMessage)

		testutil.Equal(t, nil, err)
		testutil.Equal(t, llmMessage, b.sendMessageIn1.Text)

		history := dto.OpenAIConversationHistory{
			Messages: []dto.OpenAIConversationMessage{
				{
					OfUser: &openai.ChatCompletionUserMessageParam{
						Content: openai.ChatCompletionUserMessageParamContentUnion{
							OfString: param.Opt[string]{Value: userMessage},
						},
					},
				},
				{
					OfAssistant: &openai.ChatCompletionAssistantMessageParam{
						Content: openai.ChatCompletionAssistantMessageParamContentUnion{
							OfString: param.Opt[string]{Value: llmMessage},
						},
					},
				},
			},
		}
		expected := dto.Config{ConversationHistory: string(testutil.MustMarshalJson(t, history))}
		testutil.Equal(t, expected, cr.updateConfigIn3)
	})

	t.Run("should trim history when it reaches 10 messages", func(t *testing.T) {
		history := dto.OpenAIConversationHistory{
			Messages: make([]dto.OpenAIConversationMessage, 10),
		}
		for i := range len(history.Messages) {
			history.Messages[i] = dto.OpenAIConversationMessage{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: param.Opt[string]{Value: fmt.Sprintf("message_%d", i)},
					},
				},
			}
		}
		cr := &mockConfigRepo{
			findConfigByChatIDRes: dto.Config{
				ConversationHistory: string(testutil.MustMarshalJson(t, history)),
			},
		}
		sut := NewOpenAIService(mockQueryExecutor{}, cr, lc)

		err := sut.SendMessage(ctx, b, chatID, userMessage)

		testutil.Equal(t, nil, err)

		history.Messages = append(history.Messages, dto.OpenAIConversationMessage{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: param.Opt[string]{Value: userMessage},
				},
			},
		})
		history.Messages = append(history.Messages, dto.OpenAIConversationMessage{
			OfAssistant: &openai.ChatCompletionAssistantMessageParam{
				Content: openai.ChatCompletionAssistantMessageParamContentUnion{
					OfString: param.Opt[string]{Value: llmMessage},
				},
			},
		})
		history.Messages = history.Messages[2:]
		expected := dto.Config{ConversationHistory: string(testutil.MustMarshalJson(t, history))}
		testutil.Equal(t, expected, cr.updateConfigIn3)
	})

	t.Run("should return error if FindConfigByChatID returns error", func(t *testing.T) {
		cr := &mockConfigRepo{findConfigByChatIDErr: testutil.TestError}
		sut := NewOpenAIService(mockQueryExecutor{}, cr, lc)

		err := sut.SendMessage(ctx, b, chatID, userMessage)

		testutil.ErrorIs(t, testutil.TestError, err)
	})

	t.Run("should return error if Unmarshal returns error", func(t *testing.T) {
		cr := &mockConfigRepo{findConfigByChatIDRes: dto.Config{}}
		sut := NewOpenAIService(mockQueryExecutor{}, cr, lc)

		err := sut.SendMessage(ctx, b, chatID, userMessage)

		testutil.ErrorContains(t, "unexpected end of JSON input", err)
	})

	t.Run("should return error if CreateChatCompletion returns error", func(t *testing.T) {
		lc := &mockLLMClient{createChatCompletionErr: testutil.TestError}
		sut := NewOpenAIService(mockQueryExecutor{}, cr, lc)

		err := sut.SendMessage(ctx, b, chatID, userMessage)

		testutil.ErrorIs(t, testutil.TestError, err)
	})

	t.Run("should return error if UpdateConfig returns error", func(t *testing.T) {
		cr := &mockConfigRepo{
			findConfigByChatIDRes: dto.Config{ConversationHistory: "{}"},
			updateConfigErr:       testutil.TestError,
		}
		sut := NewOpenAIService(mockQueryExecutor{}, cr, lc)

		err := sut.SendMessage(ctx, b, chatID, userMessage)

		testutil.ErrorIs(t, testutil.TestError, err)
	})

	t.Run("should return error if SendMessage returns error", func(t *testing.T) {
		b := &mockBot{
			sendMessageErr: testutil.TestError,
		}
		sut := NewOpenAIService(mockQueryExecutor{}, cr, lc)

		err := sut.SendMessage(ctx, b, chatID, userMessage)

		testutil.ErrorIs(t, testutil.TestError, err)
	})

	t.Run("should send LLM response text to correct chat ID", func(t *testing.T) {

	})

	t.Run("should serialize messages for the same chat ID", func(t *testing.T) {

	})

	t.Run("should process messages for different chats in parallel", func(t *testing.T) {

	})
}

type mockQueryExecutor struct {
}

func (mockQueryExecutor) ExecuteInTx(func(tx *sql.Tx) error) error {
	return nil
}

func (mockQueryExecutor) Executor() db.Executor {
	return nil
}

type mockConfigRepo struct {
	findConfigByChatIDRes dto.Config
	findConfigByChatIDErr error
	updateConfigIn3       dto.Config
	updateConfigErr       error
}

func (m *mockConfigRepo) FindConfigByChatID(db.Executor, context.Context, int64) (dto.Config, error) {
	return m.findConfigByChatIDRes, m.findConfigByChatIDErr
}

func (m *mockConfigRepo) UpdateConfig(_ db.Executor, _ context.Context, _ string, in3 dto.Config) error {
	m.updateConfigIn3 = in3
	return m.updateConfigErr
}

type mockLLMClient struct {
	createChatCompletionRes *openai.ChatCompletion
	createChatCompletionErr error
	modelRes                string
}

func (c *mockLLMClient) CreateChatCompletion(context.Context, openai.ChatCompletionNewParams) (
	*openai.ChatCompletion, error) {
	return c.createChatCompletionRes, c.createChatCompletionErr
}

func (c *mockLLMClient) Model() string {
	return c.modelRes
}

type mockBot struct {
	sendMessageIn1 *bot.SendMessageParams
	sendMessageRes *models.Message
	sendMessageErr error
}

func (b *mockBot) SendMessage(_ context.Context, in1 *bot.SendMessageParams) (*models.Message, error) {
	b.sendMessageIn1 = in1
	return b.sendMessageRes, b.sendMessageErr
}
