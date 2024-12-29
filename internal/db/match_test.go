package db

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/yujin/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// func TestEnsureMatchlist(t *testing.T) {
// 	t.Parallel()
//
// 	ctx := context.Background()
//
// 	db := setupDB(t)
//
// 	err := ensureMatchList(ctx, db.pool, db.riot, "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q", 0, 1)
// 	assert.NoError(t, err)
// }

func TestCreateMatch(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	db := setupDB(t)

	for _, test := range []struct {
		desc     string
		expected internal.Match
	}{
		{
			"Real data",
			internal.Match{
				ID:              "NA1_5011055088",
				DataVersion:     "2",
				Patch:           internal.GameVersion("14.11.589.9418"),
				EndOfGameResult: "GameComplete",
				CreateTimestamp: time.Unix(1717303471, 0),
				EndTimestamp:    time.Unix(1717304694, 0),
				StartTimestamp:  time.Unix(1717303529, 0),
				Duration:        1165 * time.Second,
			},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			actual, err := createMatch(ctx, db.pool, test.expected)

			require.NoError(t, err)

			assert.Equal(t, test.expected, actual)
		})
	}
}
