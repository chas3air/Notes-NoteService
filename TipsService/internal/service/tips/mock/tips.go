package mock

import (
	"context"
	"errors"
	"tipsservice/internal/models/domain"

	"github.com/google/uuid"
)

type MockStorage struct {
	GetTipsFunc       func(ctx context.Context, offset int, limit int) ([]domain.Tip, error)
	GetTipsByUserFunc func(ctx context.Context, userId uuid.UUID, offset int, limit int) ([]domain.Tip, error)
	GetTipByIdFunc    func(ctx context.Context, id uuid.UUID) (domain.Tip, error)
	InsertFunc        func(ctx context.Context, tip domain.Tip) error
	UpdateFunc        func(ctx context.Context, id uuid.UUID, tip domain.Tip) error
	DeleteFunc        func(ctx context.Context, id uuid.UUID) error
}

func (m *MockStorage) GetTips(ctx context.Context, offset int, limit int) ([]domain.Tip, error) {
	if m.GetTipsFunc == nil {
		return nil, errors.New("GetTips not implemented")
	}
	return m.GetTipsFunc(ctx, offset, limit)
}

func (m *MockStorage) GetTipsByUser(ctx context.Context, userId uuid.UUID, offset int, limit int) ([]domain.Tip, error) {
	if m.GetTipsByUserFunc == nil {
		return nil, errors.New("GetTipsByUser not implemented")
	}
	return m.GetTipsByUserFunc(ctx, userId, offset, limit)
}

func (m *MockStorage) GetTipById(ctx context.Context, id uuid.UUID) (domain.Tip, error) {
	if m.GetTipByIdFunc == nil {
		return domain.Tip{}, errors.New("GetTipById not implemented")
	}
	return m.GetTipByIdFunc(ctx, id)
}

func (m *MockStorage) Insert(ctx context.Context, tip domain.Tip) error {
	if m.InsertFunc == nil {
		return errors.New("Insert not implemented")
	}
	return m.InsertFunc(ctx, tip)
}

func (m *MockStorage) Update(ctx context.Context, id uuid.UUID, tip domain.Tip) error {
	if m.UpdateFunc == nil {
		return errors.New("Update not implemented")
	}
	return m.UpdateFunc(ctx, id, tip)
}

func (m *MockStorage) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc == nil {
		return errors.New("Delete not implemented")
	}
	return m.DeleteFunc(ctx, id)
}
