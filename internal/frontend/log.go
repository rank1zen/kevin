package frontend

import (
	"context"
	"log/slog"
)

type ctxKey struct{}

func LoggerNewContext(parent context.Context, logger *slog.Logger) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if lp, ok := parent.Value(ctxKey{}).(*slog.Logger); ok {
		// if parent already has the same loggger
		if lp == logger {
			return parent
		}
	}

	return context.WithValue(parent, ctxKey{}, logger)
}

func LoggerFromContext(parent context.Context) *slog.Logger {
	if logger, ok := parent.Value(ctxKey{}).(*slog.Logger); ok {
		return logger
	}

	return slog.Default()
}
