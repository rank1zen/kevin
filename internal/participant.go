package internal

import (
	"fmt"

	"github.com/rank1zen/kevin/internal/riot"
)

// Participant represents a record of a summoner in a ranked match.
type Participant struct {
	PUUID                  string
	MatchID                string
	TeamID                 int
	ChampionID             int
	ChampionLevel          int
	SummonerIDs            [2]int
	Runes                  RunePage
	Items                  [7]int
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

// WithDefaultParticipant instantiates some valid [Participant], usually used
// for testing.
func WithDefaultParticipant() ParticipantOption {
	return func(p *Participant) error {
		p.PUUID = "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg"
		p.MatchID = "NA1_5304757838"
		p.TeamID = 100
		p.ChampionID = 63
		p.ChampionLevel = 12
		p.SummonerIDs = [2]int{4, 14}
		p.Runes = RunePage{}
		p.Items = [7]int{1056,3116,3020,2508,3802,0,3363}
		p.Kills = 2
		p.Deaths = 0
		p.Assists = 8
		p.KillParticipation = 1
		p.CreepScore = 131
		p.CreepScorePerMinute = 1
		p.DamageDealt = 1
		p.DamageTaken = 1
		p.DamageDeltaEnemy = 1
		p.DamagePercentageTeam = 1
		p.GoldEarned = 1
		p.GoldDeltaEnemy = 1
		p.GoldPercentageTeam = 1
		p.VisionScore = 1
		p.PinkWardsBought = 1
		return nil
	}
}

func RiotMatchToParticipant(match riot.Match, puuid string) ParticipantOption {
	var s *riot.MatchParticipant
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
		p.PUUID = s.PUUID
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
		return nil
	}
}

type LiveParticipant struct {
	PUUID string
	MatchID string
	ChampionID int
	Runes RunePage
	TeamID int
	SummonersIDs [2]int
}

func NewLiveParticipant(opts ...LiveParticipantOption) LiveParticipant {
	var m LiveParticipant
	for _, f := range opts {
		f(&m)
	}
	return m
}

type LiveParticipantOption func(*LiveParticipant) error

func WithRiotCurrentGame(r riot.LiveMatch, puuid string) LiveParticipantOption {
	var selected *riot.LiveMatchParticipant
	for _, p := range r.Participants {
		if p.PUUID == puuid {
			selected = &p
		}
	}
	if selected == nil {
		panic("bro.")
	}

	matchID := fmt.Sprintf("%s_%d", r.PlatformID, r.GameID)

	return func(m *LiveParticipant) error {
		m.PUUID = selected.PUUID
		m.MatchID = matchID
		m.ChampionID = selected.ChampionID
		m.Runes = NewRunePage(WithRiotSpectatorPerks(&selected.Perks))
		m.TeamID = selected.TeamID
		m.SummonersIDs = [2]int{selected.Spell1ID,selected.Spell2ID}
		return nil
	}
}
