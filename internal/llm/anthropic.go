package llm

import (
	"ai_chat/internal/config"
	"ai_chat/pkg"
	"context"
	"fmt"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AnthropicClient struct {
	cli *anthropic.Client
	cfg config.AnthropicConfig
}

func NewAnthropicClient(cfg config.AnthropicConfig) *AnthropicClient {
	cli := anthropic.NewClient(
		option.WithAPIKey(cfg.APIKey),
	)
	return &AnthropicClient{cli: &cli, cfg: cfg}
}

func (a *AnthropicClient) ChatComplete(ctx context.Context, history []Message) (Message, error) {
	var messages []anthropic.MessageParam
	for _, m := range history {
		switch m.Role {
		case "user":
			messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(m.Content)))
		case "assistant":
			messages = append(messages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(m.Content)))
		}
	}

	resp, err := a.cli.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_0,
		MaxTokens: a.cfg.MaxTokens,
		Messages:  messages,
	})
	if err != nil {
		return Message{}, fmt.Errorf("anthropic ChatComplete: %w", err)
	}
	var reply string
	for _, block := range resp.Content {
		reply += block.Text
	}

	return Message{
		Role:    "assistant",
		Content: reply,
	}, nil
}

func (a *AnthropicClient) MessageFromModel(m pkg.Message) Message {
	return Message{
		Role:    m.Role,
		Content: m.Content,
	}
}
