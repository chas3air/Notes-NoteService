package app

import (
	"context"
	grpcapp "notesservice/internal/app/grpc"
	restapp "notesservice/internal/app/rest"
	"notesservice/internal/models/domain"
	"notesservice/internal/service/notes"

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

type App struct {
	RESTServer *restapp.App
	GRPCServer *grpcapp.App
}

func New(log *zap.Logger, storage Storage, restPort int, grpcPort int) *App {
	notesService := notes.New(log, storage)

	restApp := restapp.New(log, notesService, restPort)
	grpcApp := grpcapp.New(log, notesService, grpcPort)

	return &App{
		RESTServer: restApp,
		GRPCServer: grpcApp,
	}
}
