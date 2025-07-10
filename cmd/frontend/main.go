package main

import (
	"context"
	"flag"
	"os"
)

var (
	mode = os.Getenv("KEVIN_MODE")
	riotAPIKey = os.Getenv("KEVIN_RIOT_API_KEY")
	postgresConn = os.Getenv("KEVIN_POSTGRES_CONN")
)

func main() {
	flag.StringVar(&mode, "m", "", "set app mode")
	flag.StringVar(&mode, "k", "", "set riot api key")
	flag.StringVar(&mode, "c", "", "set postgres connection string")

	flag.Parse()

	ctx := context.Background()

	app := New(riotAPIKey, postgresConn, WithMode(AppModeDevelopment))

	app.Run(ctx)
}
