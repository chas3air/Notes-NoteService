package logger

import (
	"errors"
	"notesservice/pkg/config"

	"go.uber.org/zap"
)

func Setup(env string) (*zap.Logger, error) {
	var cfg zap.Config

	switch env {
	case config.EnvLocal, config.EnvDev:
		cfg = zap.NewDevelopmentConfig()
	case config.EnvProd:
		cfg = zap.NewProductionConfig()
	default:
		return nil, errors.New("invalid environment")
	}

	cfg.Level.SetLevel(zap.DebugLevel)

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
