package llm

import (
	"ai_chat/pkg"
	"context"
)

type Message struct {
	Role    string
	Content string
}

type Client interface {
	ChatComplete(ctx context.Context, history []Message) (Message, error)
	MessageFromModel(m pkg.Message) Message
	//ChatStream(ctx context.Context, history []Message) (<-chan Message, error)
}
