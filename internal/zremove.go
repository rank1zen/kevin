package internal

import (
	"github.com/rank1zen/kevin/internal/riot"
)

// ALL deprecated.

func NewRankStatus(option RankStatusOption) RankStatus {
	m := RankStatus{}
	option(&m)
	return m
}

type RankStatusOption func(*RankStatus)

type ParticipantOption func(*Participant) error

// NOTE: should be unexported.
func NewParticipant(opts ...ParticipantOption) Participant {
	var p Participant
	for _, f := range opts {
		f(&p)
	}

	return p
}

func WithRiotMatchParticipant(match riot.Match, puuid riot.PUUID) ParticipantOption {
	var s *riot.MatchParticipant
	for _, p := range match.Info.Participants {
		if p.PUUID == string(puuid) {
			s = p
		}
	}

	if s == nil {
		panic("yo puuid is not in this match")
	}

	var teamDamage, teamGold, teamKills int
	for _, p := range match.Info.Participants {
		if p.TeamID == s.TeamID {
			teamDamage += p.TotalDamageDealtToChampions
			teamGold += p.GoldEarned
			teamKills += p.Kills
		}
	}

	var counterpart *riot.MatchParticipant
	for _, p := range match.Info.Participants {
		if p.TeamPosition == s.TeamPosition && p.PUUID != s.PUUID {
			counterpart = p
		}
	}

	var (
		cs       = s.TotalMinionsKilled + s.NeutralMinionsKilled
		csPerMin = float32(cs*60) / float32(match.Info.GameDuration)
		kp       = float32(s.Assists+s.Kills) / float32(teamKills)

		damageDelta = s.TotalDamageDealtToChampions - counterpart.TotalDamageDealtToChampions
		goldDelta   = s.GoldEarned - counterpart.GoldEarned
		damageShare = float32(s.TotalDamageDealtToChampions) / float32(teamDamage)
		goldShare   = float32(s.GoldEarned) / float32(teamGold)

		runes     = NewRunePage(WithRiotPerks(s.Perks))
		items     = [7]int{s.Item0, s.Item1, s.Item2, s.Item3, s.Item4, s.Item5, s.Item6}
		summoners = [2]int{s.Summoner1ID, s.Summoner2ID}
	)

	return func(p *Participant) error {
		p.PUUID = NewPUUIDFromString(s.PUUID)
		p.MatchID = match.Metadata.MatchID
		p.TeamID = s.TeamID
		p.ChampionID = s.ChampionID
		p.ChampionLevel = s.ChampLevel
		p.SummonerIDs = summoners
		p.Runes = runes
		p.Items = items
		p.Kills = s.Kills
		p.Deaths = s.Deaths
		p.Assists = s.Assists
		p.KillParticipation = kp
		p.CreepScore = cs
		p.CreepScorePerMinute = csPerMin
		p.DamageDealt = s.TotalDamageDealtToChampions
		p.DamageTaken = s.TotalDamageTaken
		p.DamageDeltaEnemy = damageDelta
		p.DamagePercentageTeam = damageShare
		p.GoldEarned = s.GoldEarned
		p.GoldDeltaEnemy = goldDelta
		p.GoldPercentageTeam = goldShare
		p.VisionScore = s.VisionScore
		p.PinkWardsBought = s.DetectorWardsPlaced
		p.TeamPosition = convertRiotTeamPosition(s.TeamPosition)
		return nil
	}
}

func NewSummonerMatch(option SummonerMatchOption) SummonerMatch {
	detail := SummonerMatch{}
	option(&detail)
	return detail
}

type SummonerMatchOption func(*SummonerMatch)

func NewParticipantDetail(option ParticipantDetailOption) ParticipantDetail {
	detail := ParticipantDetail{}
	option(&detail)
	return detail
}

type ParticipantDetailOption func(*ParticipantDetail)
