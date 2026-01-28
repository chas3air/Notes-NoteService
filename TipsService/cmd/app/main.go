package main

import (
	"tipsservice/pkg/config"
	"tipsservice/pkg/logger"

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

	// psotgres := ...

	// application := ...

	// application starting
	// grpc.Starting
	// rest.Starting

	// application stopping
	// postgres.Stop()
}
