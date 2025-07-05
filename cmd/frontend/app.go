package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/riot"
)

type AppMode int

const (
	AppModeProduction AppMode = iota
	AppModeDevelopment
)

type App struct {
	mode AppMode

	handler http.Handler

	conn *pgxpool.Pool // this should be swapped for internal.Store

	riotClient *riot.Client

	datasource *internal.Datasource

	address string
}

func (app *App) Run(ctx context.Context) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	switch app.mode {
	case AppModeDevelopment:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	go func() {
		http.ListenAndServe(app.address, app.handler)
	}()

	slog.Info("starting server", "addr", app.address)

	<-ctx.Done()

	slog.Info("stopping server")

}

type logger struct {
	handler http.Handler
}

func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	logger := slog.Default()

	logger.With(
		slog.Group("request",
			slog.String("method", r.Method),
			slog.Any("url", &r.URL),
		),
	)

	l.handler.ServeHTTP(w, r)

	slog.Info(fmt.Sprintf("%s %s %v", r.Method, r.URL.Path, time.Since(start)))
}

func (app *App) Close(ctx context.Context) {
	app.conn.Close()
}

func New(riotAPIKey string, pgConnStr string, opts ...AppOption) *App {
	ctx := context.Background()

	app := App{
		address:    "localhost:4001",
	}

	app.riotClient = riot.NewClient(riotAPIKey)


	pool, err := pgxpool.New(ctx, pgConnStr)
	if err != nil {
		panic(err)
	}

	app.conn = pool

	conn, err := app.conn.Acquire(ctx)
	if err != nil {
		panic(err)
	}
	defer conn.Release()

	// always migrate for now
	m, err := migrate.NewMigrator(ctx, conn.Conn(), "public.schema_version")
	if err != nil {
		panic(err)
	}

	if err := m.LoadMigrations(os.DirFS("./migrations")); err != nil {
		panic(fmt.Errorf("loading migrations: %w", err))
	}

	if err = m.Migrate(ctx); err != nil {
		panic(err)
	}

	store := postgres.NewStore(app.conn)

	app.datasource = internal.NewDatasource(app.riotClient, store)

	router := http.NewServeMux()

	frontend := frontend.New(store, app.datasource)

	router.Handle("/", frontend)

	app.handler = &logger{router}

	for _, opt := range opts {
		if err := opt(&app); err != nil {
			panic(err)
		}
	}

	return &app
}

type AppOption func(*App) error

func WithMode(mode AppMode) AppOption {
	return func(a *App) error {
		a.mode = mode
		return nil
	}
}
