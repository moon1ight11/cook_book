package logger

import (
	"log/slog"
	"os"
	"path/filepath"

	"cook_book/internal/config"

	"gopkg.in/natefinch/lumberjack.v2"
)

type appLogger struct {
	*slog.Logger
}

func NewLogger(cfg config.LoggerConfig) (Logger, error) {
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "error":
		level = slog.LevelError
	case "warn":
		level = slog.LevelWarn
	default:
		level = slog.LevelInfo
	}

	if cfg.FilePath != "" {
		logDir := filepath.Dir(cfg.FilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, err
		}
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename: cfg.FilePath,
		MaxSize:  cfg.MaxSize,
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	handler := slog.NewTextHandler(lumberjackLogger, opts)

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return &appLogger{Logger: logger}, nil
}
