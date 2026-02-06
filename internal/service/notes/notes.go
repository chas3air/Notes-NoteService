package notes

import (
	"context"
	"errors"
	"fmt"
	"notesservice/internal/models/domain"
	"notesservice/internal/service"
	storageerrors "notesservice/internal/storage"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Storage interface {
	GetNotes(ctx context.Context, offset int, limit int) ([]domain.Note, error)
	GetNotesByUser(ctx context.Context, userId uuid.UUID, offset int, limit int) ([]domain.Note, error)
	GetNoteById(ctx context.Context, id uuid.UUID) (domain.Note, error)
	Insert(ctx context.Context, note domain.Note) error
	Update(ctx context.Context, id uuid.UUID, note domain.Note) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	log     *zap.Logger
	storage Storage
}

func New(log *zap.Logger, storage Storage) *Service {
	return &Service{
		log:     log,
		storage: storage,
	}
}

func (s *Service) GetNotes(ctx context.Context, offset int, limit int) ([]domain.Note, error) {
	const op = "service.notes.GetNotes"
	log := s.log.With(zap.String("op", op))

	notes, err := s.storage.GetNotes(ctx, offset, limit)
	if err != nil {
		log.Error("failed to get notes from storage", zap.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return notes, nil
}

func (s *Service) GetNotesByUser(ctx context.Context, userId uuid.UUID, offset int, limit int) ([]domain.Note, error) {
	const op = "service.notes.GetNotesByUser"
	log := s.log.With(zap.String("op", op))

	notes, err := s.storage.GetNotesByUser(ctx, userId, offset, limit)
	if err != nil {
		log.Error("failed to get notes by user from storage", zap.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return notes, nil
}

func (s *Service) GetNoteById(ctx context.Context, noteId uuid.UUID) (domain.Note, error) {
	const op = "service.notes.GetNoteById"
	log := s.log.With(zap.String("op", op))

	note, err := s.storage.GetNoteById(ctx, noteId)
	if err != nil {
		if errors.Is(err, storageerrors.ErrNotFound) {
			log.Warn("note not found", zap.String("noteId", noteId.String()))
			return domain.Note{}, fmt.Errorf("%s: %w", op, service.ErrNotFound)
		}

		log.Error("failed to get note by id from storage", zap.Error(err))
		return domain.Note{}, fmt.Errorf("%s: %w", op, err)
	}

	return note, nil
}

func (s *Service) Insert(ctx context.Context, note domain.Note) error {
	const op = "service.notes.Insert"
	log := s.log.With(zap.String("op", op))

	err := s.storage.Insert(ctx, note)
	if err != nil {
		if errors.Is(err, storageerrors.ErrAlreadyExists) {
			log.Warn("note already exists", zap.String("noteId", note.Id.String()))
			return fmt.Errorf("%s: %w", op, service.ErrAlreadyExists)
		}

		log.Error("failed to insert note into storage", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, note domain.Note) error {
	const op = "service.notes.Update"
	log := s.log.With(zap.String("op", op))

	err := s.storage.Update(ctx, id, note)
	if err != nil {
		if errors.Is(err, storageerrors.ErrNotFound) {
			log.Warn("note not found for update", zap.String("id", id.String()))
			return fmt.Errorf("%s: %w", op, service.ErrNotFound)
		}

		log.Error("failed to update note in storage", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
func (s *Service) Delete(ctx context.Context, noteId uuid.UUID) error {
	const op = "service.notes.Delete"
	log := s.log.With(zap.String("op", op))

	err := s.storage.Delete(ctx, noteId)
	if err != nil {
		if errors.Is(err, storageerrors.ErrNotFound) {
			log.Warn("note not found for deletion", zap.String("noteId", noteId.String()))
			return fmt.Errorf("%s: %w", op, service.ErrNotFound)
		}

		log.Error("failed to delete note from storage", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
