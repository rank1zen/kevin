package internal

import (
	"fmt"
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

type RiotToProfileMapper struct {
	Account riot.Account

	Rank *riot.LeagueEntry
}

func (m RiotToProfileMapper) Convert() Profile {
	profile := Profile{
		PUUID:   "",
		Name:    "",
		Tagline: "",
		Rank:    RankStatus{},
	}

	return profile
}

type RiotToMatchMapper struct {
	Match riot.Match
}

func (m RiotToMatchMapper) Map() Match {
	riotMatch := m.Match

	var winner int
	if m.Match.Info.Teams[0].Win {
		winner = m.Match.Info.Teams[0].TeamID
	} else {
		winner = m.Match.Info.Teams[1].TeamID
	}

	match := Match{
		ID:           m.Match.Metadata.MatchID,
		Date:         convertRiotUnixToTimestamp(m.Match.Info.GameEndTimestamp),
		Duration:     convertRiotTimeToDuration(m.Match.Info.GameDuration),
		Version:      m.Match.Info.GameVersion,
		WinnerID:     winner,
		Participants: [10]Participant{},
	}

	for i, p := range m.Match.Info.Participants {
		match.Participants[i] = RiotToParticipantMapper{
			Participant:       *p,
			MatchID:           riotMatch.Metadata.MatchID,
			MatchDuration:     convertRiotTimeToDuration(riotMatch.Info.GameDuration),
			TeamKills:         1, // TODO: compute these
			TeamGold:          1,
			TeamDamage:        1,
			CounterpartGold:   1,
			CounterpartDamage: 1,
		}.Map()
	}

	return match
}

type RiotToParticipantMapper struct {
	Participant riot.MatchParticipant

	MatchID string
	MatchDuration time.Duration

	TeamKills int
	TeamGold  int
	TeamDamage int

	CounterpartGold int
	CounterpartDamage int
}

func (m RiotToParticipantMapper) Map() Participant {
	riotParticipant := m.Participant

	participant := Participant{
		PUUID:                NewPUUIDFromString(riotParticipant.PUUID),
		MatchID:              m.MatchID,
		TeamID:               riotParticipant.TeamID,
		ChampionID:           riotParticipant.ChampionID,
		ChampionLevel:        riotParticipant.ChampLevel,
		TeamPosition:         convertRiotTeamPosition(riotParticipant.TeamPosition),
		SummonerIDs:          convertRiotSummonerSpells(riotParticipant),
		Runes:                NewRunePage(WithRiotPerks(riotParticipant.Perks)),
		Items:                convertRiotItems(riotParticipant),
		Kills:                riotParticipant.Kills,
		Deaths:               riotParticipant.Deaths,
		Assists:              riotParticipant.Assists,
		KillParticipation:    computeKillParticipation(riotParticipant, m.TeamKills),
		CreepScore:           computeCreepScore(riotParticipant),
		CreepScorePerMinute:  computeCreepScorePerMinute(riotParticipant, m.MatchDuration),
		DamageDealt:          riotParticipant.TotalDamageDealtToChampions,
		DamageTaken:          riotParticipant.TotalDamageTaken,
		DamageDeltaEnemy:     riotParticipant.TotalDamageDealtToChampions - m.CounterpartDamage,
		DamagePercentageTeam: float32(riotParticipant.TotalDamageDealtToChampions) / float32(m.TeamDamage),
		GoldEarned:           riotParticipant.GoldEarned,
		GoldDeltaEnemy:       riotParticipant.GoldEarned - m.CounterpartGold,
		GoldPercentageTeam:   float32(riotParticipant.GoldEarned) / float32(m.TeamGold),
		VisionScore:          riotParticipant.VisionScore,
		PinkWardsBought:      riotParticipant.DetectorWardsPlaced,
	}

	return participant
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
		PUUID:        NewPUUIDFromString(riotParticipant.PUUID),
		MatchID:      m.MatchID,
		ChampionID:   m.Participant.ChampionID,
		Runes:        NewRunePage(WithRiotSpectatorPerks(&riotParticipant.Perks)),
		SummonersIDs: convertRiotLiveSummonerSpells(riotParticipant),
		TeamID:       riotParticipant.TeamID,
	}

	return liveparticipant
}

func computeCreepScore(p riot.MatchParticipant) int {
	return p.TotalMinionsKilled + p.NeutralMinionsKilled
}

func computeCreepScorePerMinute(p riot.MatchParticipant, duration time.Duration) float32 {
	cs := p.TotalMinionsKilled + p.NeutralMinionsKilled
	return float32(cs*60) / float32(duration)
}

func computeKillParticipation(p riot.MatchParticipant, teamKills int) float32 {
	return float32(p.Assists+p.Kills) / float32(teamKills)
}

func convertRiotUnixToTimestamp(ts int64) time.Time {
	return time.UnixMilli(ts)
}

func convertRiotTimeToDuration(t int64) time.Duration {
	return time.Second * time.Duration(t)
}

func convertRiotItems(p riot.MatchParticipant) [7]int {
	return [7]int{p.Item0, p.Item1, p.Item2, p.Item3, p.Item4, p.Item5, p.Item6}
}

func convertRiotSummonerSpells(p riot.MatchParticipant) [2]int {
	return [2]int{p.Summoner1ID, p.Summoner2ID}
}

func convertRiotLiveSummonerSpells(p riot.LiveMatchParticipant) [2]int {
	return [2]int{p.Spell1ID, p.Spell2ID}
}
