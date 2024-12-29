package riot

import (
	"context"

	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/riotclient"
)

func (r *Riot) GetProfile(ctx context.Context, puuid internal.PUUID) (internal.Profile, error) {
	account, err := r.client.AccountGetByPuuid(ctx, puuid.String())
	if err != nil {
		return internal.Profile{}, err
	}

	summoner, err := r.client.GetSummonerByPuuid(ctx, puuid.String())
	if err != nil {
		return internal.Profile{}, err
	}

	leagueEntries, err := r.client.GetLeagueEntriesForSummoner(ctx, puuid.String())
	if err != nil {
		return internal.Profile{}, err
	}

	var soloqEntry *riotclient.LeagueEntry = nil
	for _, entry := range leagueEntries {
		if entry.QueueType == riotclient.QueueTypeRankedSolo5x5 {
			soloqEntry = entry
			break
		}
	}

	var rank *internal.RankRecord = nil
	if soloqEntry != nil {
		rank = &internal.RankRecord{
			Puuid:    internal.PUUID(summoner.Puuid),
			LeagueID: soloqEntry.LeagueId,
			Wins:     soloqEntry.Wins,
			Losses:   soloqEntry.Losses,
			Tier:     soloqEntry.Tier,
			Division: soloqEntry.Rank,
			LP:       soloqEntry.LeaguePoints,
		}
	}

	return internal.Profile{
		Name:          account.GameName,
		Tagline:       account.TagLine,
		Puuid:         internal.PUUID(summoner.Puuid),
		SummonerID:    internal.SummonerID(summoner.Id),
		AccountID:     internal.AccountID(summoner.AccountId),
		Level:         int(summoner.SummonerLevel),
		ProfileIconID: internal.ProfileIconID(summoner.ProfileIconId),
		Rank:          rank,
	}, nil
}
