package internal

import (
	"fmt"
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

type RiotToProfileMapper struct {
	Account riot.Account

	Rank *riot.LeagueEntry

	EffectiveDate time.Time
}

func (m RiotToProfileMapper) Convert() Profile {
	profile := Profile{
		PUUID:   m.Account.PUUID,
		Name:    m.Account.GameName,
		Tagline: m.Account.TagLine,
		Rank: RankStatus{
			PUUID:         m.Account.PUUID,
			EffectiveDate: m.EffectiveDate,
			Detail:        nil,
		},
	}

	if m.Rank != nil {
		profile.Rank.Detail = &RankDetail{
			Wins:   m.Rank.Wins,
			Losses: m.Rank.Losses,
			Rank: Rank{
				Tier:     m.Rank.Tier,
				Division: m.Rank.Division,
				LP:       m.Rank.LeaguePoints,
			},
		}
	}

	return profile
}

type RiotToLiveMatchMapper struct {
	Match riot.LiveMatch
}

func (m RiotToLiveMatchMapper) Map() LiveMatch {
	riotMatch := m.Match

	livematch := LiveMatch{
		ID:           fmt.Sprintf("%s_%d", riotMatch.PlatformID, riotMatch.GameID),
		Date:         convertRiotUnixToTimestamp(riotMatch.GameStartTime),
		Participants: [10]LiveParticipant{},
	}

	for i := range livematch.Participants {
		livematch.Participants[i] = RiotToLiveMatchParticipantMapper{
			Participant: riotMatch.Participants[i],
			MatchID:     fmt.Sprintf("%s_%d", riotMatch.PlatformID, riotMatch.GameID),
		}.Map()
	}

	return livematch
}

type RiotToLiveMatchParticipantMapper struct {
	Participant riot.LiveMatchParticipant

	MatchID string
}

func (m RiotToLiveMatchParticipantMapper) Map() LiveParticipant {
	riotParticipant := m.Participant

	liveparticipant := LiveParticipant{
		PUUID:        riot.PUUID(riotParticipant.PUUID),
		MatchID:      m.MatchID,
		ChampionID:   m.Participant.ChampionID,
		Runes:        NewRunePage(WithRiotSpectatorPerks(&riotParticipant.Perks)),
		SummonersIDs: convertRiotLiveSummonerSpells(riotParticipant),
		TeamID:       riotParticipant.TeamID,
	}

	return liveparticipant
}

func convertRiotUnixToTimestamp(ts int64) time.Time {
	return time.UnixMilli(ts)
}

func convertRiotLiveSummonerSpells(p riot.LiveMatchParticipant) [2]int {
	return [2]int{p.Spell1ID, p.Spell2ID}
}
