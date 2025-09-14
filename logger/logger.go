// Package logger handles logging needs
package logger

import (
	"errors"
	"log/slog"
	"os"
)

func Init() error {
	// Log Level
	var level slog.Level
	switch os.Getenv("LISTINATOR_LOG_LEVEL") {
	case "":
		level = slog.LevelInfo
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		return errors.New("unknown log level")
	}

	// Log Type
	options := slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	var h slog.Handler
	switch os.Getenv("LISTINATOR_LOG_TYPE") {
	case "":
		h = slog.NewTextHandler(os.Stderr, &options)
	case "text":
		h = slog.NewTextHandler(os.Stderr, &options)
	case "json":
		h = slog.NewJSONHandler(os.Stderr, &options)
	default:
		return errors.New("unknown log type")
	}

	logger := slog.New(h)

	slog.SetDefault(logger)
	return nil
}
