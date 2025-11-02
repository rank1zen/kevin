package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

// Should be unexported
func CreateListRankOption(history []Match) [][2]ListRankOption {
	opt := [][2]ListRankOption{}

	for i, m := range history {
		rankBefore := ListRankOption{
			Offset: 0,
			Limit:  100,
			Recent: true,
		}

		rankAfter := ListRankOption{
			Offset: 0,
			Limit:  100,
			Recent: true,
		}

		if i > 0 {
			rankBefore.Start = &history[i-1].Date
		}

		rankBefore.End = &m.Date

		if i < len(history)-1 {
			rankAfter.End = &history[i+1].Date
		}

		rankAfter.Start = &m.Date

		opt = append(opt, [2]ListRankOption{rankBefore, rankAfter})
	}

	return opt
}

// Store manages connections with a postgres database.
type Store struct {
	Pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *internal.Store {
	return &Store{Pool: pool}
}

func ParticipantDetailFromPG(participant Participant, summoner Summoner, currentRank, rankBefore, rankAfter *RankFull) internal.ParticipantDetailOption {
	return func(m *internal.ParticipantDetail) {
		m.PUUID = internal.NewPUUIDFromString(participant.PUUID)

		m.MatchID = participant.MatchID

		m.TeamID = participant.TeamID

		m.ChampionID = participant.ChampionID

		m.ChampionLevel = participant.ChampionLevel

		m.TeamPosition = convertStringToTeamPosition(participant.TeamPosition)

		m.SummonerIDs = participant.SummonerIDs

		m.Runes = convertListToRunePage(participant.Runes)

		m.Items = participant.Items

		m.Kills = participant.Kills

		m.Deaths = participant.Deaths

		m.Assists = participant.Assists

		m.KillParticipation = participant.KillParticipation

		m.CreepScore = participant.CreepScore

		m.CreepScorePerMinute = participant.CreepScorePerMinute

		m.DamageDealt = participant.DamageDealt

		m.DamageTaken = participant.DamageTaken

		m.DamageDeltaEnemy = participant.DamageDeltaEnemy

		m.DamagePercentageTeam = participant.DamagePercentageTeam

		m.GoldEarned = participant.GoldEarned

		m.GoldDeltaEnemy = participant.GoldDeltaEnemy

		m.GoldPercentageTeam = participant.GoldPercentageTeam

		m.VisionScore = participant.VisionScore

		m.PinkWardsBought = participant.PinkWardsBought

		m.Name = summoner.Name

		m.Tag = summoner.Tagline

		m.RankAfter = nil
		if rankAfter != nil {
			rank := internal.NewRankStatus(WithPostgresRankStatus2(rankAfter))
			m.RankBefore = &rank
		}

		m.RankBefore = nil
		if rankBefore != nil {
			rank := internal.NewRankStatus(WithPostgresRankStatus2(rankBefore))
			m.RankBefore = &rank
		}

		m.CurrentRank = nil
		if currentRank != nil {
			rank := internal.NewRankStatus(WithPostgresRankStatus2(currentRank))
			m.RankBefore = &rank
		}
	}
}

// Deprecated: not using.
func WithPostgresSummonerMatch(match Match, participant Participant, rankBefore, rankAfter *RankFull) internal.SummonerMatchOption {
	return func(m *internal.SummonerMatch) {
		m.PUUID = internal.NewPUUIDFromString(participant.PUUID)

		m.MatchID = match.ID

		m.TeamID = participant.TeamID

		m.ChampionID = participant.ChampionID

		m.ChampionLevel = participant.ChampionLevel

		m.TeamPosition = convertStringToTeamPosition(participant.TeamPosition)

		m.SummonerIDs = participant.SummonerIDs

		m.Runes = convertListToRunePage(participant.Runes)

		m.Items = participant.Items

		m.Kills = participant.Kills

		m.Deaths = participant.Deaths

		m.Assists = participant.Assists

		m.KillParticipation = participant.KillParticipation

		m.CreepScore = participant.CreepScore

		m.CreepScorePerMinute = participant.CreepScorePerMinute

		m.DamageDealt = participant.DamageDealt

		m.DamageTaken = participant.DamageTaken

		m.DamageDeltaEnemy = participant.DamageDeltaEnemy

		m.DamagePercentageTeam = participant.DamagePercentageTeam

		m.GoldEarned = participant.GoldEarned

		m.GoldDeltaEnemy = participant.GoldDeltaEnemy

		m.GoldPercentageTeam = participant.GoldPercentageTeam

		m.VisionScore = participant.VisionScore

		m.PinkWardsBought = participant.PinkWardsBought

		m.Date = match.Date

		m.Duration = match.Duration

		m.Win = false
		if participant.TeamID == match.WinnerID {
			m.Win = true
		}

		m.RankBefore = nil
		if rankBefore != nil {
			rank := internal.NewRankStatus(WithPostgresRankStatus2(rankBefore))
			m.RankBefore = &rank
		}

		m.RankAfter = nil
		if rankAfter != nil {
			rank := internal.NewRankStatus(WithPostgresRankStatus2(rankAfter))
			m.RankAfter = &rank
		}
	}
}

// Deprecated: not using.
func WithPostgresRankStatus2(rank *RankFull) internal.RankStatusOption {
	return func(m *internal.RankStatus) {
		m.PUUID = internal.NewPUUIDFromString(rank.Status.PUUID)
		m.EffectiveDate = rank.Status.EffectiveDate

		m.Detail = nil
		if detail := rank.Detail; detail != nil {
			m.Detail = &internal.RankDetail{
				Wins:   detail.Wins,
				Losses: detail.Losses,
				Rank: internal.Rank{
					Tier:     convertStringToRiotTier(detail.Tier),
					Division: convertStringToRiotRank(detail.Division),
					LP:       detail.LP,
				},
			}
		}
	}
}

type PostgresSearchResult2 struct {
	Summoner Summoner

	Status *RankStatus
	Detail *RankDetail
}

func (m *PostgresSearchResult2) Convert() internal.SearchResult {
	result := internal.SearchResult{
		PUUID:   m.Summoner.PUUID,
		Name:    m.Summoner.Name,
		Tagline: m.Summoner.Tagline,
		Rank:    nil,
	}

	if m.Status != nil {
		result.Rank = &internal.RankStatus{
			PUUID:         m.Summoner.PUUID,
			EffectiveDate: m.Status.EffectiveDate,
			Detail:        nil,
		}

		if m.Detail != nil {
			result.Rank.Detail = &internal.RankDetail{
				Wins:   m.Detail.Wins,
				Losses: m.Detail.Losses,
				Rank: internal.Rank{
					Tier:     riot.Tier(m.Detail.Tier),
					Division: riot.Division(m.Detail.Division),
					LP:       m.Detail.LP,
				},
			}
		}
	}

	return result
}

// chooseStatusID chooses some id that is suitable.
func chooseStatusID(statusIDs []int) *int {
	if len(statusIDs) == 0 {
		return nil
	}
	return &statusIDs[0]
}

// getRank returns both status and detail. It will error if both are not found.
func (db *Store) getRank(ctx context.Context, statusID int) (RankStatus, *RankDetail, error) {
	rankStore := RankStore{Tx: db.Pool}

	status, err := rankStore.GetRankStatus(ctx, statusID)
	if err != nil {
		return RankStatus{}, nil, err
	}

	detail, err := rankStore.GetRankDetail(ctx, statusID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return status, nil, nil
		}

		return RankStatus{}, nil, err
	}

	return status, &detail, err
}

func (db *Store) getMostRecentRank(ctx context.Context, puuid riot.PUUID) (m RankStatus, n *RankDetail, err error) {
	defer errWrap(&err, "getMostRecentRank(ctx, %q)", puuid)

	rankStore := RankStore{Tx: db.Pool}

	ids, err := rankStore.ListRankIDs(ctx, puuid, ListRankOption{Limit: 1, Recent: true})
	if err != nil {
		return m, n, err
	}

	if len(ids) != 1 {
		return m, n, errors.New("ListRankIDS did not return exactly one id")
	}

	id := ids[0]

	return db.getRank(ctx, id)
}
