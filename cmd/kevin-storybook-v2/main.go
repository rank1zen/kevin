package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/rank1zen/kevin/internal/frontend"
)

func main() {
	config := Config{}

	ctx := context.Background()

	if err := config.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

	os.Exit(0)
}

type Config struct{}

func (c *Config) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	router := http.NewServeMux()

	router.Handle("GET /static/", http.FileServer(http.FS(frontend.StaticAssets)))

	router.Handle("/assets/js/", http.StripPrefix("/assets/js/", http.FileServer(http.Dir("./assets/js"))))

	registerProfileRoutes(router)

	go func() {
		err := http.ListenAndServe(":3001", router)
		slog.Error("error starting server", "err", err)
	}()

	slog.Info("Listening on :3001")

	<-ctx.Done()

	slog.Info("shutting down server")

	return nil
}
