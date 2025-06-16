package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
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
	mode AppMode

	handler http.Handler

	conn *pgxpool.Pool // this should be swapped for internal.Store

	riotClient *riot.Client

	datasource *internal.Datasource

	// address is the http address to run App. If empty, localhost:4001 is
	// used.
	address string

	logger *slog.Logger
}

func New(riotAPIKey string, pgConnStr string, opts ...AppOption) *App {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, pgConnStr)
	if err != nil {
		panic(err)
	}

	// defaults
	app := App{
		mode:       AppModeProduction,
		handler:    nil,
		conn:       pool,
		riotClient: riot.NewClient(riotAPIKey),
		datasource: &internal.Datasource{},
		address:    "0.0.0.0:4001",
		logger:     slog.Default(),
	}

	for _, f := range opts {
		f(&app)
	}

	var logger *slog.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	if app.mode == AppModeDevelopment {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	logger = logger.With("env", app.mode.String())
	app.logger = logger

	conn, err := app.conn.Acquire(ctx)
	if err != nil {
		panic(err)
	}
	defer conn.Release()

	store := postgres.NewStore(app.conn)

	app.datasource = internal.NewDatasource(app.riotClient, store)

	frontend := frontend.New(
		&frontend.Handler{Datasource: app.datasource,},
		frontend.WithLogger(logger),
	)

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
		a.address = addr
		return nil
	}
}

func WithMode(mode AppMode) AppOption {
	return func(a *App) error {
		a.mode = mode
		return nil
	}
}

func (app *App) Run(ctx context.Context) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	go func() {
		http.ListenAndServe(app.address, app.handler)
	}()

	app.logger.Info("start server", "addr", app.address)

	<-ctx.Done()

	app.logger.Info("stop server")
}

func (app *App) Close(ctx context.Context) {
	app.conn.Close()
}
