package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/rank1zen/yujin/internal/db"
	"github.com/rank1zen/yujin/internal/logging"
	"github.com/rank1zen/yujin/internal/ui"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	var port int
	port = 4001
	if os.Getenv("YUJIN_PORT") != "" {
		port, _ = strconv.Atoi(os.Getenv("YUJIN_PORT"))
	}

	logger := logging.Get()
	defer logger.Sync()

	db, err := db.NewDB(ctx, os.Getenv("YUJIN_POSTGRES_POOL_URL"))
	if err != nil {
		logger.Sugar().Fatalf("connecting to db: %v", err)
	}
	defer db.Close()

	ui := ui.Routes(db)
	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", port), ui)
	}()

	logger.Sugar().Infof("started server on %v", zap.Int("port", port))

	<-ctx.Done()

	logger.Sugar().Info("shutting down")
}
