package logger

import (
	"go.uber.org/zap"
)

func NewLogger(logLvl string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(logLvl)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
