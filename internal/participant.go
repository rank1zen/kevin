package internal

import (
	"github.com/rank1zen/kevin/internal/riot"
)

type Participant struct {
	Puuid                  string
	MatchID                string
	Team                   Team
	Champion               Champion
	ChampionLevel          int
	Summoners              [2]Spell
	Runes                  RunePage
	Items                  [7]Item
	Kills, Deaths, Assists int
	KillParticipation      float32
	CreepScore             int
	CreepScorePerMinute    float32
	DamageDealt            int
	DamageTaken            int
	DamageDeltaEnemy       int
	DamagePercentageTeam   float32
	GoldEarned             int
	GoldDeltaEnemy         int
	GoldPercentageTeam     float32
	VisionScore            int
	PinkWardsBought        int
}

type ParticipantOption func(*Participant) error

func NewParticipant(opts ...ParticipantOption) Participant {
	var p Participant
	for _, f := range opts {
		f(&p)
	}
	return p
}

func RiotMatchToParticipant(match riot.Match, puuid string) ParticipantOption {
	var s *riot.Participant
	for _, p := range match.Info.Participants {
		if p.PUUID == puuid {
			s = p
		}
	}

	if s == nil {
		panic("yo puuid is not in this match")
	}

	var teamDamage, teamGold, teamKills int
	for _, p := range match.Info.Participants {
		if p.TeamId == s.TeamId {
			teamDamage += p.TotalDamageDealtToChampions
			teamGold += p.GoldEarned
			teamKills += p.Kills
		}
	}

	var counterpart *riot.Participant
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

		runes     = NewRunePage(RiotPerksToRunePage(s.Perks))
		items     = makeItems(s.Item0, s.Item1, s.Item2, s.Item3, s.Item4, s.Item5, s.Item6)
		summoners = [2]Spell{Spell(s.Summoner1Id), Spell(s.Summoner2Id)}
	)

	return func(p *Participant) error {
		p.Puuid = s.PUUID
		p.MatchID = match.Metadata.MatchId
		p.Team = Team(s.TeamId)
		p.Champion = Champion(s.ChampionId)
		p.ChampionLevel = s.ChampLevel
		p.Summoners = summoners
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
		return nil
	}
}
