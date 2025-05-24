package chat

import (
	"ai_chat/pkg"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresChatRepository struct {
	db *gorm.DB
}

func NewPostgresChatRepository(db *gorm.DB) *PostgresChatRepository {
	return &PostgresChatRepository{db: db}
}

func (r *PostgresChatRepository) CreateChat(ctx context.Context, chat pkg.Chat) (uuid.UUID, error) {
	if chat.ID == uuid.Nil {
		chat.ID = uuid.New()
	}
	if chat.CreatedAt.IsZero() {
		chat.CreatedAt = time.Now()
	}
	if err := r.db.WithContext(ctx).Create(&chat).Error; err != nil {
		return uuid.Nil, err
	}
	return chat.ID, nil
}

func (r *PostgresChatRepository) GetChat(ctx context.Context, id uuid.UUID) (pkg.Chat, error) {
	var chat pkg.Chat
	if err := r.db.WithContext(ctx).First(&chat, "id = ?", id).Error; err != nil {
		return pkg.Chat{}, err
	}
	return chat, nil
}

func (r *PostgresChatRepository) ListChats(ctx context.Context) ([]pkg.Chat, error) {
	var chats []pkg.Chat
	if err := r.db.WithContext(ctx).Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}
