package internal_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/stretchr/testify/require"
)

var DefaultPGInstance *postgres.PGInstance

func TestMain(t *testing.M) {
	ctx := context.Background()

	DefaultPGInstance = postgres.NewPGInstance(context.Background(), "../migrations/")

	code := t.Run()

	if err := DefaultPGInstance.Terminate(ctx); err != nil {
		log.Fatalf("terminating: %s", err)
	}

	os.Exit(code)
}

var T1OKGOODYESNA1PUUID = internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")

func findT1(tb testing.TB, match internal.MatchDetail) internal.ParticipantDetail {
	var actualParticipant *internal.ParticipantDetail
	for _, p := range match.Participants {
		if p.PUUID == T1OKGOODYESNA1PUUID {
			actualParticipant = &p
		}
	}

	require.NotNil(tb, actualParticipant)
	return *actualParticipant
}
