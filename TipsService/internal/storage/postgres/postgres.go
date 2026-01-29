package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"tipsservice/internal/models/domain"
	storageerrors "tipsservice/internal/storage"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pressly/goose"
	"go.uber.org/zap"
)

type Storage struct {
	log *zap.Logger
	db  *sql.DB
}

func New(log *zap.Logger, conn string) (*Storage, func() error, error) {
	const op = "storage.postgres.New"
	tmpLog := log.With(zap.String("op", op))

	db, err := sql.Open("postgres", conn)
	if err != nil {
		tmpLog.Error("failed to open database connection", zap.Error(err))
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		tmpLog.Error("failed to ping database", zap.Error(err))
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	tmpLog.Info("database connection established")

	return &Storage{
		log: log,
		db:  db,
	}, db.Close, nil
}

func NewWithDB(log *zap.Logger, db *sql.DB) (*Storage, error) {
	return &Storage{
		log: log,
		db:  db,
	}, nil
}

func ApplyMigrations(db *sql.DB, migrationPath string) error {
	const op = "storage.postgres.apply"
	if err := goose.Up(db, migrationPath); err != nil {
		goose.SetLogger(nil)
		return goose.Up(db, migrationPath)
	}

	return nil
}

func (s *Storage) GetTips(ctx context.Context, offset int, limit int) ([]domain.Tip, error) {
	const op = "storage.postgres.GetTips"
	log := s.log.With(zap.String("op", op))

	query := `SELECT id, user_id, title, content, created_at
			  FROM tips
			  WHERE is_private = false
			  ORDER BY created_at DESC
			  OFFSET $1
			  LIMIT $2;`

	rows, err := s.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		log.Error("failed to execute query", zap.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var tips = make([]domain.Tip, 0, limit)
	var tmp domain.Tip

	for rows.Next() {
		if err := rows.Scan(&tmp.Id, &tmp.UserId, &tmp.Title, &tmp.Content, &tmp.CreatedAt); err != nil {
			log.Error("failed to scan row", zap.Error(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		tips = append(tips, tmp)
	}

	return tips, nil
}

func (s *Storage) GetTipsByUser(ctx context.Context, userId uuid.UUID, offset int, limit int) ([]domain.Tip, error) {
	const op = "storage.postgres.GetTipsByUser"
	log := s.log.With(zap.String("op", op))

	query := `SELECT id, user_id, title, content, created_at
			  FROM tips
			  WHERE user_id = $1
			  ORDER BY created_at DESC
			  LIMIT $2
			  OFFSET $3;`

	rows, err := s.db.QueryContext(ctx, query, userId, limit, offset)
	if err != nil {
		log.Error("failed to execute query", zap.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var tips []domain.Tip

	for rows.Next() {
		var tmp domain.Tip

		if err := rows.Scan(&tmp.Id, &tmp.UserId, &tmp.Title, &tmp.Content, &tmp.CreatedAt); err != nil {
			log.Error("failed to scan row", zap.Error(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		tips = append(tips, tmp)
	}

	return tips, nil
}

func (s *Storage) GetTipById(ctx context.Context, id uuid.UUID) (domain.Tip, error) {
	const op = "storage.postgres.GetTipById"
	log := s.log.With(zap.String("op", op))

	query := `SELECT id, user_id, title, content, created_at, is_private
			  FROM tips
			  WHERE id = $1;`

	row := s.db.QueryRowContext(ctx, query, id)

	var tmp domain.Tip
	err := row.Scan(&tmp.Id, &tmp.UserId, &tmp.Title, &tmp.Content, &tmp.CreatedAt, &tmp.IsPrivate)
	if err == sql.ErrNoRows {
		log.Warn("tip not found", zap.String("tip_id", id.String()))
		return domain.Tip{}, fmt.Errorf("%s: %w", op, storageerrors.ErrNotFound)
	}
	if err != nil {
		log.Error("failed to scan row", zap.Error(err))
		return domain.Tip{}, fmt.Errorf("%s: %w", op, err)
	}

	return tmp, nil
}

func (s *Storage) Insert(ctx context.Context, tip domain.Tip) error {
	const op = "storage.postgres.Insert"
	log := s.log.With(zap.String("op", op))

	query := `INSERT INTO tips (id, user_id, title, content, created_at, is_private)
			  VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := s.db.ExecContext(ctx, query, tip.Id, tip.UserId, tip.Title, tip.Content, tip.CreatedAt, tip.IsPrivate)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				log.Warn("duplicate tip", zap.Error(err))
				return fmt.Errorf("already exists: %w", storageerrors.ErrAlreadyExists)
			}
		}

		log.Error("failed to execute insert", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Update(ctx context.Context, id uuid.UUID, tip domain.Tip) error {
	const op = "storage.postgres.Update"
	log := s.log.With(zap.String("op", op))

	query := `UPDATE tips
			  SET title = $1, content = $2, is_private = $3
			  WHERE id = $4;`

	res, err := s.db.ExecContext(ctx, query, tip.Title, tip.Content, tip.IsPrivate, id)
	if err != nil {
		log.Error("failed to execute update", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Error("failed to get rows affected", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		log.Warn("tip not found for update", zap.String("tip_id", id.String()))
		return fmt.Errorf("%s: %w", op, storageerrors.ErrNotFound)
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "storage.postgres.Delete"
	log := s.log.With(zap.String("op", op))

	query := `DELETE FROM tips WHERE id = $1;`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Error("failed to execute delete", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("failed to get rows affected", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		log.Warn("tip not found for deletion", zap.String("tip_id", id.String()))
		return fmt.Errorf("%s: %w", op, storageerrors.ErrNotFound)
	}

	return nil
}
