package handlers

import (
	"context"
	"notesservice/internal/models/domain"

	"github.com/google/uuid"
)

type Service interface {
	GetNotes(ctx context.Context, offset int, limit int) ([]domain.Note, error)
	GetNotesByUser(ctx context.Context, userId uuid.UUID, offset int, limit int) ([]domain.Note, error)
	GetNoteById(ctx context.Context, id uuid.UUID) (domain.Note, error)
	Insert(ctx context.Context, note domain.Note) error
	Update(ctx context.Context, id uuid.UUID, note domain.Note) error
	Delete(ctx context.Context, id uuid.UUID) error
}
