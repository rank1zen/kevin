package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

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

func (m AppMode) String() string {
	switch m {
	case AppModeDevelopment:
		return "Developement"
	default:
		return "Production"
	}
}

type App struct {
	Mode AppMode

	handler http.Handler

	conn *pgxpool.Pool // this should be swapped for internal.Store

	riotClient *riot.Client

	datasource *internal.Datasource

	// Address is the http address to run App. If empty, localhost:4001 is
	// used.
	Address string

	Logger *slog.Logger
}

func New(riotAPIKey string, pgConnStr string, opts ...AppOption) *App {
	ctx := context.Background()

	app := App{
		Address:    "localhost:4001",
	}

	var logger *slog.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	if app.Mode == AppModeDevelopment {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	logger = logger.With("env", app.Mode.String())

	app.Logger = logger

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

	frontend := frontend.New(&frontend.Handler{
		Datasource: app.datasource,
		Store:      store,
	}, frontend.WithLogger(logger))

	for _, opt := range opts {
		if err := opt(&app); err != nil {
			panic(err)
		}
	}

	app.handler = frontend

	return &app
}

type AppOption func(*App) error

func WithAddress(addr string) AppOption {
	return func(a *App) error {
		a.Address = addr
		return nil
	}
}

func WithMode(mode AppMode) AppOption {
	return func(a *App) error {
		a.Mode = mode
		return nil
	}
}

func (app *App) Run(ctx context.Context) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	go func() {
		http.ListenAndServe(app.Address, app.handler)
	}()

	app.Logger.Info("start server", "addr", app.Address)

	<-ctx.Done()

	app.Logger.Info("stop server")
}

func (app *App) Close(ctx context.Context) {
	app.conn.Close()
}
