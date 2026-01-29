package tips

import (
	"context"
	"errors"
	"fmt"
	"tipsservice/internal/models/domain"
	"tipsservice/internal/service"
	storageerrors "tipsservice/internal/storage"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Storage interface {
	GetTips(ctx context.Context, offset int, limit int) ([]domain.Tip, error)
	GetTipsByUser(ctx context.Context, userId uuid.UUID, offset int, limit int) ([]domain.Tip, error)
	GetTipById(ctx context.Context, id uuid.UUID) (domain.Tip, error)
	Insert(ctx context.Context, tip domain.Tip) error
	Update(ctx context.Context, id uuid.UUID, tip domain.Tip) error
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

func (s *Service) GetTips(ctx context.Context, offset int, limit int) ([]domain.Tip, error) {
	const op = "service.tips.GetTips"
	log := s.log.With(zap.String("op", op))

	tips, err := s.storage.GetTips(ctx, offset, limit)
	if err != nil {
		log.Error("failed to get tips from storage", zap.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tips, nil
}

func (s *Service) GetTipsByUser(ctx context.Context, userId uuid.UUID, offset int, limit int) ([]domain.Tip, error) {
	const op = "service.tips.GetTipsByUser"
	log := s.log.With(zap.String("op", op))

	tips, err := s.storage.GetTipsByUser(ctx, userId, offset, limit)
	if err != nil {
		log.Error("failed to get tips by user from storage", zap.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tips, nil
}

func (s *Service) GetTipById(ctx context.Context, tipId uuid.UUID) (domain.Tip, error) {
	const op = "service.tips.GetTipById"
	log := s.log.With(zap.String("op", op))

	tip, err := s.storage.GetTipById(ctx, tipId)
	if err != nil {
		if errors.Is(err, storageerrors.ErrNotFound) {
			log.Warn("tip not found", zap.String("tipId", tipId.String()))
			return domain.Tip{}, fmt.Errorf("%s: %w", op, service.ErrNotFound)
		}

		log.Error("failed to get tip by id from storage", zap.Error(err))
		return domain.Tip{}, fmt.Errorf("%s: %w", op, err)
	}

	return tip, nil
}

func (s *Service) Insert(ctx context.Context, tip domain.Tip) error {
	const op = "service.tips.Insert"
	log := s.log.With(zap.String("op", op))

	err := s.storage.Insert(ctx, tip)
	if err != nil {
		if errors.Is(err, storageerrors.ErrAlreadyExists) {
			log.Warn("tip already exists", zap.String("tipId", tip.Id.String()))
			return fmt.Errorf("%s: %w", op, service.ErrAlreadyExists)
		}

		log.Error("failed to insert tip into storage", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, tip domain.Tip) error {
	const op = "service.tips.Update"
	log := s.log.With(zap.String("op", op))

	err := s.storage.Update(ctx, tip.Id, tip)
	if err != nil {
		if errors.Is(err, storageerrors.ErrNotFound) {
			log.Warn("tip not found for update", zap.String("tipId", tip.Id.String()))
			return fmt.Errorf("%s: %w", op, service.ErrNotFound)
		}

		log.Error("failed to update tip in storage", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
func (s *Service) Delete(ctx context.Context, tipId uuid.UUID) error {
	const op = "service.tips.Delete"
	log := s.log.With(zap.String("op", op))

	err := s.storage.Delete(ctx, tipId)
	if err != nil {
		if errors.Is(err, storageerrors.ErrNotFound) {
			log.Warn("tip not found for deletion", zap.String("tipId", tipId.String()))
			return fmt.Errorf("%s: %w", op, service.ErrNotFound)
		}

		log.Error("failed to delete tip from storage", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
