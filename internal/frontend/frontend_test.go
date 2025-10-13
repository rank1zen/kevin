package frontend_test

import (
	"context"
	"errors"
	"flag"
	"log"
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

func TestMain(t *testing.M) {
	ctx := context.Background()

	flag.Parse()

	if !testing.Short() {
		DefaultPGInstance = postgres.NewPGInstance(context.Background(), "../../migrations/")
	}

	code := t.Run()

	if !testing.Short() {
		if err := DefaultPGInstance.Terminate(ctx); err != nil {
			log.Fatalf("terminating: %s", err)
		}
	}

	os.Exit(code)
}

func SetupDatasource(ctx context.Context, t testing.TB) *internal.Datasource {
	pool := DefaultPGInstance.SetupConn(ctx, t)

	client := riot.NewClient(os.Getenv("KEVIN_RIOT_API_KEY"))

	store := postgres.NewStore(pool)

	return internal.NewDatasource(client, store)
}

func TestGetDay(t *testing.T) {
	timezone, err := time.LoadLocation("America/Toronto")
	require.NoError(t, err)

	ts := time.Date(2025, time.April, 4, 0, 0, 0, 0, timezone)

	days := frontend.GetDays(ts)

	expected := []time.Time{
		time.Date(2025, time.April, 5, 0, 0, 0, 0, timezone),
		time.Date(2025, time.April, 4, 0, 0, 0, 0, timezone),
		time.Date(2025, time.April, 3, 0, 0, 0, 0, timezone),
		time.Date(2025, time.April, 2, 0, 0, 0, 0, timezone),
		time.Date(2025, time.April, 1, 0, 0, 0, 0, timezone),
		time.Date(2025, time.March, 31, 0, 0, 0, 0, timezone),
		time.Date(2025, time.March, 30, 0, 0, 0, 0, timezone),
		time.Date(2025, time.March, 29, 0, 0, 0, 0, timezone),
	}

	require.Len(t, days, len(expected))
	assert.Equal(t, expected, days)
}

func TestParseRiotID(t *testing.T) {
	for _, test := range []struct {
		TestName string
		RiotID   string
		Err      error
		Name     string
		Tag      string
	}{
		{
			TestName: "expects orrange#NA1",
			RiotID:   "orrange-NA1",
			Name:     "orrange",
			Tag:      "NA1",
		},
		{
			TestName: "expects invalid error for missing tag",
			RiotID:   "orrange-",
			Err:      errors.New("invalid riot id"),
		},
		{
			TestName: "expects invalid error for missing tag",
			RiotID:   "orrange",
			Err:      errors.New("invalid riot id"),
		},
		{
			TestName: "expects invalid  error for two seperators",
			RiotID:   "orrange-NA1-NA1",
			Err:      errors.New("invalid riot id"),
		},
	} {
		t.Run(
			test.TestName,
			func(t *testing.T) {
				name, tag, err := frontend.ParseRiotID(test.RiotID)
				if test.Err != nil {
					assert.ErrorIs(t, err, test.Err)
				} else {
					assert.Equal(t, test.Name, name)
					assert.Equal(t, test.Tag, tag)
				}
			},
		)
	}
}
