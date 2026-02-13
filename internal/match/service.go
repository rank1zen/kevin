package match

import (
	"context"
	"fmt"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

type MatchService struct {
	riot  *riot.Client
	store MatchStore
}

// GetMatchDetailRequest represents the request payload for retrieving a single match's details.
type GetMatchDetailRequest struct {
	Region  string `json:"region"`
	MatchID string `json:"matchId"`
}

// GetMatchDetail retrieves the details of a specific match.
func (s *MatchService) GetMatchDetail(ctx context.Context, req GetMatchDetailRequest) (*MatchDetail, error) {
	if req.Region == "" {
		req.Region = "NA1" // Default to NA1 if not specified
	}

	found, err := s.store.CheckMatch(ctx, req.MatchID)
	if err != nil {
		return nil, fmt.Errorf("failed to check match IDs in store: %w", err)
	}

	if !found {
		// If the match is new, fetch it from Riot API and record it
		riotMatch, err := s.riot.Match.GetMatch(ctx, riot.Region(req.Region), req.MatchID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch match from Riot API: %w", err)
		}

		match := mapMatch(riotMatch)

		err = s.store.RecordMatch(ctx, match)
		if err != nil {
			return nil, fmt.Errorf("failed to record match in store: %w", err)
		}
	}

	storeMatch, err := s.store.GetMatchDetail(ctx, req.MatchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get match detail from store: %w", err)
	}

	return storeMatch, nil
}

func mapMatch(riotMatch *riot.Match) *Match {
	var winner int
	if riotMatch.Info.Teams[0].Win {
		winner = riotMatch.Info.Teams[0].TeamID
	} else {
		winner = riotMatch.Info.Teams[1].TeamID
	}

	riotMatchInfo := riotMatch.Info

	var (
		gameDate     = time.UnixMilli(riotMatchInfo.GameEndTimestamp)
		gameDuration = time.Duration(riotMatchInfo.GameDuration) * time.Second
	)

	match := Match{
		ID:           riotMatch.Metadata.MatchID,
		Date:         gameDate,
		Duration:     gameDuration,
		Version:      riotMatch.Info.GameVersion,
		WinnerID:     winner,
		Participants: [10]Participant{},
	}

	//
	var (
		blueTeamKills  int
		blueTeamGold   int
		blueTeamDamage int
		redTeamKills   int
		redTeamGold    int
		redTeamDamage  int
	)

	for _, riotParticipant := range riotMatch.Info.Participants {
		if riotParticipant.TeamID == 100 {
			blueTeamKills += riotParticipant.Kills
			blueTeamGold += riotParticipant.GoldEarned
			blueTeamDamage += riotParticipant.TotalDamageDealtToChampions
		} else {
			redTeamKills += riotParticipant.Kills
			redTeamGold += riotParticipant.GoldEarned
			redTeamDamage += riotParticipant.TotalDamageDealtToChampions
		}
	}

	for i, riotParticipant := range riotMatch.Info.Participants {
		counterpartIndex := getCounterPartIndex(i)

		counterpart := riotMatch.Info.Participants[counterpartIndex]

		var (
			teamDamage int
			teamGold   int
			teamKills  int
		)

		if riotParticipant.TeamID == 100 {
			teamKills = blueTeamKills
			teamDamage = blueTeamDamage
			teamGold = blueTeamGold
		} else {
			teamKills = redTeamKills
			teamDamage = redTeamDamage
			teamGold = redTeamGold
		}

		participant := Participant{
			PUUID:                riot.PUUID(riotParticipant.PUUID),
			MatchID:              riotMatch.Metadata.MatchID,
			TeamID:               riotParticipant.TeamID,
			ChampionID:           riotParticipant.ChampionID,
			ChampionLevel:        riotParticipant.ChampLevel,
			TeamPosition:         convertRiotTeamPosition(riotParticipant.TeamPosition),
			SummonerIDs:          convertRiotSummonerSpells(riotParticipant),
			Runes:                newRunePage(withRiotPerks(riotParticipant.Perks)),
			Items:                convertRiotItems(riotParticipant),
			Kills:                riotParticipant.Kills,
			Deaths:               riotParticipant.Deaths,
			Assists:              riotParticipant.Assists,
			KillParticipation:    computeKillParticipation(riotParticipant, teamKills),
			CreepScore:           computeCreepScore(riotParticipant),
			CreepScorePerMinute:  computeCreepScorePerMinute(riotParticipant, gameDuration),
			DamageDealt:          riotParticipant.TotalDamageDealtToChampions,
			DamageTaken:          riotParticipant.TotalDamageTaken,
			DamageDeltaEnemy:     riotParticipant.TotalDamageDealtToChampions - counterpart.TotalDamageDealtToChampions,
			DamagePercentageTeam: float32(riotParticipant.TotalDamageDealtToChampions) / float32(teamDamage),
			GoldEarned:           riotParticipant.GoldEarned,
			GoldDeltaEnemy:       riotParticipant.GoldEarned - counterpart.GoldEarned,
			GoldPercentageTeam:   float32(riotParticipant.GoldEarned) / float32(teamGold),
			VisionScore:          riotParticipant.VisionScore,
			PinkWardsBought:      riotParticipant.DetectorWardsPlaced,
		}

		match.Participants[i] = participant
	}

	return &match
}

func computeCreepScore(p *riot.MatchParticipant) int {
	return p.TotalMinionsKilled + p.NeutralMinionsKilled
}

func computeCreepScorePerMinute(p *riot.MatchParticipant, duration time.Duration) float32 {
	cs := p.TotalMinionsKilled + p.NeutralMinionsKilled
	return float32(cs) / float32(duration.Minutes())
}

func computeKillParticipation(p *riot.MatchParticipant, teamKills int) float32 {
	return float32(p.Assists+p.Kills) / float32(teamKills)
}

func convertRiotUnixToTimestamp(ts int64) time.Time {
	return time.UnixMilli(ts)
}

func convertRiotItems(p *riot.MatchParticipant) [7]int {
	return [7]int{p.Item0, p.Item1, p.Item2, p.Item3, p.Item4, p.Item5, p.Item6}
}

func convertRiotSummonerSpells(p *riot.MatchParticipant) [2]int {
	return [2]int{p.Summoner1ID, p.Summoner2ID}
}

func convertRiotLiveSummonerSpells(p riot.LiveMatchParticipant) [2]int {
	return [2]int{p.Spell1ID, p.Spell2ID}
}

func getCounterPartIndex(index int) int {
	if index >= 5 {
		return index % 5
	}
	return index + 5
}

var teamPositions = map[string]internal.TeamPosition{
	"TOP":     0,
	"JUNGLE":  1,
	"MIDDLE":  2,
	"BOTTOM":  3,
	"UTILITY": 4,
}

func convertRiotTeamPosition(s string) internal.TeamPosition {
	pos, ok := teamPositions[s]
	if !ok {
		panic(fmt.Sprintf("team position %s is not valid", s))
	}

	return pos
}

type RunePageOption func(*internal.RunePage) error

func newRunePage(opts ...RunePageOption) (runes internal.RunePage) {
	for _, f := range opts {
		_ = f(&runes)
	}
	return runes
}

func withRiotPerks(perks *riot.MatchPerks) RunePageOption {
	return func(runes *internal.RunePage) error {
		runes.PrimaryTree = perks.Styles[0].Style
		runes.PrimaryKeystone = perks.Styles[0].Selections[0].Perk
		runes.PrimaryA = perks.Styles[0].Selections[1].Perk
		runes.PrimaryB = perks.Styles[0].Selections[2].Perk
		runes.PrimaryC = perks.Styles[0].Selections[3].Perk
		runes.SecondaryTree = perks.Styles[1].Style
		runes.SecondaryA = perks.Styles[1].Selections[0].Perk
		runes.SecondaryB = perks.Styles[1].Selections[1].Perk
		runes.MiniOffense = perks.StatPerks.Offense
		runes.MiniFlex = perks.StatPerks.Flex
		runes.MiniDefense = perks.StatPerks.Defense
		return nil
	}
}
