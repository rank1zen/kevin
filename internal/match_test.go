package internal_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMatch(t *testing.T) {
	testdata := os.DirFS("../testdata")

	matchFile, err := testdata.Open("NA1_5304757838.json")
	require.NoError(t, err)

	var riotMatch riot.Match
	err = json.NewDecoder(matchFile).Decode(&riotMatch)
	require.NoError(t, err)

	expected := internal.Match{
		ID:       "NA1_5304757838",
		Date:     time.UnixMilli(1749596377340),
		Duration: 1131 * time.Second,
		Version:  "15.11.685.5259",
		WinnerID: 100,
	}

	t.Run(
		"create match from riot",
		func(t *testing.T) {
			actual := internal.NewMatch(internal.WithRiotMatch(&riotMatch))
			assert.Equal(t, expected, actual)
		},
	)
}
