package riot_test

// import (
// 	"context"
// 	"os"
// 	"testing"
//
// 	"github.com/rank1zen/kevin/internal/riot"
// 	"github.com/stretchr/testify/assert"
// )
//
// func TestMatch(t *testing.T) {
// 	client := riot.NewClient(riot.WithApiKey(os.Getenv("KEVIN_RIOT_API_KEY")))
// 	_, err := client.GetMatchIDsByPUUID(context.Background(), riot.RegionAmericas, "xpzpxnzLQX12ACv3iHZfqgdA8RGZQBLCiqJVa1rfVO8Z3KRiYD7YikD2RZC5mot0YhJNKn1UDxu-Ng", "queue=420&start=0&count=10")
// 	assert.NoError(t, err)
// }
//
// func TestTimeline(t *testing.T) {
// 	client := riot.NewClient(riot.WithApiKey(os.Getenv("KEVIN_RIOT_API_KEY")))
//
// 	t.Run(
// 		"fetch riot",
// 		func(t *testing.T) {
// 			_, err := client.GetMatchTimeline(context.Background(), riot.RegionAmericas, "NA1_5194546103")
// 			assert.NoError(t, err)
// 		},
// 	)
// }
