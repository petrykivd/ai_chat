package messages

import (
	"ai_chat/pkg"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresMessageRepository struct {
	db *gorm.DB
}

func NewPostgresMessageRepository(db *gorm.DB) *PostgresMessageRepository {
	return &PostgresMessageRepository{db: db}
}

func (r *PostgresMessageRepository) CreateMessage(ctx context.Context, message pkg.Message) (uuid.UUID, error) {
	if message.ID == uuid.Nil {
		message.ID = uuid.New()
	}
	if message.CreatedAt.IsZero() {
		message.CreatedAt = time.Now()
	}
	if err := r.db.WithContext(ctx).Create(&message).Error; err != nil {
		return uuid.Nil, err
	}
	return message.ID, nil
}

func (r *PostgresMessageRepository) GetMessage(ctx context.Context, id uuid.UUID) (pkg.Message, error) {
	var message pkg.Message
	if err := r.db.WithContext(ctx).First(&message, "id = ?", id).Error; err != nil {
		return pkg.Message{}, err
	}
	return message, nil
}

func (r *PostgresMessageRepository) ListMessages(ctx context.Context, chatID uuid.UUID) ([]pkg.Message, error) {
	var messages []pkg.Message
	if err := r.db.WithContext(ctx).Find(&messages, "chat_id = ?", chatID).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
