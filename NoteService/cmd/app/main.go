// @title           User Service API
// @version         1.0
// @description     REST and gRPC API for managing users and outbox events
// @host            localhost:8080
// @BasePath        /api/v1
package main

import (
	"notesservice/internal/app"
	"notesservice/internal/storage/postgres"
	"notesservice/pkg/config"
	"notesservice/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log, err := logger.Setup(cfg.Env)
	if err != nil {
		panic(err)
	}

	log.Info("application trying to setting up", zap.Any("config", cfg))

	storage, close, err := postgres.New(log, cfg.Postgres.DSN())

	application := app.New(log, storage, cfg.Rest.Port, cfg.Grpc.Port)

	go func() {
		if err := application.RESTServer.Run(); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err := application.GRPCServer.Run(); err != nil {
			panic(err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)
	<-done

	application.RESTServer.Shutdown()
	application.GRPCServer.Stop()
	close()
}
