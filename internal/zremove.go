package internal

import (
	"errors"
	"fmt"

	"github.com/rank1zen/kevin/internal/riot"
)

// ALL deprecated.

func NewMatch(opts ...MatchOption) Match {
	var match Match
	for _, f := range opts {
		if err := f(&match); err != nil {
			panic(err)
		}
	}

	return match
}

type MatchOption func(*Match) error

func WithRiotMatch(match *riot.Match) MatchOption {
	return func(m *Match) error {
		m.ID = match.Metadata.MatchID
		m.Date = convertRiotUnixToTimestamp(match.Info.GameEndTimestamp)
		m.Duration = convertRiotTimeToDuration(match.Info.GameDuration)
		m.Version = match.Info.GameVersion

		var winner int
		if match.Info.Teams[0].Win {
			winner = match.Info.Teams[0].TeamID
		} else {
			winner = match.Info.Teams[1].TeamID
		}

		m.WinnerID = winner

		for i, p := range match.Info.Participants {
			m.Participants[i] = NewParticipant(WithRiotMatchParticipant(*match, riot.PUUID(p.PUUID)))
		}

		return nil
	}
}

func NewRankDetail(option RankDetailOption) RankDetail {
	m := RankDetail{}
	option(&m)
	return m
}

type RankDetailOption func(*RankDetail)

func WithRiotLeagueEntry(rank riot.LeagueEntry) RankDetailOption {
	return func(m *RankDetail) {
		m.Wins = rank.Wins
		m.Losses = rank.Losses

		m.Rank = Rank{
			Tier:     rank.Tier,
			Division: rank.Division,
			LP:       rank.LeaguePoints,
		}
	}
}

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

func NewLiveMatch(opts ...LiveMatchOption) LiveMatch {
	var match LiveMatch

	for _, f := range opts {
		f(&match)
	}

	return match
}

type LiveMatchOption func(*LiveMatch) error

func WithRiotLiveMatch(match riot.LiveMatch) LiveMatchOption {
	return func(m *LiveMatch) error {
		matchID := fmt.Sprintf("%s_%d", match.PlatformID, match.GameID)

		participants := [10]LiveParticipant{}
		for i, p := range match.Participants {
			participants[i] = NewLiveParticipant(WithRiotLiveMatchParticipant(match, NewPUUIDFromString(p.PUUID)))
		}

		m.ID = matchID
		m.Date = convertRiotUnixToTimestamp(match.GameStartTime)
		m.Participants = [10]LiveParticipant(participants)
		return nil
	}
}

// WithRiotLiveMatchParticipant uses creates participant with puuid using
// [riot.LiveMatch].
func WithRiotLiveMatchParticipant(match riot.LiveMatch, puuid riot.PUUID) LiveParticipantOption {
	var selected *riot.LiveMatchParticipant
	for _, p := range match.Participants {
		if p.PUUID == puuid.String() {
			selected = &p
		}
	}

	return func(m *LiveParticipant) error {
		matchID := fmt.Sprintf("%s_%d", match.PlatformID, match.GameID)

		if selected == nil {
			return errors.New(fmt.Sprintf("puuid %s is not when creating live match", puuid.String()))
		}

		m.PUUID = NewPUUIDFromString(selected.PUUID)
		m.MatchID = matchID
		m.ChampionID = selected.ChampionID
		m.Runes = NewRunePage(WithRiotSpectatorPerks(&selected.Perks))
		m.TeamID = selected.TeamID
		m.SummonersIDs = [2]int{selected.Spell1ID, selected.Spell2ID}
		return nil
	}
}

func NewSummonerMatch(option SummonerMatchOption) SummonerMatch {
	detail := SummonerMatch{}
	option(&detail)
	return detail
}

type SummonerMatchOption func(*SummonerMatch)

// NOTE: should be unexported.
func NewLiveParticipant(opts ...LiveParticipantOption) LiveParticipant {
	var m LiveParticipant

	for _, f := range opts {
		f(&m)
	}

	return m
}

type LiveParticipantOption func(*LiveParticipant) error

func NewParticipantDetail(option ParticipantDetailOption) ParticipantDetail {
	detail := ParticipantDetail{}
	option(&detail)
	return detail
}

type ParticipantDetailOption func(*ParticipantDetail)
