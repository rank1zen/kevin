// NOTE: all tests fetch live API, probably need to start mocking
package frontend_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/component/profile"
	"github.com/rank1zen/kevin/internal/ddragon"
	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerUpdateSummoner(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	store := DefaultPGInstance.SetupStore(ctx, t)

	handler := frontend.Handler{internal.NewDatasource(riot.NewClient(os.Getenv("KEVIN_RIOT_API_KEY")), store)}

	err := handler.UpdateSummoner(ctx, riot.RegionNA1, "orrange", "NA1")
	assert.NoError(t, err)
}

func TestHandlerGetLiveMatch(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	store := DefaultPGInstance.SetupStore(ctx, t)

	handler := frontend.Handler{internal.NewDatasource(riot.NewClient(os.Getenv("KEVIN_RIOT_API_KEY")), store)}

	req := frontend.GetLiveMatchRequest{
		Region: riot.RegionNA1,
		PUUID:  internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg"),
	}

	component, err := handler.GetLiveMatch(ctx, req)
	require.NoError(t, err)

	t.Run(
		"expects no live match",
		func(t *testing.T) {
			_, ok := component.(profile.LiveMatchNotFound)
			require.True(t, ok)
		},
	)

	// NOTE: how to write tests for checking live games? Mock?
}

// NOTE: should probably mock so that it is clear what result is expected.
func TestHandlerZGetSummonerChampions(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	store := DefaultPGInstance.SetupStore(ctx, t)

	handler := frontend.Handler{internal.NewDatasource(riot.NewClient(os.Getenv("KEVIN_RIOT_API_KEY")), store)}

	location, err := time.LoadLocation("America/Toronto")
	require.NoError(t, err)

	req := frontend.ZGetSummonerChampionsRequest{
		Region: riot.RegionNA1,
		PUUID:  internal.NewPUUIDFromString("44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg"),
		Week:   time.Date(2025, 7, 5, 0, 0, 0, 0, location),
	}

	component, err := handler.ZGetSummonerChampions(ctx, req)
	require.NoError(t, err)

	actual, ok := component.(profile.SummonerChampionStatContent)
	require.True(t, ok)

	t.Run(
		"expects 4 champions played",
		func(t *testing.T) {
			assert.Equal(t, 4, len(actual.List.Champions))
		},
	)

	require.Greater(t, len(actual.List.Champions), 0)
	require.Equal(t, ddragon.ChampionIllaoiID, actual.List.Champions[0].ChampionWidget.ChampionSprite.ChampionID)
	illaoi := actual.List.Champions[0]

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct win rate for Illaoi",
			Expected: float32(1) / float32(2),
			Actual:   illaoi.WinRate,
		},
		{
			Name:     "expects correct number of games for Illaoi",
			Expected: 6,
			Actual:   illaoi.GamesPlayed,
		},
		{
			Name:     "expects correct KDA for Illaoi",
			Expected: frontend.KDAWidget{Kills: 17 / 6, Deaths: 22 / 6, Assists: 15 / 6},
			Actual:   illaoi.KDAWidget,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}

func TestHandlerGetSummonerPage(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	store := DefaultPGInstance.SetupStore(ctx, t)

	handler := frontend.Handler{internal.NewDatasource(riot.NewClient(os.Getenv("KEVIN_RIOT_API_KEY")), store)}

	t.Run(
		"expects update then returns page",
		func(t *testing.T) {
			component, err := handler.GetSummonerPage(ctx, riot.RegionNA1, "orrange", "NA1")
			require.NoError(t, err)

			expectedPUUID := internal.NewPUUIDFromString("0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q")

			if assert.IsType(t, frontend.SummonerPage{}, component) {
				page, ok := component.(frontend.SummonerPage)
				require.True(t, ok)

				assert.Equal(t, page.Region, riot.RegionNA1)
				assert.Equal(t, page.Name, "orrange")
				assert.Equal(t, page.Tag, "NA1")
				assert.Equal(t, page.PUUID, expectedPUUID)
			}
		},
	)

	t.Run(
		"expects matching fields",
		func(t *testing.T) {
			actual, err := handler.GetSummonerPage(ctx, riot.RegionNA1, "orrange", "NA1")
			require.NoError(t, err)

			expectedPUUID := internal.NewPUUIDFromString("0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q")

			// these are all the fields that we _can_ test.
			expected := frontend.SummonerPage{
				Region:              riot.RegionNA1,
				PUUID:               expectedPUUID,
				Name:                "orrange",
				Tag:                 "NA1",
				LiveMatchLoader:     frontend.LiveMatchModalWindowLoader{Request: frontend.GetLiveMatchRequest{PUUID: expectedPUUID, Region: riot.RegionNA1}},
				MatchHistoryLoaders: []frontend.MatchHistoryListLoader{},
			}

			for i := range 7 {
				expected.MatchHistoryLoaders = append(
					expected.MatchHistoryLoaders,
					frontend.MatchHistoryListLoader{
						Request: frontend.MatchHistoryRequest{
							Region: riot.RegionNA1,
							PUUID:  expectedPUUID,
							Date:   frontend.GetDay(i),
						},
					},
				)
			}

			if assert.IsType(t, expected, actual) {
				page, _ := actual.(frontend.SummonerPage)

				assert.Equal(t, page.Region, riot.RegionNA1)
				assert.Equal(t, page.Name, "orrange")
				assert.Equal(t, page.Tag, "NA1")
				assert.Equal(t, page.PUUID, expectedPUUID)
				assert.Equal(t, expected.LiveMatchLoader, page.LiveMatchLoader)
				assert.Equal(t, expected.MatchHistoryLoaders, page.MatchHistoryLoaders)
			}
		},
	)

	t.Run(
		"expects no summoner page",
		func(t *testing.T) {
			component, err := handler.GetSummonerPage(ctx, riot.RegionKR, "orrange", "KR")
			require.NoError(t, err)

			expected := frontend.NoSummonerPage{
				Region: riot.RegionKR,
				Name:   "orrange",
				Tag:    "KR",
			}

			assert.Equal(t, expected, component)
		},
	)
}

