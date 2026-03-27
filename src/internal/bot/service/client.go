package service

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIConfig struct {
	BaseURL string
	Token   string
	Model   string
}

type OpenAIClient struct {
	client openai.Client
	model  string
}

func NewOpenAIClient(cfg OpenAIConfig) *OpenAIClient {
	return &OpenAIClient{
		model: cfg.Model,
		client: openai.NewClient(
			option.WithBaseURL(cfg.BaseURL),
			option.WithAPIKey(cfg.Token),
		),
	}
}

func (c *OpenAIClient) CreateChatCompletion(ctx context.Context, params openai.ChatCompletionNewParams) (
	*openai.ChatCompletion, error) {
	return c.client.Chat.Completions.New(ctx, params)
}

func (c *OpenAIClient) Model() string {
	return c.model
}
