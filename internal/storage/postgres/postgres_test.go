package postgres_test

import (
	"context"
	"database/sql"
	"fmt"
	"notesservice/internal/models/domain"
	storageerrors "notesservice/internal/storage"
	"notesservice/internal/storage/postgres"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap/zaptest"
)

func setupPostgres(t *testing.T) (*sql.DB, testcontainers.Container) {
	t.Helper()
	ctx := context.Background()

	postgresC, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "postgres:15",
				ExposedPorts: []string{"5432/tcp"},
				Env: map[string]string{
					"POSTGRES_USER":     "testuser",
					"POSTGRES_PASSWORD": "testpass",
					"POSTGRES_DB":       "testdb",
				},
				WaitingFor: wait.ForListeningPort("5432/tcp"),
			},
			Started: true,
		})
	require.NoError(t, err)

	host, err := postgresC.Host(ctx)
	require.NoError(t, err)

	port, err := postgresC.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=testuser password=testpass dbname=testdb sslmode=disable",
		host, port.Port(),
	)

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	err = postgres.ApplyMigrations(db, "./../../../migrations/postgres")
	require.NoError(t, err)

	return db, postgresC
}

func newStorage(t *testing.T, db *sql.DB) *postgres.Storage {
	t.Helper()
	log := zaptest.NewLogger(t)
	st, err := postgres.NewWithDB(log, db)
	require.NoError(t, err)

	return st
}

func TestInsertAndGetById(t *testing.T) {
	ctx := context.Background()
	db, container := setupPostgres(t)
	defer container.Terminate(ctx)
	defer db.Close()

	storage := newStorage(t, db)

	id := uuid.New()
	note := domain.Note{
		Id:        id,
		UserId:    uuid.New(),
		Title:     "Hello",
		Content:   "World",
		CreatedAt: time.Now(),
		IsPrivate: false,
	}

	err := storage.Insert(ctx, note)
	require.NoError(t, err)

	got, err := storage.GetNoteById(ctx, id)
	require.NoError(t, err)
	require.Equal(t, note.Id, got.Id)
	require.Equal(t, note.Title, got.Title)
	require.Equal(t, note.IsPrivate, got.IsPrivate)
}

func TestGetNoteById_NotFound(t *testing.T) {
	ctx := context.Background()
	db, container := setupPostgres(t)
	defer container.Terminate(ctx)
	defer db.Close()

	storage := newStorage(t, db)

	_, err := storage.GetNoteById(ctx, uuid.New())
	require.Error(t, err)
	require.ErrorIs(t, err, storageerrors.ErrNotFound)
}

func TestGetNotesByUser(t *testing.T) {
	ctx := context.Background()
	db, container := setupPostgres(t)
	defer container.Terminate(ctx)
	defer db.Close()

	storage := newStorage(t, db)

	userID := uuid.New()

	note1 := domain.Note{
		Id:        uuid.New(),
		UserId:    userID,
		Title:     "Note1",
		Content:   "Content1",
		CreatedAt: time.Now(),
		IsPrivate: false,
	}
	note2 := domain.Note{
		Id:        uuid.New(),
		UserId:    userID,
		Title:     "Note2",
		Content:   "Content2",
		CreatedAt: time.Now(),
		IsPrivate: true,
	}

	require.NoError(t, storage.Insert(ctx, note1))
	require.NoError(t, storage.Insert(ctx, note2))

	notes, err := storage.GetNotesByUser(ctx, userID, 0, 2)
	require.NoError(t, err)
	require.Len(t, notes, 2)
}

func TestGetNotes_PublicOnly(t *testing.T) {
	ctx := context.Background()
	db, container := setupPostgres(t)
	defer container.Terminate(ctx)
	defer db.Close()

	storage := newStorage(t, db)

	publicNote := domain.Note{
		Id:        uuid.New(),
		UserId:    uuid.New(),
		Title:     "Public",
		Content:   "Visible",
		CreatedAt: time.Now(),
		IsPrivate: false,
	}
	privateNote := domain.Note{
		Id:        uuid.New(),
		UserId:    uuid.New(),
		Title:     "Private",
		Content:   "Hidden",
		CreatedAt: time.Now(),
		IsPrivate: true,
	}

	require.NoError(t, storage.Insert(ctx, publicNote))
	require.NoError(t, storage.Insert(ctx, privateNote))

	notes, err := storage.GetNotes(ctx, 0, 10)
	require.NoError(t, err)
	require.Len(t, notes, 1)
	require.Equal(t, publicNote.Id, notes[0].Id)
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	db, container := setupPostgres(t)
	defer container.Terminate(ctx)
	defer db.Close()

	storage := newStorage(t, db)

	id := uuid.New()
	original := domain.Note{
		Id:        id,
		UserId:    uuid.New(),
		Title:     "Old",
		Content:   "Old content",
		CreatedAt: time.Now(),
		IsPrivate: false,
	}

	require.NoError(t, storage.Insert(ctx, original))

	update := domain.Note{
		Title:     "New",
		Content:   "New content",
		IsPrivate: true,
	}

	err := storage.Update(ctx, id, update)
	require.NoError(t, err)

	got, err := storage.GetNoteById(ctx, id)
	require.NoError(t, err)
	require.Equal(t, "New", got.Title)
	require.Equal(t, "New content", got.Content)
	require.Equal(t, true, got.IsPrivate)
}

func TestUpdate_NotFound(t *testing.T) {
	ctx := context.Background()
	db, container := setupPostgres(t)
	defer container.Terminate(ctx)
	defer db.Close()

	storage := newStorage(t, db)

	err := storage.Update(ctx, uuid.New(), domain.Note{
		Title: "X",
	})
	require.Error(t, err)
	require.ErrorIs(t, err, storageerrors.ErrNotFound)
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	db, container := setupPostgres(t)
	defer container.Terminate(ctx)
	defer db.Close()

	storage := newStorage(t, db)

	id := uuid.New()
	note := domain.Note{
		Id:        id,
		UserId:    uuid.New(),
		Title:     "ToDelete",
		Content:   "Bye",
		CreatedAt: time.Now(),
		IsPrivate: false,
	}

	require.NoError(t, storage.Insert(ctx, note))

	err := storage.Delete(ctx, id)
	require.NoError(t, err)

	_, err = storage.GetNoteById(ctx, id)
	require.Error(t, err)
	require.ErrorIs(t, err, storageerrors.ErrNotFound)
}

func TestDelete_NotFound(t *testing.T) {
	ctx := context.Background()
	db, container := setupPostgres(t)
	defer container.Terminate(ctx)
	defer db.Close()

	storage := newStorage(t, db)

	err := storage.Delete(ctx, uuid.New())
	require.Error(t, err)
	require.ErrorIs(t, err, storageerrors.ErrNotFound)
}

func TestInsert_Duplicate(t *testing.T) {
	ctx := context.Background()
	db, container := setupPostgres(t)
	defer container.Terminate(ctx)
	defer db.Close()

	storage := newStorage(t, db)

	id := uuid.New()
	note := domain.Note{
		Id:        id,
		UserId:    uuid.New(),
		Title:     "Dup",
		Content:   "Dup",
		CreatedAt: time.Now(),
		IsPrivate: false,
	}

	require.NoError(t, storage.Insert(ctx, note))

	err := storage.Insert(ctx, note)
	require.Error(t, err)
	require.ErrorIs(t, err, storageerrors.ErrAlreadyExists)
}
