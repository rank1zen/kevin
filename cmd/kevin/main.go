package main

import (
	"context"
	"os"

	"github.com/rank1zen/kevin/internal/app"
)

func main() {
	ctx := context.Background()
	a := app.New(ctx)
	os.Exit(a.Run(ctx))
}
