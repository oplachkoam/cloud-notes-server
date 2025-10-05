package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"cloud-notes/internal/config"
)

type logger struct {
	logger *slog.Logger
}

func (l *logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

func (l *logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

func (l *logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}

func (l *logger) With(args ...any) Logger {
	l.logger.With(args...)
	return l
}

func Load(cfg *config.Logger) (Logger, error) {
	var level slog.Level
	switch strings.Trim(strings.ToLower(cfg.Level), " \n\t") {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		return nil, fmt.Errorf("unknown log level: %s", cfg.Level)
	}

	var output io.Writer
	switch strings.Trim(strings.ToLower(cfg.Output), " \n\t") {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	case "discard":
		output = io.Discard
	default:
		return nil, fmt.Errorf("unknown output: %s", cfg.Output)
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	switch strings.Trim(strings.ToLower(cfg.Format), " \n\t") {
	case "text":
		return &logger{logger: slog.New(slog.NewTextHandler(output, opts))}, nil
	case "json":
		return &logger{logger: slog.New(slog.NewJSONHandler(output, opts))}, nil
	default:
		return nil, fmt.Errorf("unknown log format: %s", cfg.Format)
	}
}

func MustLoad(cfg *config.Logger) Logger {
	l, err := Load(cfg)
	if err != nil {
		panic(err)
	}

	return l
}
