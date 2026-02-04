package domain

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	Id        uuid.UUID `json:"id,omitempty"`
	UserId    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	IsPrivate bool      `json:"is_private"`
}
