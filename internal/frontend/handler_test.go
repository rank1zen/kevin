package frontend_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/postgres"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	os.Exit(code)
}

func TestGetSummonerPage(t *testing.T) {
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

			if assert.IsType(t, frontend.SummonerPage{}, component) {
				page, ok := component.(frontend.SummonerPage)
				require.True(t, ok)

				assert.Equal(t, page.Region, riot.RegionNA1)
				assert.Equal(t, page.Name, "orrange")
				assert.Equal(t, page.Tag, "NA1")
				assert.Equal(t, page.PUUID, "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q")
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

func TestGetSummonerMatchHistory(t *testing.T) {
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
			localTimeUnix := localTime.Unix()

			component, err := handler.GetSummonerMatchHistory(ctx, riot.RegionNA1, "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg", localTimeUnix)
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
