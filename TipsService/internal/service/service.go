package service

import (
	"context"
	"errors"
	"tipsservice/internal/models/domain"

	"github.com/google/uuid"
)

var (
	ErrNotFound      = errors.New("tip not found")
	ErrAlreadyExists = errors.New("tip already exists")
)

type QueryStorage interface {
	GetTips(ctx context.Context, offset int, limit int) ([]domain.Tip, error)
	GetTipsByUser(ctx context.Context, userId uuid.UUID, offset int, limit int) ([]domain.Tip, error)
	GetTipById(ctx context.Context, id uuid.UUID) (domain.Tip, error)
}

type CommandStorage interface {
	Insert(ctx context.Context, tip domain.Tip) error
	Update(ctx context.Context, id uuid.UUID, tip domain.Tip) error
	Delete(ctx context.Context, id uuid.UUID) error
}