func TestHandlerGetMatchHistory(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	store := DefaultPGInstance.SetupStore(ctx, t)

	handler := frontend.Handler{internal.NewDatasource(riot.NewClient(os.Getenv("KEVIN_RIOT_API_KEY")), store)}

	t.Run(
		"expects 14 matches",
		func(t *testing.T) {
			location, _ := time.LoadLocation("America/Toronto")
			localTime := time.Date(2025, 7, 4, 0, 0, 0, 0, location)

			req := frontend.MatchHistoryRequest{
				Region: riot.RegionNA1,
				PUUID:  "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg",
				Date:   localTime,
			}

			component, err := handler.GetMatchHistory(ctx, req)
			require.NoError(t, err)

			list, ok := component.(frontend.MatchHistoryList)
			if assert.True(t, ok) {
				expectedIDs := []string{
					"NA1_5319611168",
					"NA1_5319592152",
					"NA1_5319579789",
					"NA1_5319551702",
					"NA1_5319526470",
					"NA1_5319509894",
					"NA1_5319489324",
					"NA1_5319337632",
					"NA1_5319319764",
					"NA1_5319307543",
					"NA1_5319296051",
					"NA1_5319287528",
					"NA1_5319275283",
					"NA1_5319263238",
				}

				actualIDs := []string{}
				for _, match := range list.Matches {
					actualIDs = append(actualIDs, match.MatchID)
				}

				assert.Equal(t, expectedIDs, actualIDs)
			}
		},
	)
}

func TestGetCurrentWeek(t *testing.T) {
	week := frontend.GetCurrentWeek()

	require.Equal(t, time.UTC, week.Location(), "expects server time")

	t.Run(
		"expects week starts on start of day",
		func(t *testing.T) {
			assert.Equal(t, 0, week.Hour())
			assert.Equal(t, 0, week.Minute())
			assert.Equal(t, 0, week.Second())
		},
	)

	t.Run(
		"expects at most 7 days prior",
		func(t *testing.T) {
			assert.Less(t, time.Since(week), 7*24*time.Hour)
		},
	)
}

func TestGetDay(t *testing.T) {
	t.Run(
		"expects week starts on start of day",
		func(t *testing.T) {
			day := frontend.GetDay(0)

			require.Equal(t, time.UTC, day.Location(), "expects server time")

			assert.Equal(t, 0, day.Hour())
			assert.Equal(t, 0, day.Minute())
			assert.Equal(t, 0, day.Second())
		},
	)

	t.Run(
		"expects start of yesterday",
		func(t *testing.T) {
			day := frontend.GetDay(1)

			require.Equal(t, time.UTC, day.Location(), "expects server time")

			now := time.Now().In(time.UTC).Add(-24 * time.Hour)

			assert.Equal(t, now.Year(), day.Year())
			assert.Equal(t, now.Month(), day.Month())
			assert.Equal(t, now.Day(), day.Day())
		},
	)
}

var DefaultPGInstance *postgres.PGInstance

func TestMain(m *testing.M) {
	ctx := context.Background()

	flag.Parse()

	if testing.Short() {
		fmt.Println("Skipping integration tests in short mode")
	} else {
		DefaultPGInstance = postgres.NewPGInstance(ctx)
	}

	code := m.Run()

	if DefaultPGInstance != nil {
		DefaultPGInstance.Terminate(ctx)
	}

	os.Exit(code)
}
