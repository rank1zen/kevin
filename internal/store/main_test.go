package store_test

import (
	"context"
	"log"
	"testing"

	"github.com/rank1zen/kevin/internal/pgtestcontainer"
)

const ExamplePUUID = "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg"
const ExamplePUUID2 = "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNq1"

var DefaultPGInstance *pgtestcontainer.PGInstance

func TestMain(m *testing.M) {
	DefaultPGInstance = pgtestcontainer.NewPGInstance(context.Background())
	defer func(DefaultPGInstance *pgtestcontainer.PGInstance, ctx context.Context) {
		err := DefaultPGInstance.Terminate(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(DefaultPGInstance, context.Background())
	m.Run()
}
