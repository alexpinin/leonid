package handler

import (
	"context"
	"fmt"
	"os"

	"github.com/go-telegram/bot"
	"github.com/revrost/go-openrouter"
)

type MessageSender struct {
	basicHandler
	llmClient *openrouter.Client
	promptProvider
}

func NewMessageSender(pp promptProvider) *MessageSender {
	return &MessageSender{
		llmClient: openrouter.NewClient(
			os.Getenv("LLM_TOKEN"),
		),
		promptProvider: pp,
	}
}

type promptProvider interface {
	FindSystemPrompt(context.Context, int64) string
}

func (h *MessageSender) Handle(ctx context.Context, b *bot.Bot, u *UpdateContext) {
	systemPrompt := h.promptProvider.FindSystemPrompt(ctx, u.Message.Chat.ID)
	if systemPrompt == "" {
		return
	}

	llmRequest := fmt.Sprintf(`%s "%s"`, systemPrompt, u.Message.Text)
	resp, err := h.llmClient.CreateChatCompletion(
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
		return
	}

	if len(resp.Choices) == 0 {
		return
	}

	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: u.Message.Chat.ID,
		Text:   resp.Choices[0].Message.Content.Text,
	})
	h.nextHandle(ctx, b, u)
}
