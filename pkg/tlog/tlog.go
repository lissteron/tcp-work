package tlog

import (
	"context"
	"log/slog"
)

type Logger interface {
	Infow(ctx context.Context, format string, args ...interface{})
	Errorw(ctx context.Context, format string, args ...interface{})
}

type LoggerPKG struct{}

func New() (*LoggerPKG, error) {
	return &LoggerPKG{}, nil
}

//nolint:goprintffuncname // lint mistake
func (l *LoggerPKG) Infow(ctx context.Context, format string, args ...interface{}) {
	slog.InfoContext(ctx, format, args...)
}

//nolint:goprintffuncname // lint mistake
func (l *LoggerPKG) Errorw(ctx context.Context, format string, args ...interface{}) {
	slog.ErrorContext(ctx, format, args...)
}
