package web

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// GetDays returns 8 datetimes, representing the past 7 days, including the
// given timestamp. Each date time is at midnight in the provided timezone. The
// slice is ordered in most recent order.
func GetDays(ts time.Time) []time.Time {
	year, month, day := ts.Date()
	ref := time.Date(year, month, day, 0, 0, 0, 0, ts.Location())

	days := []time.Time{}

	for i := range 8 {
		offset := -i + 1
		days = append(days, ref.AddDate(0, 0, offset))
	}

	return days
}

// ParseRiotID returns name and tag from riotID.
func ParseRiotID(riotID string) (name, tag string) {
	index := strings.Index(riotID, "-")
	if index == -1 {
		return riotID, ""
	}

	name = riotID[:index]
	tag = riotID[index+1:]

	return name, tag
}

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

func LogError(r *http.Request, err error) {
	logger := LoggerFromContext(r.Context())
	logger.Error("[http] error", "method", r.Method, "path", r.URL.Path, "error", err)
}
