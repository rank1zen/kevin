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
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/riot"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	addr := getAddr()

	pool, err := pgxpool.New(ctx, os.Getenv("KEVIN_DATABASE_URL"))
	if err != nil {
		slog.Error("connecting to postgres", "err", err)
		os.Exit(1)
	}

	defer pool.Close()
	store := internal.NewStore(pool)
	client := riot.NewClient(riot.WithApiKey(os.Getenv("KEVIN_RIOT_API_KEY")))
	ds := internal.NewDatasource(client, store)

	router := http.NewServeMux()
	router.Handle("/fetch/", http.StripPrefix("/fetch", internal.FetchRoutes(ds)))
	router.Handle("/", frontend.Routes(ds))

	go func() {
		http.ListenAndServe(addr, &logger{router})
	}()

	slog.Info("starting server", "addr", addr)

	<-ctx.Done()

	slog.Info("stopping server")
}

func getAddr() string {
	if addr := os.Getenv("YUJIN_PORT"); addr != "" {
		return addr
	} else {
		return ":4001"
	}
}

type logger struct {
	handler http.Handler
}

func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	slog.Info(fmt.Sprintf("%s %s %v", r.Method, r.URL.Path, time.Since(start)))
}
