package service_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/service"
	"github.com/stretchr/testify/require"
)

func SetupDatasource(ctx context.Context, t testing.TB) *service.Service {
	pool := DefaultPGInstance.SetupConn(ctx, t)

	client := riot.NewClient("RGAPI-4bd06d26-4a55-4c92-adf5-a0aba99a5e35")

	store := postgres.NewStore(pool)

	return service.NewService(client, store, pool)
}

var T1OKGOODYESNA1PUUID = riot.PUUID("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg")

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

var DefaultPGInstance *postgres.PGInstance

func TestMain(t *testing.M) {
	ctx := context.Background()

	DefaultPGInstance = postgres.NewPGInstance(context.Background(), "../../migrations/")

	code := t.Run()

	if err := DefaultPGInstance.Terminate(ctx); err != nil {
		log.Fatalf("terminating: %s", err)
	}

	os.Exit(code)
}
