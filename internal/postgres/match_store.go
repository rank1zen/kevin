package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

// TODO: should rename this to something else
type ZMatchStore Store

func NewZMatchStore(pool *pgxpool.Pool) *ZMatchStore {
	return &ZMatchStore{Pool: pool}
}

func (db *ZMatchStore) RecordMatch(ctx context.Context, match internal.Match) error {
	matchStore := MatchStore{Tx: db.Pool}

	batch := pgx.Batch{}

	matchStore.CreateMatchInBatch(&batch, Match{
		ID:       match.ID,
		Date:     match.Date,
		Duration: match.Duration,
		Version:  match.Version,
		WinnerID: match.WinnerID,
	})

	for _, p := range match.Participants {
		matchStore.CreateParticipantInBatch(&batch, Participant{
			PUUID:                p.PUUID.String(),
			MatchID:              p.MatchID,
			TeamID:               p.TeamID,
			ChampionID:           p.ChampionID,
			ChampionLevel:        p.ChampionLevel,
			TeamPosition:         convertTeamPositionToString(p.TeamPosition),
			SummonerIDs:          p.SummonerIDs,
			Runes:                convertRunePageToList(p.Runes),
			Items:                p.Items,
			Kills:                p.Kills,
			Deaths:               p.Deaths,
			Assists:              p.Assists,
			KillParticipation:    p.KillParticipation,
			CreepScore:           p.CreepScore,
			CreepScorePerMinute:  p.CreepScorePerMinute,
			DamageDealt:          p.DamageDealt,
			DamageTaken:          p.DamageTaken,
			DamageDeltaEnemy:     p.DamageDeltaEnemy,
			DamagePercentageTeam: p.DamagePercentageTeam,
			GoldEarned:           p.GoldEarned,
			GoldDeltaEnemy:       p.GoldDeltaEnemy,
			GoldPercentageTeam:   p.GoldPercentageTeam,
			VisionScore:          p.VisionScore,
			PinkWardsBought:      p.PinkWardsBought,
		})
	}

	br := db.Pool.SendBatch(ctx, &batch)

	return br.Close()
}

func (db *ZMatchStore) GetMatchDetail(ctx context.Context, id string) (_ internal.MatchDetail, err error) {
	defer errWrap(&err, "GetMatchDetail(ctx, %s)", id)

	m := internal.MatchDetail{}

	matchStore := MatchStore{Tx: db.Pool}

	summonerStore := SummonerStore{Tx: db.Pool}

	match, err := matchStore.GetMatch(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return m, internal.ErrMatchNotFound
		}

		return m, fmt.Errorf("%w: %w", internal.ErrUnknownStoreError, err)
	}

	participants, err := matchStore.GetParticipants(ctx, id)
	if err != nil {
		return m, fmt.Errorf("%w: %w", internal.ErrUnknownStoreError, err)
	}

	if len(participants) != 10 {
		return m, internal.ErrMatchMissingParticipants
	}

	detail := internal.MatchDetail{
		ID:           id,
		Date:         match.Date,
		Duration:     match.Duration,
		Version:      match.Version,
		WinnerID:     match.WinnerID,
		Participants: [10]internal.ParticipantDetail{},
	}

	for i := range 10 {
		puuid := participants[i].PUUID
		summoner, err := summonerStore.GetSummoner(ctx, riot.PUUID(puuid))
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return m, err
			}

			detail.Participants[i] = toParticipantDetail(participants[i], Summoner{PUUID: riot.PUUID(puuid)}, nil, nil, nil)
			continue
		}

		detail.Participants[i] = toParticipantDetail(participants[i], summoner, nil, nil, nil)
	}

	return detail, nil
}

func (db *ZMatchStore) GetMatchlist(ctx context.Context, puuid riot.PUUID, start, end time.Time) (matches []internal.SummonerMatch, err error) {
	defer errWrap(&err, "GetMatchHistory")

	matchStore := MatchStore{Tx: db.Pool}

	rankStore := RankStore{Tx: db.Pool}

	ids, err := matchStore.ListMatchHistoryIDs(ctx, puuid, start, end)
	if err != nil {
		return nil, err
	}

	matchHistory := []internal.SummonerMatch{}

	var (
		matchList       = []Match{}
		participantList = []Participant{}
		rankBeforeList  = []*RankFull{}
		rankAfterList   = []*RankFull{}
	)

	for _, id := range ids {
		match, err := matchStore.GetMatch(ctx, id)
		if err != nil {
			return nil, err
		}

		matchList = append(matchList, match)

		participant, err := matchStore.GetParticipant(ctx, puuid, id)
		if err != nil {
			return nil, err
		}

		participantList = append(participantList, participant)
	}

	listRankOptions := CreateListRankOption(matchList)
	for _, opt := range listRankOptions {
		statusIDs, err := rankStore.ListRankIDs(ctx, puuid, opt[0])
		if err != nil {
			return nil, err
		}

		id := chooseStatusID(statusIDs)
		if id != nil {
			status, detail, err := (*Store)(db).getRank(ctx, *id)
			if err != nil {
				return nil, err
			}

			rankBeforeList = append(rankBeforeList, &RankFull{Status: status, Detail: detail})
		} else {
			rankBeforeList = append(rankBeforeList, nil)
		}

		statusIDs, err = rankStore.ListRankIDs(ctx, puuid, opt[1])
		if err != nil {
			return nil, err
		}

		id = chooseStatusID(statusIDs)
		if id != nil {
			status, detail, err := (*Store)(db).getRank(ctx, *id)
			if err != nil {
				return nil, err
			}

			rankAfterList = append(rankAfterList, &RankFull{Status: status, Detail: detail})
		} else {
			rankAfterList = append(rankAfterList, nil)
		}
	}

	for i := range ids {
		match := toSummonerMatch(matchList[i], participantList[i], rankBeforeList[i], rankAfterList[i])
		matchHistory = append(matchHistory, match)
	}

	return matchHistory, nil
}

