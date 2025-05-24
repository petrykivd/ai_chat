package messages

import (
	"ai_chat/pkg"
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateMessage(ctx context.Context, message pkg.Message) (uuid.UUID, error)
	ListMessages(ctx context.Context, chatID uuid.UUID) ([]pkg.Message, error)
	GetMessage(ctx context.Context, id uuid.UUID) (pkg.Message, error)
}
