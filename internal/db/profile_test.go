package db

import (
	"context"
	"log"
	"testing"

	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/riotclient"
	"github.com/rank1zen/yujin/internal/riotclient/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type profileDB interface {
	ProfileExists(ctx context.Context, puuid riotclient.PUUID) (bool, error)
	ProfileGetChampionStatList(ctx context.Context, puuid riotclient.PUUID, season internal.Season) (ProfileChampionStatList, error)
	ProfileGetHeader(ctx context.Context, puuid riotclient.PUUID) (ProfileHeader, error)
	ProfileGetLiveGame(ctx context.Context, puuid string) (ProfileLiveGame, error)
	ProfileGetMatchList(ctx context.Context, puuid riotclient.PUUID, page int, ensure bool) (ProfileMatchList, error)
	ProfileGetRankHistory(ctx context.Context, puuid riotclient.PUUID) (ProfileRankHistoryList, error)
	ProfileUpdate(ctx context.Context, puuid riotclient.PUUID) error
}

func testFetch(t testing.TB, ctx context.Context, db profileDB, puuid riotclient.PUUID) {
	exists, err := db.ProfileExists(ctx, puuid)
	require.NoError(t, err)
	if exists {
	} else {
	}
}

func testMatchList(t testing.TB, ctx context.Context, db profileDB) {
	// assume the app has been running for months
	// we have a bunch of matches for puuid and they are all valid
	// lets fetch their match list and verify against what?
	// I guess we just check for error

	var match riotclient.Match
	// ->
	var out ProfileMatch
	assert.Equal(t, match.Info.GameCreation, out.Date)
}

func TestFindProfile(t *testing.T) {
	ctx := context.Background()

	db := setupDB(t)

	exists, err := db.ProfileExists(ctx, "xpzpxnzLQX12ACv3iHZfqgdA8RGZQBLCiqJVa1rfVO8Z3KRiYD7YikD2RZC5mot0YhJNKn1UDxu-Ng")
	require.NoError(t, err)
	require.True(t, exists)

	m, err := db.ProfileGetHeader(ctx, "xpzpxnzLQX12ACv3iHZfqgdA8RGZQBLCiqJVa1rfVO8Z3KRiYD7YikD2RZC5mot0YhJNKn1UDxu-Ng")
	assert.NoError(t, err)

	log.Print(m)
}

func TestFindDneProfile(t *testing.T) {
	t.Parallel()
}

func TestProfileGetMatchListTestdata(t *testing.T) {
	ctx := context.Background()

	db := setupDB(t)

	err := matchInsert(ctx, db.pool, testdata.GetMatch("NA1_5011055088"))
	require.NoError(t, err)

	m, err := db.ProfileGetMatchList(ctx, "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q", 0, false)
	require.NoError(t, err)
	if assert.Len(t, m.List, 1) {
		log.Println(m.List[0])
	}
}

func TestProfileMatch(t *testing.T) {
	ctx := context.Background()

	db := setupDB(t)

	var match riotclient.Match

	db.ProfileGetMatchList(ctx, "", 0, false)
	var res ProfileMatch

	assert.Equal(t, match.Metadata.MatchId, res.MatchID)
	assert.Equal(t, match.Info.GameEndTimestamp, res.Patch)
	assert.Equal(t, match.Info.GameVersion, res.Date)
	assert.Equal(t, match.Info.GameDuration, res.Duration)
}

func TestProfileGetChampionStatList(t *testing.T) {
	ctx := context.Background()

	db := setupDB(t)

	err := matchInsert(ctx, db.pool, testdata.GetMatch("NA1_5011055088"))
	require.NoError(t, err)

	stats, err := db.ProfileGetChampionStatList(ctx, "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q", internal.SeasonAll)
	require.NoError(t, err)
	log.Print(stats)
}
