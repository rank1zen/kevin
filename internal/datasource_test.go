package internal_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SetupDatasource(ctx context.Context, t testing.TB) *internal.Datasource {
	pool := DefaultPGInstance.SetupConn(ctx, t)

	client := riot.NewClient(os.Getenv("KEVIN_RIOT_API_KEY"))

	store := postgres.NewStore(pool)

	return internal.NewDatasource(client, store)
}

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
				Expected: [2]int{4, 12},
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

func TestDatasource_GetMatchHistory(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	ds := SetupDatasource(ctx, t)

	tz, err := time.LoadLocation("America/Toronto")
	require.NoError(t, err)

	startTS := time.Date(2025, time.August, 26, 22, 0, 0, 0, tz)
	endTS := time.Date(2025, time.August, 27, 0, 0, 0, 0, tz)

	matches, err := ds.GetMatchHistory(ctx, riot.RegionNA1, T1OKGOODYESNA1PUUID, startTS, endTS)
	require.NoError(t, err)

	actualIDs := []string{}
	for _, match := range matches {
		actualIDs = append(actualIDs, match.MatchID)
	}

	assert.Equal(t, []string{"NA1_5356002067", "NA1_5355962057"}, actualIDs)
}

func getEnvRiotAPIKey() string {
	return os.Getenv("KEVIN_RIOT_API_KEY")
}
