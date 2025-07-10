package riot_test

import (
	"context"
	"os"
	"testing"

	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
)

func TestGetMatchList(t *testing.T) {
	ctx := context.Background()
	client := riot.NewClient(os.Getenv("KEVIN_RIOT_API_KEY"))

	t.Run(
		"sanity test",
		func(t *testing.T) {
			_, err := client.Match.GetMatchList(ctx, riot.ContinentAmericas, "xpzpxnzLQX12ACv3iHZfqgdA8RGZQBLCiqJVa1rfVO8Z3KRiYD7YikD2RZC5mot0YhJNKn1UDxu-Ng", riot.MatchListOptions{Count: 20})
			assert.NoError(t, err)
		},
	)

	t.Run(
		"all matches in a day",
		func(t *testing.T) {
			options := riot.MatchListOptions{
				StartTime: new(int64),
				EndTime:   new(int64),
				Queue:     new(int),
				Start:     0,
				Count:     100,
			}

			*options.StartTime = 1751601600
			*options.EndTime = 1751687999
			*options.Queue = 420

			matches, err := client.Match.GetMatchList(ctx, riot.ContinentAmericas, "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg", options)
			assert.NoError(t, err)

			expected := riot.MatchList{
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

			assert.Equal(t, expected, matches)
		},
	)
}
