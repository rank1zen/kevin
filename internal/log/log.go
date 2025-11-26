package log

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func InitLogger(env string) {
	var handler slog.Handler
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}
	Logger = slog.New(handler)
	slog.SetDefault(Logger) // Set as default logger for standard library
}
