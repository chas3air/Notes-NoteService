package domain

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	Id        uuid.UUID `json:"id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserId    uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Title     string    `json:"title" example:"Купить продукты"`
	Content   string    `json:"content" example:"Молоко, сыр, хлеб"`
	CreatedAt time.Time `json:"created_at" example:"2026-02-05T13:40:00Z"`
	IsPrivate bool      `json:"is_private" example:"true"`
}
