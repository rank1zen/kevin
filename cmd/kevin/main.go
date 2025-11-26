package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/kevin/internal/config"
	"github.com/rank1zen/kevin/internal/frontend/server"
	"github.com/rank1zen/kevin/internal/log"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/service"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Default().Error("failed to load configuration", "err", err)
		os.Exit(1)
	}

	log.InitLogger(cfg.Environment)

	tp, err := initTelemetry(context.Background(), cfg.IsProduction())
	if err != nil {
		slog.Default().Error("failed to initialize OpenTelemetry", "err", err)
		os.Exit(1)
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			slog.Default().Error("Error shutting down tracer provider", "err", err)
		}
	}()

	if err := run(context.Background(), cfg); err != nil {
		slog.Default().Error("application exited with error", "err", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run(ctx context.Context, cfg *config.Config) error {
	pool, err := connectPostgres(ctx, cfg.PostgresConnection)
	if err != nil {
		return errors.New("failed to connect to postgres")
	}
	defer pool.Close()

	store := postgres.NewStore(pool)
	client := riot.NewClient(cfg.RiotAPIKey)
	datasource := service.NewService(client, store)

	srvr := server.New(datasource, server.WithLogger(slog.Default()), server.WithAddress(cfg.Address))

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	go func() {
		err := srvr.Open()
		slog.Default().Error("error starting server", "err", err)
	}()

	slog.Default().Info("server started", "address", cfg.Address, "environment", cfg.Environment)

	<-ctx.Done()

	slog.Default().Info("shutting down server")

	return nil
}

func connectPostgres(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres connection pool: %w", err)
	}

	return pool, nil
}

func initTelemetry(ctx context.Context, isProduction bool) (*trace.TracerProvider, error) {
	var exporter trace.SpanExporter
	var err error

	if isProduction {
		// In a production environment, you would typically use an OTLP exporter
		// to send traces to a collector (e.g., Jaeger, Tempo).
		// For example:
		// exporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
		// }
		// For now, we'll use stdouttrace even in production for simplicity in this exercise.
		// In a real scenario, this would be an OTLP exporter.
		exporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithoutTimestamps(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout trace exporter: %w", err)
		}
	} else {
		// For development, print traces to stdout
		exporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
			stdouttrace.WithoutTimestamps(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout trace exporter: %w", err)
		}
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("kevin-frontend"),
			semconv.ServiceVersion("0.1.0"), // TODO: Get version from build info
		)),
	)
	otel.SetTracerProvider(tp)
	// otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})) // For context propagation

	slog.Default().Info("OpenTelemetry initialized", "exporter", "stdouttrace", "isProduction", isProduction)
	return tp, nil
}
