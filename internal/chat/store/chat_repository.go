package chat

import (
	"ai_chat/pkg"
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	CreateChat(ctx context.Context, chat pkg.Chat) (uuid.UUID, error)
	GetChat(ctx context.Context, id uuid.UUID) (pkg.Chat, error)
	ListChats(ctx context.Context) ([]pkg.Chat, error)
}
