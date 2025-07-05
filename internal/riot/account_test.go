package riot_test

import (
	"context"
	"os"
	"testing"

	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAccountByRiotID(t *testing.T) {
	client := riot.NewClient(os.Getenv("KEVIN_RIOT_API_KEY"))

	ctx := context.Background()

	t.Run(
		"account is independent of route",
		func(t *testing.T) {
			americas, err := client.Account.GetAccountByRiotID(ctx, riot.ContinentAmericas, "orrange", "NA1")
			require.NoError(t, err)

			europe, err := client.Account.GetAccountByRiotID(ctx, riot.ContinentEurope, "orrange", "NA1")
			require.NoError(t, err)

			asia, err := client.Account.GetAccountByRiotID(ctx, riot.ContinentAsia, "orrange", "NA1")
			require.NoError(t, err)

			assert.Equal(t, americas.PUUID, europe.PUUID)
			assert.Equal(t, americas.PUUID, asia.PUUID)
		},
	)
}
