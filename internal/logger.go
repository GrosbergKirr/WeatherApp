package internal

import (
	"errors"
	"log/slog"
	"os"
)

func SetupLogger(cfg *Config) (*slog.Logger, error) {
	var log *slog.Logger
	switch {
	case cfg.LogLevel == "Debug":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		log.Info("Logger set success. LEVEL DEBUG")
		return log, nil
	case cfg.LogLevel == "Info":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		log.Info("Logger set success. LEVEL INFO")
		return log, nil
	case cfg.LogLevel == "Warn":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
		log.Info("Logger set success. LEVEL WARN")
		return log, nil
	case cfg.LogLevel == "Error":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		log.Info("Logger set success. LEVEL ERROR")
		return log, nil
	}

	return nil, errors.New("incorrect logger level")
}
