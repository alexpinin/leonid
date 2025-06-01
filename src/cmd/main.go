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
	if update == nil || update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	validChat, err := isValidChat(chatID)
	if err != nil {
		return
	}

	message := strings.ToLower(update.Message.Text)
	if !validChat {
		validPass, err := isPassPhrase(message)
		if err != nil {
			return
		}
		if validPass {
			err = connectPassPhraseToChat(message, chatID)
			if err != nil {
				return
			}
		} else {
			return
		}
	}
	if !isMessageToAnswer(message) && !isMessageResponse(message) {
		return
	}

	message = removeNickname(message)

	handle(ctx, b, message, chatID)
}

func isMessageResponse(message string) bool {
	return true
}

func handle(ctx context.Context, b *bot.Bot, message string, chatID int64) {
	llmRequest := fmt.Sprintf(`Ты не очень образованный и обидчивый человек.
		Очень кратко ответь на вопрос, обращаясь к собеседнику "братишка".
		Вопрос: "%s"`, message)

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
		ChatID: chatID,
		Text:   resp.Choices[0].Message.Content.Text,
	})
}

func isMessageToAnswer(s string) bool {
	for _, nickname := range nicknames {
		if strings.Contains(s, nickname) {
			return true
		}
	}
	return false
}

func isValidChat(chatID int64) (bool, error) {
	return true, nil
}

func isPassPhrase(s string) (bool, error) {
	return true, nil
}

func connectPassPhraseToChat(passPhrase string, chatID int64) error {
	return nil
}

func removePassPhrase(s string) error {
	return nil
}

func activateChat(chatID int64) error {
	return nil
}

func removeNickname(s string) string {
	for _, nickname := range nicknames {
		if strings.Contains(s, nickname) {
			s = strings.ReplaceAll(s, nickname, "")
		}
	}
	return s
}
