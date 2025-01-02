package db

import (
	"context"
	"testing"
	"time"

	"github.com/rank1zen/yujin/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChampion(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	db := setupDB(t)

	var a internal.RiotClient

	match, err := a.GetMatch(ctx, "ada")
	require.NoError(t, err)

	err = db.CreateMatch(ctx, match)
	require.NoError(t, err)

	champions, err := db.GetChampionList(ctx, "a")
	assert.NoError(t, err)

	// example
	assert.Len(t, champions.List, 3)
	assert.Equal(t, champions.List[0].Kills, 4.5)
}
