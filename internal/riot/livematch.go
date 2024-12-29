package riot

import (
	"context"
	"time"

	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/riotclient"
)

// NOTE: don't know what this is.
func getI(match riotclient.SpectatorCurrentGameInfo, position int) riotclient.SpectatorCurrentGameParticipant {
	return match.Participants[position]
}

// getBan finds the champion player at position banned.
// NOTE: I DONT KNOW IF THIS WORKS
func getBan(match riotclient.SpectatorCurrentGameInfo, position int) *internal.ChampionID {
	return (*internal.ChampionID)(&match.BannedChampions[position].ChampionId)
}

func (r *Riot) GetLiveMatch(ctx context.Context, puuid internal.PUUID) (internal.LiveMatch, error) {
	gameInfo, err := r.client.GetCurrentGameInfoByPuuid(ctx, puuid.String())
	if err != nil {
		return internal.LiveMatch{}, err
	}

	// probably check that there is indeed some number of things
	participants := internal.LiveMatchParticipantList{}

	for i := range 10 {
		p := getI(gameInfo, i)
		banned := getBan(gameInfo, i)

		participants[i] = internal.LiveMatchParticipant{
			Champion:   internal.ChampionID(p.ChampionId),
			Puuid:      internal.PUUID(p.Puuid),
			TeamID:     internal.TeamID(p.TeamId),
			SummonerID: internal.SummonerID(p.SummonerId),
			Summoners: internal.SummsIDs{
				internal.SummsID(p.Spell1Id),
				internal.SummsID(p.Spell2Id),
			},
			BannedChampion: banned,
		}
	}

	return internal.LiveMatch{
		StartTimestamp: riotUnixToDate(gameInfo.GameStartTime),
		Length:         time.Duration(gameInfo.GameLength),
		Participant:    participants,
	}, nil
}
