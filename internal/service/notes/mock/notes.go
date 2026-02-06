package mock

import (
	"context"
	"errors"
	"notesservice/internal/models/domain"

	"github.com/google/uuid"
)

type MockStorage struct {
	GetNotesFunc       func(ctx context.Context, offset int, limit int) ([]domain.Note, error)
	GetNotesByUserFunc func(ctx context.Context, userId uuid.UUID, offset int, limit int) ([]domain.Note, error)
	GetNoteByIdFunc    func(ctx context.Context, id uuid.UUID) (domain.Note, error)
	InsertFunc         func(ctx context.Context, note domain.Note) error
	UpdateFunc         func(ctx context.Context, id uuid.UUID, note domain.Note) error
	DeleteFunc         func(ctx context.Context, id uuid.UUID) error
}

func (m *MockStorage) GetNotes(ctx context.Context, offset int, limit int) ([]domain.Note, error) {
	if m.GetNotesFunc == nil {
		return nil, errors.New("GetNotes not implemented")
	}
	return m.GetNotesFunc(ctx, offset, limit)
}

func (m *MockStorage) GetNotesByUser(ctx context.Context, userId uuid.UUID, offset int, limit int) ([]domain.Note, error) {
	if m.GetNotesByUserFunc == nil {
		return nil, errors.New("GetNotesByUser not implemented")
	}
	return m.GetNotesByUserFunc(ctx, userId, offset, limit)
}

func (m *MockStorage) GetNoteById(ctx context.Context, id uuid.UUID) (domain.Note, error) {
	if m.GetNoteByIdFunc == nil {
		return domain.Note{}, errors.New("GetNoteById not implemented")
	}
	return m.GetNoteByIdFunc(ctx, id)
}

func (m *MockStorage) Insert(ctx context.Context, note domain.Note) error {
	if m.InsertFunc == nil {
		return errors.New("Insert not implemented")
	}
	return m.InsertFunc(ctx, note)
}

func (m *MockStorage) Update(ctx context.Context, id uuid.UUID, note domain.Note) error {
	if m.UpdateFunc == nil {
		return errors.New("Update not implemented")
	}
	return m.UpdateFunc(ctx, id, note)
}

func (m *MockStorage) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc == nil {
		return errors.New("Delete not implemented")
	}
	return m.DeleteFunc(ctx, id)
}
