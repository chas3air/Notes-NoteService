package tips_test

import (
	"context"
	"errors"
	"testing"
	"tipsservice/internal/models/domain"
	"tipsservice/internal/service"
	"tipsservice/internal/service/tips"
	"tipsservice/internal/service/tips/mock"
	storageerrors "tipsservice/internal/storage"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestGetTips(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		GetTipsFunc: func(ctx context.Context, offset, limit int) ([]domain.Tip, error) {
			return []domain.Tip{}, nil
		},
	}

	svc := tips.New(log, mockStorage)

	result, err := svc.GetTips(ctx, 0, 2)
	require.NoError(t, err)
	require.Len(t, result, 0)
}

func TestGetTips_Error(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		GetTipsFunc: func(ctx context.Context, offset, limit int) ([]domain.Tip, error) {
			return nil, errors.New("db error")
		},
	}

	svc := tips.New(log, mockStorage)

	_, err := svc.GetTips(ctx, 0, 2)
	require.Error(t, err)
}

func TestGetTipsByUser_Success(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	userID := uuid.New()

	mockStorage := &mock.MockStorage{
		GetTipsByUserFunc: func(ctx context.Context, id uuid.UUID, offset, limit int) ([]domain.Tip, error) {
			require.Equal(t, userID, id)
			return []domain.Tip{
				{Id: uuid.New(), UserId: userID, Title: "A"},
			}, nil
		},
	}

	svc := tips.New(log, mockStorage)

	result, err := svc.GetTipsByUser(ctx, userID, 0, 10)
	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestGetTipsByUser_Error(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	userID := uuid.New()

	mockStorage := &mock.MockStorage{
		GetTipsByUserFunc: func(ctx context.Context, id uuid.UUID, offset, limit int) ([]domain.Tip, error) {
			return nil, errors.New("db error")
		},
	}

	svc := tips.New(log, mockStorage)

	_, err := svc.GetTipsByUser(ctx, userID, 0, 10)
	require.Error(t, err)
}

func TestGetTipById_Success(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	id := uuid.New()
	expected := domain.Tip{Id: id, Title: "Test"}

	mockStorage := &mock.MockStorage{
		GetTipByIdFunc: func(ctx context.Context, tipID uuid.UUID) (domain.Tip, error) {
			require.Equal(t, id, tipID)
			return expected, nil
		},
	}

	svc := tips.New(log, mockStorage)

	result, err := svc.GetTipById(ctx, id)
	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestGetTipById_NotFound(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		GetTipByIdFunc: func(ctx context.Context, id uuid.UUID) (domain.Tip, error) {
			return domain.Tip{}, storageerrors.ErrNotFound
		},
	}

	svc := tips.New(log, mockStorage)

	_, err := svc.GetTipById(ctx, uuid.New())
	require.Error(t, err)
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestGetTipById_Error(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		GetTipByIdFunc: func(ctx context.Context, id uuid.UUID) (domain.Tip, error) {
			return domain.Tip{}, errors.New("db error")
		},
	}

	svc := tips.New(log, mockStorage)

	_, err := svc.GetTipById(ctx, uuid.New())
	require.Error(t, err)
}

func TestInsert_Success(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		InsertFunc: func(ctx context.Context, tip domain.Tip) error {
			require.Equal(t, "Hello", tip.Title)
			return nil
		},
	}

	svc := tips.New(log, mockStorage)

	err := svc.Insert(ctx, domain.Tip{Id: uuid.New(), Title: "Hello"})
	require.NoError(t, err)
}

func TestInsert_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		InsertFunc: func(ctx context.Context, tip domain.Tip) error {
			return storageerrors.ErrAlreadyExists
		},
	}

	svc := tips.New(log, mockStorage)

	err := svc.Insert(ctx, domain.Tip{Id: uuid.New()})
	require.Error(t, err)
	require.ErrorIs(t, err, service.ErrAlreadyExists)
}

func TestInsert_Error(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		InsertFunc: func(ctx context.Context, tip domain.Tip) error {
			return errors.New("db error")
		},
	}

	svc := tips.New(log, mockStorage)

	err := svc.Insert(ctx, domain.Tip{Id: uuid.New()})
	require.Error(t, err)
}

func TestUpdate_Success(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	id := uuid.New()

	mockStorage := &mock.MockStorage{
		UpdateFunc: func(ctx context.Context, tipID uuid.UUID, tip domain.Tip) error {
			require.Equal(t, id, tipID)
			require.Equal(t, "New", tip.Title)
			return nil
		},
	}

	svc := tips.New(log, mockStorage)

	err := svc.Update(ctx, domain.Tip{Id: id, Title: "New"})
	require.NoError(t, err)
}

func TestUpdate_NotFound(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		UpdateFunc: func(ctx context.Context, id uuid.UUID, tip domain.Tip) error {
			return storageerrors.ErrNotFound
		},
	}

	svc := tips.New(log, mockStorage)

	err := svc.Update(ctx, domain.Tip{Id: uuid.New()})
	require.Error(t, err)
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestUpdate_Error(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	mockStorage := &mock.MockStorage{
		UpdateFunc: func(ctx context.Context, id uuid.UUID, tip domain.Tip) error {
			return errors.New("db error")
		},
	}

	svc := tips.New(log, mockStorage)

	err := svc.Update(ctx, domain.Tip{Id: uuid.New()})
	require.Error(t, err)
}

func TestDelete_Success(t *testing.T) {
	ctx := context.Background()
	log := zaptest.NewLogger(t)

	id := uuid.New()

	mockStorage := &mock.MockStorage{
		DeleteFunc: func(ctx context.Context, tipID uuid.UUID) error {
			require.Equal(t, id, tipID)
			return nil
		},
	}

	svc := tips.New(log, mockStorage)

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

	svc := tips.New(log, mockStorage)

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

	svc := tips.New(log, mockStorage)

	err := svc.Delete(ctx, uuid.New())
	require.Error(t, err)
}
