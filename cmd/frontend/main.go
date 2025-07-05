package main

import (
	"context"
	"flag"
	"os"
)

var (
	mode = flag.String("mode", "s", "h")
	conn = flag.String("conn", os.Getenv("KEVIN_POSTGRES_CONN"), "postgres connection string")
	// have some command flags here
)

func main() {
	flag.Parse()

	ctx := context.Background()

	app := New(os.Getenv("KEVIN_RIOT_API_KEY"), *conn, WithMode(AppModeDevelopment))

	app.Run(ctx)
}
