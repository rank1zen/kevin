package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
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

	router := NewRouter()
	go func() {
		err := http.ListenAndServe(":3001", router)
		slog.Error("error starting server", "err", err)
	}()

	slog.Info("Listening on :3001")

	<-ctx.Done()

	slog.Info("shutting down server")

	return nil
}

func aaa(x int) *int {
	return &x
}
