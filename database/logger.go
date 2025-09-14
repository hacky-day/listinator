package database

import (
	"context"
	"log/slog"
	"time"

	gormLogger "gorm.io/gorm/logger"
)

type slogLogger struct{}

func (l *slogLogger) LogMode(_ gormLogger.LogLevel) gormLogger.Interface {
	// We ignore the level. We just want to use the global one
	return l
}

func (l *slogLogger) Info(ctx context.Context, msg string, data ...any) {
	slog.InfoContext(ctx, msg, data...)
}

func (l *slogLogger) Warn(ctx context.Context, msg string, data ...any) {
	slog.WarnContext(ctx, msg, data...)
}

func (l *slogLogger) Error(ctx context.Context, msg string, data ...any) {
	slog.ErrorContext(ctx, msg, data...)
}

func (l *slogLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// Only Log on Debug Level
	if !slog.Default().Enabled(ctx, slog.LevelDebug) {
		return
	}
	sql, rows := fc()
	slog.Debug("Database Trace", "sql", sql, "rows", rows, "error", err, "duration", time.Since(begin))
}
