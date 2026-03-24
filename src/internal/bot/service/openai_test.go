package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/openai/openai-go"

	"leonid/src/internal/bot/dto"
	"leonid/src/internal/db"
	"leonid/src/internal/testutil"
)

func TestOpenAIService(t *testing.T) {
	b := &mockBot{}
	//cr := &mockConfigRepo{
	//	findConfigByChatIDRes: dto.Config{
	//		ConversationHistory: "{}",
	//	},
	//}
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
		testutil.Equal(t, dto.Config{
			ConversationHistory: `{"messages":[{"ofUser":{"content":"user message","role":"user"},"ofAssistant":null},{"ofUser":null,"ofAssistant":{"content":"LLM message","role":"assistant"}}]}`,
		}, cr.updateConfigIn3)
	})

	t.Run("should return error when config not found", func(t *testing.T) {

	})

	t.Run("should return error when config repo fails", func(t *testing.T) {

	})

	t.Run("should unmarshal existing conversation history", func(t *testing.T) {

	})

	t.Run("should return error when conversation history is invalid JSON", func(t *testing.T) {

	})

	t.Run("should trim history when it reaches 10 messages", func(t *testing.T) {

	})

	t.Run("should return error when LLM call fails", func(t *testing.T) {

	})

	t.Run("should return error when LLM returns no choices", func(t *testing.T) {

	})

	t.Run("should return error when config update fails", func(t *testing.T) {

	})

	t.Run("should append assistant response to persisted history", func(t *testing.T) {

	})

	t.Run("should return error when Telegram send fails", func(t *testing.T) {

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
