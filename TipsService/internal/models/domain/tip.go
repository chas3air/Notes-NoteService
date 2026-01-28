package domain

import (
	"time"

	"github.com/google/uuid"
)

type Tip struct {
	Id        uuid.UUID `json:"id,omitempty"`
	UserId    uuid.UUID `json:"user_id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	IsPrivate bool      `json:"is_private,omitempty"`
}
