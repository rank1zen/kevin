package pgtestcontainer_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/rank1zen/kevin/internal/pgtestcontainer"
)

var DefaultPGInstance *pgtestcontainer.PGInstance

func TestMain(t *testing.M) {
	ctx := context.Background()

	DefaultPGInstance = pgtestcontainer.NewPGInstance(context.Background())

	code := t.Run()

	if err := DefaultPGInstance.Terminate(ctx); err != nil {
		log.Fatalf("terminating: %s", err)
	}

	os.Exit(code)
}

func TestConnect(t *testing.T) {
	_ = DefaultPGInstance.SetupConn(context.Background(), t)
}
