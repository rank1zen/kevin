package store_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/rank1zen/kevin/internal/pgtestcontainer"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/store"
)

var DefaultPGInstance *pgtestcontainer.PGInstance

func TestMain(t *testing.M) {
	ctx := context.Background()

	DefaultPGInstance = pgtestcontainer.NewPGInstance(context.Background(), "../../migrations")

	code := t.Run()

	if err := DefaultPGInstance.Terminate(ctx); err != nil {
		log.Fatalf("terminating: %s", err)
	}

	os.Exit(code)
}

func setupStore(ctx context.Context, t testing.TB) *store.Store {
	conn := DefaultPGInstance.SetupConn(ctx, t)
	return &store.Store{Pool: conn}
}

var T1OKGOODYESNA1PUUID = riot.PUUID("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")