func (db *ZMatchStore) GetNewMatchIDs(ctx context.Context, ids []string) (newIDs []string, err error) {
	matchStore := MatchStore{Tx: db.Pool}
	return matchStore.GetNewMatchIDs(ctx, ids)
}

func (db *ZMatchStore) GetChampions(ctx context.Context, puuid riot.PUUID, start, end time.Time) ([]internal.SummonerChampion, error) {
	matchStore := MatchStore{Tx: db.Pool}

	return matchStore.GetSummonerChampions(ctx, puuid, start, end)
}

func toSummonerMatch(match Match, participant Participant, rankBefore, rankAfter *RankFull) internal.SummonerMatch {
	result := internal.SummonerMatch{
		Participant: internal.Participant{
			PUUID:                riot.PUUID(participant.PUUID),
			MatchID:              match.ID,
			TeamID:               participant.TeamID,
			ChampionID:           participant.ChampionID,
			ChampionLevel:        participant.ChampionLevel,
			TeamPosition:         convertStringToTeamPosition(participant.TeamPosition),
			SummonerIDs:          [2]int{},
			Runes:                internal.RunePage{},
			Items:                [7]int{},
			Kills:                participant.Kills,
			Deaths:               participant.Deaths,
			Assists:              participant.Assists,
			KillParticipation:    participant.KillParticipation,
			CreepScore:           participant.CreepScore,
			CreepScorePerMinute:  participant.CreepScorePerMinute,
			DamageDealt:          participant.DamageDealt,
			DamageTaken:          participant.DamageTaken,
			DamageDeltaEnemy:     participant.DamageDeltaEnemy,
			DamagePercentageTeam: participant.DamagePercentageTeam,
			GoldEarned:           participant.GoldEarned,
			GoldDeltaEnemy:       participant.GoldDeltaEnemy,
			GoldPercentageTeam:   participant.GoldPercentageTeam,
			VisionScore:          participant.VisionScore,
			PinkWardsBought:      participant.PinkWardsBought,
		},
		Date:       match.Date,
		Duration:   match.Duration,
		Win:        false,
		RankBefore: nil,
		RankAfter:  nil,
	}

	if participant.TeamID == match.WinnerID {
		result.Win = true
	}

	if rankBefore != nil {
		rank := toRankStatus(rankBefore)
		result.RankBefore = &rank
	}

	if rankAfter != nil {
		rank := toRankStatus(rankAfter)
		result.RankAfter = &rank
	}

	return result
}

func toParticipantDetail(participant Participant, summoner Summoner, currentRank, rankBefore, rankAfter *RankFull) internal.ParticipantDetail {
	result := internal.ParticipantDetail{
		Participant: internal.Participant{
			PUUID:                riot.PUUID(participant.PUUID),
			MatchID:              participant.MatchID,
			TeamID:               participant.TeamID,
			ChampionID:           participant.ChampionID,
			ChampionLevel:        participant.ChampionLevel,
			TeamPosition:         convertStringToTeamPosition(participant.TeamPosition),
			SummonerIDs:          [2]int{},
			Runes:                internal.RunePage{},
			Items:                [7]int{},
			Kills:                participant.Kills,
			Deaths:               participant.Deaths,
			Assists:              participant.Assists,
			KillParticipation:    participant.KillParticipation,
			CreepScore:           participant.CreepScore,
			CreepScorePerMinute:  participant.CreepScorePerMinute,
			DamageDealt:          participant.DamageDealt,
			DamageTaken:          participant.DamageTaken,
			DamageDeltaEnemy:     participant.DamageDeltaEnemy,
			DamagePercentageTeam: participant.DamagePercentageTeam,
			GoldEarned:           participant.GoldEarned,
			GoldDeltaEnemy:       participant.GoldDeltaEnemy,
			GoldPercentageTeam:   participant.GoldPercentageTeam,
			VisionScore:          participant.VisionScore,
			PinkWardsBought:      participant.PinkWardsBought,
		},
		Name:        summoner.Name,
		Tag:         summoner.Tagline,
		CurrentRank: nil,
		RankBefore:  nil,
		RankAfter:   nil,
	}

	if rankAfter != nil {
		rank := toRankStatus(rankAfter)
		result.RankAfter = &rank
	}

	if rankBefore != nil {
		rank := toRankStatus(rankBefore)
		result.RankBefore = &rank
	}

	if currentRank != nil {
		rank := toRankStatus(currentRank)
		result.CurrentRank = &rank
	}

	return result
}

func toRankStatus(rank *RankFull) internal.RankStatus {
	result := internal.RankStatus{
		PUUID:         riot.PUUID(rank.Status.PUUID),
		EffectiveDate: rank.Status.EffectiveDate,
		Detail:        nil,
	}

	if rank.Detail != nil {
		result.Detail = &internal.RankDetail{
			Wins:   rank.Detail.Wins,
			Losses: rank.Detail.Losses,
			Rank: internal.Rank{
				Tier:     convertStringToRiotTier(rank.Detail.Tier),
				Division: convertStringToRiotRank(rank.Detail.Division),
				LP:       rank.Detail.LP,
			},
		}
	}

	return result
}
