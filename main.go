package main

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/revrost/go-openrouter"
	"os"
	"os/signal"
	"strings"
)

var (
	botToken = os.Getenv("BOT_TOKEN")
	llmToken = os.Getenv("LLM_TOKEN")

	llmClient = openrouter.NewClient(
		llmToken,
	)
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
		bot.WithAllowedUpdates(bot.AllowedUpdates{
			"message",
		}),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

var nicknames = []string{
	"леня",
	"лёня",
	"ленька",
	"лёнька",
	"леонид",
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	//todo check chat id and message date
	if update == nil || update.Message == nil {
		return
	}
	lower := strings.ToLower(update.Message.Text)
	contains := false
	for _, nickname := range nicknames {
		if strings.Contains(lower, nickname) {
			lower = strings.ReplaceAll(lower, nickname, "")
			contains = true
		}
	}
	if !contains {
		return
	}
	println(lower)
	llmRequest := fmt.Sprintf(`Ты не очень образованный и обидчивый человек.
		Очень кратко ответь на вопрос, обращаясь к собеседнику "братишка".
		Вопрос: "%s"`, lower)

	// todo add quota
	// todo return error
	resp, err := llmClient.CreateChatCompletion(
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

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   resp.Choices[0].Message.Content.Text,
	})
}
