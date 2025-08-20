package internal_test

import (
	"context"
	"os"
	"testing"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
)

func TestDatasource_GetMatchDetail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	riotClient := riot.NewClient(getEnvRiotAPIKey())

	pgPool := DefaultPGInstance.SetupConn(ctx, t)
	store := postgres.NewStore(pgPool)

	datasource := internal.NewDatasource(riotClient, store)


	match, err := datasource.GetMatchDetail(ctx, riot.RegionNA1, "NA1_5346312088")
	if assert.NoError(t, err) {
		T1 := findT1(t, match)

		for _, tc := range []struct {
			Name             string
			Expected, Actual any
		}{
			{
				Name:     "expects correct version",
				Expected: "15.16.702.7993",
				Actual:   match.Version,
			},
			{
				Name:     "expects correct summoner spells",
				Expected: [2]int{4,12},
				Actual:   T1.SummonerIDs,
			},
		} {
			t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
		}
	}
}

func TestDatasource_GetProfileDetail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	riotClient := riot.NewClient(getEnvRiotAPIKey())

	pgPool := DefaultPGInstance.SetupConn(ctx, t)
	store := postgres.NewStore(pgPool)

	datasource := internal.NewDatasource(riotClient, store)

	_, err := datasource.GetProfileDetail(ctx, riot.RegionNA1, T1OKGOODYESNA1PUUID)
	assert.NoError(t, err)
}

func getEnvRiotAPIKey() string {
	return os.Getenv("KEVIN_RIOT_API_KEY")
}
