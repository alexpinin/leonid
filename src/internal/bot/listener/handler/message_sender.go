package handler

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/revrost/go-openrouter"
	"os"
)

type MessageSender struct {
	basicHandler
	llmClient *openrouter.Client
}

func NewMessageSender() *MessageSender {
	return &MessageSender{
		llmClient: openrouter.NewClient(
			os.Getenv("LLM_TOKEN"),
		),
	}
}

func (h *MessageSender) Handle(ctx context.Context, b *bot.Bot, u *UpdateContext) {
	message := u.Message.Text

	llmRequest := fmt.Sprintf(`Ты не очень образованный и обидчивый человек.
		Очень кратко ответь на вопрос, обращаясь к собеседнику "братишка".
		Вопрос: "%s"`, message)

	// todo add quota
	// todo return error
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
