package grpcapp

import (
	"fmt"
	"net"
	"notesservice/internal/handlers"
	"notesservice/internal/handlers/grpc/notes"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	log        *zap.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *zap.Logger, service handlers.Service, port int) *App {
	gRPCServer := grpc.NewServer()

	notes.Register(gRPCServer, log, service)
	reflection.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "grpc.App.Run"
	a.log.Info("starting gRPC server", zap.String("op", op), zap.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		a.log.Error("failed to listen", zap.String("op", op), zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(zap.String("op", op)).Info("stopping gRPC server")
	a.gRPCServer.GracefulStop()
}
