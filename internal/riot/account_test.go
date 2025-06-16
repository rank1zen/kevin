package riot_test

import (
	"context"
	"os"
	"testing"

	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
)

func TestAccountByPuuid(t *testing.T) {
	client := riot.NewClient(riot.WithApiKey(os.Getenv("KEVIN_RIOT_API_KEY")))
	_, err := client.GetAccountByPuuid(context.Background(), riot.RegionAmericas, "xpzpxnzLQX12ACv3iHZfqgdA8RGZQBLCiqJVa1rfVO8Z3KRiYD7YikD2RZC5mot0YhJNKn1UDxu-Ng")
	assert.NoError(t, err)
}
