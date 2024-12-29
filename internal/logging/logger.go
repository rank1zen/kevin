package logging

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once   sync.Once
	logger *zap.Logger
)

type ctxKey struct{}

func Get() *zap.Logger {
	once.Do(func() {
		stdout := zapcore.AddSync(os.Stdout)

		var level zapcore.Level
		levelEnv := os.Getenv("LOG_LEVEL")
		if levelEnv != "" {
			levelFromEnv, err := zapcore.ParseLevel(levelEnv)
			if err != nil {
				log.Println(fmt.Errorf("invalid level, defaulting to INFO: %w", err))
				level = zap.InfoLevel
			} else {
				level = levelFromEnv
			}
		} else {
			level = zap.InfoLevel
		}

		logLevel := zap.NewAtomicLevelAt(level)

		devCfg := zap.NewDevelopmentEncoderConfig()
		devCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(devCfg)

		buildInfo, _ := debug.ReadBuildInfo()

		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, stdout, logLevel).
				With([]zapcore.Field{
					zap.String("go_version", buildInfo.GoVersion),
				}),
		)

		logger = zap.New(core)
	})

	return logger
}

func NewContext(ctx context.Context, fields ...zap.Field) context.Context {
	return context.WithValue(ctx, ctxKey{}, FromContext(ctx).With(fields...))
}

func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	} else if l := logger; l != nil {
		return l
	}

	return Get()
}

func WithContext(ctx context.Context, lg *zap.Logger) context.Context {
	if lp, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		if lp == lg {
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, lg)
}

func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	logger := FromContext(ctx).With(fields...)
	return WithContext(ctx, logger)
}
