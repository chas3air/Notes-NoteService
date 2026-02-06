package notes_test

import (
	"context"
	"errors"
	"notesservice/internal/models/domain"
	"notesservice/internal/service"
	notes "notesservice/internal/service/notes"
	"notesservice/internal/service/notes/mock"
	storageerrors "notesservice/internal/storage"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestGetNotes(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		GetNotesFunc: func(ctx context.Context, offset, limit int) ([]domain.Note, error) {
			return []domain.Note{}, nil
		},
	}

	svc := notes.New(log, mockStorage)

	result, err := svc.GetNotes(ctx, 0, 2)
	require.NoError(t, err)
	require.Len(t, result, 0)
}

func TestGetNotes_Error(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		GetNotesFunc: func(ctx context.Context, offset, limit int) ([]domain.Note, error) {
			return nil, errors.New("db error")
		},
	}

	svc := notes.New(log, mockStorage)

	_, err := svc.GetNotes(ctx, 0, 2)
	require.Error(t, err)
}

func TestGetNotesByUser_Success(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	userID := uuid.New()

	mockStorage := &mock.MockStorage{
		GetNotesByUserFunc: func(ctx context.Context, id uuid.UUID, offset, limit int) ([]domain.Note, error) {
			require.Equal(t, userID, id)
			return []domain.Note{
				{Id: uuid.New(), UserId: userID, Title: "A"},
			}, nil
		},
	}

	svc := notes.New(log, mockStorage)

	result, err := svc.GetNotesByUser(ctx, userID, 0, 10)
	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestGetNotesByUser_Error(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	userID := uuid.New()

	mockStorage := &mock.MockStorage{
		GetNotesByUserFunc: func(ctx context.Context, id uuid.UUID, offset, limit int) ([]domain.Note, error) {
			return nil, errors.New("db error")
		},
	}

	svc := notes.New(log, mockStorage)

	_, err := svc.GetNotesByUser(ctx, userID, 0, 10)
	require.Error(t, err)
}

func TestGetNoteById_Success(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	id := uuid.New()
	expected := domain.Note{Id: id, Title: "Test"}

	mockStorage := &mock.MockStorage{
		GetNoteByIdFunc: func(ctx context.Context, noteID uuid.UUID) (domain.Note, error) {
			require.Equal(t, id, noteID)
			return expected, nil
		},
	}

	svc := notes.New(log, mockStorage)

	result, err := svc.GetNoteById(ctx, id)
	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestGetNoteById_NotFound(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		GetNoteByIdFunc: func(ctx context.Context, id uuid.UUID) (domain.Note, error) {
			return domain.Note{}, storageerrors.ErrNotFound
		},
	}

	svc := notes.New(log, mockStorage)

	_, err := svc.GetNoteById(ctx, uuid.New())
	require.Error(t, err)
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestGetNoteById_Error(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		GetNoteByIdFunc: func(ctx context.Context, id uuid.UUID) (domain.Note, error) {
			return domain.Note{}, errors.New("db error")
		},
	}

	svc := notes.New(log, mockStorage)

	_, err := svc.GetNoteById(ctx, uuid.New())
	require.Error(t, err)
}

func TestInsert_Success(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		InsertFunc: func(ctx context.Context, note domain.Note) error {
			require.Equal(t, "Hello", note.Title)
			return nil
		},
	}

	svc := notes.New(log, mockStorage)

	err := svc.Insert(ctx, domain.Note{Id: uuid.New(), Title: "Hello"})
	require.NoError(t, err)
}

func TestInsert_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		InsertFunc: func(ctx context.Context, note domain.Note) error {
			return storageerrors.ErrAlreadyExists
		},
	}

	svc := notes.New(log, mockStorage)

	err := svc.Insert(ctx, domain.Note{Id: uuid.New()})
	require.Error(t, err)
	require.ErrorIs(t, err, service.ErrAlreadyExists)
}

func TestInsert_Error(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		InsertFunc: func(ctx context.Context, note domain.Note) error {
			return errors.New("db error")
		},
	}

	svc := notes.New(log, mockStorage)

	err := svc.Insert(ctx, domain.Note{Id: uuid.New()})
	require.Error(t, err)
}

func TestUpdate_Success(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	id := uuid.New()

	mockStorage := &mock.MockStorage{
		UpdateFunc: func(ctx context.Context, noteID uuid.UUID, note domain.Note) error {
			require.Equal(t, id, noteID)
			require.Equal(t, "New", note.Title)
			return nil
		},
	}

	svc := notes.New(log, mockStorage)

	err := svc.Update(ctx, id, domain.Note{Id: id, Title: "New"})
	require.NoError(t, err)
}

func TestUpdate_NotFound(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		UpdateFunc: func(ctx context.Context, id uuid.UUID, note domain.Note) error {
			return storageerrors.ErrNotFound
		},
	}

	svc := notes.New(log, mockStorage)

	err := svc.Update(ctx, uuid.New(), domain.Note{Id: uuid.New()})
	require.Error(t, err)
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestUpdate_Error(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		UpdateFunc: func(ctx context.Context, id uuid.UUID, note domain.Note) error {
			return errors.New("db error")
		},
	}

	svc := notes.New(log, mockStorage)

	err := svc.Update(ctx, uuid.New(), domain.Note{Id: uuid.New()})
	require.Error(t, err)
}

func TestDelete_Success(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	id := uuid.New()

	mockStorage := &mock.MockStorage{
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			require.Equal(t, id, id)
			return nil
		},
	}

	svc := notes.New(log, mockStorage)

	err := svc.Delete(ctx, id)
	require.NoError(t, err)
}

func TestDelete_NotFound(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			return storageerrors.ErrNotFound
		},
	}

	svc := notes.New(log, mockStorage)

	err := svc.Delete(ctx, uuid.New())
	require.Error(t, err)
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestDelete_Error(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			return errors.New("db error")
		},
	}

	svc := notes.New(log, mockStorage)

	err := svc.Delete(ctx, uuid.New())
	require.Error(t, err)
}
