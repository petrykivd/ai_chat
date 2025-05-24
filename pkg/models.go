package pkg

import (
	"github.com/google/uuid"
	"time"
)

type Chat struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type Message struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	Role      string
	Content   string
	CreatedAt time.Time
}
