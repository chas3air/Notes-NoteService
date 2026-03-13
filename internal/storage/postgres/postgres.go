package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"notesservice/internal/models/domain"
	storageerrors "notesservice/internal/storage"
	"os"
	"path/filepath"

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

	tmpLog.Info("database connection established")

	wd, _ := os.Getwd()
	migratiomPath := filepath.Join(wd, "migrations", "postgres")
	if err := ApplyMigrations(db, migratiomPath); err != nil {
		tmpLog.Error("failed to apply migrations", zap.Error(err))
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

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

func (s *Storage) GetNotes(ctx context.Context, offset int, limit int) ([]domain.Note, error) {
	const op = "storage.postgres.GetNotes"
	log := s.log.With(zap.String("op", op))

	query := `SELECT id, user_id, title, content, created_at
			  FROM notes
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

	var notes = make([]domain.Note, 0, limit)
	var tmp domain.Note

	for rows.Next() {
		if err := rows.Scan(&tmp.Id, &tmp.UserId, &tmp.Title, &tmp.Content, &tmp.CreatedAt); err != nil {
			log.Error("failed to scan row", zap.Error(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		notes = append(notes, tmp)
	}

	return notes, nil
}

func (s *Storage) GetNotesByUser(ctx context.Context, userId uuid.UUID, offset int, limit int) ([]domain.Note, error) {
	const op = "storage.postgres.GetNotesByUser"
	log := s.log.With(zap.String("op", op))

	query := `SELECT id, user_id, title, content, created_at
			  FROM notes
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

	var notes []domain.Note

	for rows.Next() {
		var tmp domain.Note

		if err := rows.Scan(&tmp.Id, &tmp.UserId, &tmp.Title, &tmp.Content, &tmp.CreatedAt); err != nil {
			log.Error("failed to scan row", zap.Error(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		notes = append(notes, tmp)
	}

	return notes, nil
}

func (s *Storage) GetNoteById(ctx context.Context, id uuid.UUID) (domain.Note, error) {
	const op = "storage.postgres.GetNoteById"
	log := s.log.With(zap.String("op", op))

	query := `SELECT id, user_id, title, content, created_at, is_private
			  FROM notes
			  WHERE id = $1;`

	row := s.db.QueryRowContext(ctx, query, id)

	var tmp domain.Note
	err := row.Scan(&tmp.Id, &tmp.UserId, &tmp.Title, &tmp.Content, &tmp.CreatedAt, &tmp.IsPrivate)
	if err == sql.ErrNoRows {
		log.Warn("note not found", zap.String("note_id", id.String()))
		return domain.Note{}, fmt.Errorf("%s: %w", op, storageerrors.ErrNotFound)
	}
	if err != nil {
		log.Error("failed to scan row", zap.Error(err))
		return domain.Note{}, fmt.Errorf("%s: %w", op, err)
	}

	return tmp, nil
}

func (s *Storage) Insert(ctx context.Context, note domain.Note) error {
	const op = "storage.postgres.Insert"
	log := s.log.With(zap.String("op", op))

	query := `INSERT INTO notes (id, user_id, title, content, created_at, is_private)
			  VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := s.db.ExecContext(ctx, query, note.Id, note.UserId, note.Title, note.Content, note.CreatedAt, note.IsPrivate)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				log.Warn("duplicate note", zap.Error(err))
				return fmt.Errorf("already exists: %w", storageerrors.ErrAlreadyExists)
			}
		}

		log.Error("failed to execute insert", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Update(ctx context.Context, id uuid.UUID, note domain.Note) error {
	const op = "storage.postgres.Update"
	log := s.log.With(zap.String("op", op))

	query := `UPDATE notes
			  SET title = $1, content = $2, is_private = $3
			  WHERE id = $4;`

	res, err := s.db.ExecContext(ctx, query, note.Title, note.Content, note.IsPrivate, id)
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
		log.Warn("note not found for update", zap.String("note_id", id.String()))
		return fmt.Errorf("%s: %w", op, storageerrors.ErrNotFound)
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "storage.postgres.Delete"
	log := s.log.With(zap.String("op", op))

	query := `DELETE FROM notes WHERE id = $1;`
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
		log.Warn("note not found for deletion", zap.String("note_id", id.String()))
		return fmt.Errorf("%s: %w", op, storageerrors.ErrNotFound)
	}

	return nil
}
