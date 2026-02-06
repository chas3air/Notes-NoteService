package service

import (
	"context"
	"errors"
	"notesservice/internal/models/domain"

	"github.com/google/uuid"
)

var (
	ErrNotFound      = errors.New("note not found")
	ErrAlreadyExists = errors.New("note already exists")
)

type QueryStorage interface {
	GetNotes(ctx context.Context, offset int, limit int) ([]domain.Note, error)
	GetNotesByUser(ctx context.Context, userId uuid.UUID, offset int, limit int) ([]domain.Note, error)
	GetNoteById(ctx context.Context, id uuid.UUID) (domain.Note, error)
}

type CommandStorage interface {
	Insert(ctx context.Context, note domain.Note) error
	Update(ctx context.Context, id uuid.UUID, note domain.Note) error
	Delete(ctx context.Context, id uuid.UUID) error
}
