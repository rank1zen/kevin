package frontend

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

type ZGetSummonerChampionsRequest struct {
	Region riot.Region `json:"region"`
	PUUID  riot.PUUID  `json:"puuid"`
	Week   time.Time   `json:"week"`
}

func (r ZGetSummonerChampionsRequest) Validate() (problems map[string]string) {
	problems = make(map[string]string)

	validatePUUID(problems, r.PUUID)

	if r.Week.Hour() != 0 || r.Week.Minute() != 0 || r.Week.Second() != 0 {
		problems["date"] = "date needs to be start of day"
	}

	return problems
}

type GetLiveMatchRequest struct {
	Region riot.Region `json:"region"`
	PUUID  riot.PUUID  `json:"puuid"`
}

func (r GetLiveMatchRequest) Validate() (problems map[string]string) {
	problems = make(map[string]string)
	validatePUUID(problems, r.PUUID)
	return problems
}

type MatchHistoryRequest struct {
	Region riot.Region `json:"region"`

	PUUID riot.PUUID `json:"puuid"`

	// Date should be the start of the day. The request will fetch all
	// matches played on the day.
	Date time.Time `json:"date"`
}

func (r MatchHistoryRequest) Validate() (problems map[string]string) {
	problems = make(map[string]string)

	validatePUUID(problems, r.PUUID)

	if r.Date.Hour() != 0 || r.Date.Minute() != 0 || r.Date.Second() != 0 {
		problems["date"] = "date needs to be start of day"
	}

	return problems
}

// Handler provides the API for server operations.
type Handler struct {
	Datasource *internal.Datasource
}

func NewHandler(datasource *internal.Datasource) *Handler {
	return &Handler{datasource}
}

func (h *Handler) CheckHealth(ctx context.Context) error {
	_, err := h.Datasource.GetPUUID(ctx, "orrange", "NA1")
	return err
}

// UpdateSummoner
func (h *Handler) UpdateSummoner(ctx context.Context, region riot.Region, name, tag string) error {
	puuid, err := h.Datasource.GetPUUID(ctx, name, tag)
	if err != nil {
		return err
	}

	if err := h.Datasource.UpdateSummoner(ctx, region, puuid); err != nil {
		return err
	}

	return nil
}

// GetLiveMatch returns [LiveMatchModalWindow] if the summoner is in game,
// otherwise [NoLiveMatchModalWindow].
func (h *Handler) GetLiveMatch(ctx context.Context, req GetLiveMatchRequest) (templ.Component, error) {
	match, err := h.Datasource.GetLiveMatch(ctx, req.Region, req.PUUID)
	if err != nil {
		if errors.Is(err, internal.ErrNoLiveMatch) {
			return NoLiveMatchModalWindow{}, err
		}

		return nil, err
	}

	blueSide, redSide := [5]LiveMatchRowLayout{}, [5]LiveMatchRowLayout{}

	for i, p := range match.Participants {
		name, tag, err := h.Datasource.GetRiotName(ctx, riot.PUUID(p.PUUID))
		if err != nil {
			return nil, fmt.Errorf("fetching participant %s: %w", p.PUUID, err)
		}

		card := LiveMatchRowLayout{
			MatchID:   p.MatchID,
			TeamID:    p.TeamID,
			PUUID:     riot.PUUID(p.PUUID),
			Champion:  p.ChampionID,
			Summoners: p.SummonersIDs,
			RunePage:  p.Runes,
			Name:      name,
			Tag:       tag,
			Rank:      nil, // TODO: fetch the rank.
		}

		if p.TeamID == 100 {
			blueSide[i % 5] = card
		} else {
			redSide[i % 5] = card
		}
	}

	window := LiveMatchModalWindow{
		AverageRank: &internal.RankDetail{},
		StartTime:   match.Date,
		RedSide:     redSide,
		BlueSide:    blueSide,
	}

	return window, nil
}

// GetHomePage returns [HomePage].
func (h *Handler) GetHomePage(ctx context.Context, region riot.Region) (templ.Component, error) {
	page := HomePage{Region: region}
	return page, nil
}

// GetSummonerPage returns [SummonerPage] if summoner exists in store,
// otherwise, it will complete a update for summoner, then return
// [SummonerPage]. If no summoner with name#tag exists, return
// [NoSummonerPage].
func (h *Handler) GetSummonerPage(ctx context.Context, region riot.Region, name, tag string) (templ.Component, error) {
	puuid, err := h.Datasource.GetPUUID(ctx, name, tag)
	if err != nil {
		if errors.Is(err, internal.ErrSummonerDoesNotExist) {
			return NoSummonerPage{
				Region: region,
				Name:   name,
				Tag:    tag,
			}, nil
		}

		return nil, fmt.Errorf("get puuid: %w", err)
	}

	summoner, err := h.Datasource.GetStore().GetSummoner(ctx, puuid)
	if err != nil {
		if errors.Is(err, internal.ErrSummonerNotFound) {
			fromCtx(ctx).Info("first time visit", "puuid", puuid)

			if err := h.Datasource.UpdateSummoner(ctx, region, puuid); err != nil {
				return nil, fmt.Errorf("getting puuid: %w", err)
			}
		}
	}

	// try again
	summoner, err = h.Datasource.GetStore().GetSummoner(ctx, puuid)
	if err != nil {
		return nil, fmt.Errorf("get summoner: %w", err)
	}

	rank, err := h.Datasource.GetStore().GetRank(ctx, puuid, time.Now(), true)
	if err != nil {
		return nil, fmt.Errorf("get rank: %w", err)
	}

	page := SummonerPage{
		Region:      region,
		PUUID:       puuid,
		Name:        summoner.Name,
		Tag:         summoner.Tagline,
		Rank:        rank.Detail,
		LastUpdated: rank.EffectiveDate,

		GetChampionsRequest: ZGetSummonerChampionsRequest{
			Region: region,
			PUUID:  puuid,
			Week:   GetCurrentWeek(),
		},
	}

	return page, nil
}

// GetSummonerChampions returns [SummonerChampionList] consisting of stats for
// the 7 days starting at week. The method will fetch all games played in the
// specified interval.
func (h *Handler) ZGetSummonerChampions(ctx context.Context, req ZGetSummonerChampionsRequest) (templ.Component, error) {
	start := req.Week
	end := start.Add(7 * 24 * time.Hour)

	err := h.Datasource.ZUpdateMatchHistory(ctx, req.Region, req.PUUID, start, end)
	if err != nil {
		return nil, fmt.Errorf("updating matchlist failed: %w", err)
	}

	storeChampions, err := h.Datasource.GetStore().GetChampions(ctx, req.PUUID, start, end)
	if err != nil {
		return nil, err
	}

	totalGamesPlayed := 0
	for _, c := range storeChampions {
		totalGamesPlayed += c.GamesPlayed
	}

	list := SummonerChampionList{Champions: []ChampionPopover{}}

	for _, champion := range storeChampions {
		list.Champions = append(
			list.Champions,
			ChampionPopover{
				Champion:             int(champion.Champion),
				TotalGamesPlayed:     totalGamesPlayed,
				GamesPlayed:          champion.GamesPlayed,
				Wins:                 champion.Wins,
				Losses:               champion.Losses,
				Kills:                champion.Kills,
				Deaths:               champion.Deaths,
				Assists:              champion.Assists,
				KillParticipation:    champion.KillParticipation,
				CS:                   champion.CreepScore,
				CSPerMinute:          champion.CreepScorePerMinute,
				DamageDealt:          champion.DamageDealt,
				DamageTaken:          champion.DamageTaken,
				DamageDeltaEnemy:     champion.DamageDeltaEnemy,
				DamagePercentageTeam: champion.DamagePercentageTeam,
				GoldEarned:           champion.GoldEarned,
				GoldDeltaEnemy:       champion.GoldDeltaEnemy,
				GoldPercentageTeam:   champion.GoldPercentageTeam,
				VisionScore:          champion.VisionScore,
				PinkWardsBought:      champion.PinkWardsBought,
				LPGain:               0,
			},
		)
	}

	return list, nil
}

// GetMatchHistory returns [MatchHistoryList], being the matches played on date
// to date + 24 hours. The method will fetch riot first to ensure all matches
// played on date are in store.
func (h *Handler) GetMatchHistory(ctx context.Context, req MatchHistoryRequest) (templ.Component, error) {
	// TODO: consider daylight savings
	end := req.Date.Add(24*time.Hour)

	if err := h.Datasource.ZUpdateMatchHistory(ctx, req.Region, req.PUUID, req.Date, end); err != nil {
		return nil, err
	}

	storeMatches, err := h.Datasource.GetStore().GetZMatches(ctx, req.PUUID, req.Date, end)
	if err != nil {
		return nil, fmt.Errorf("storage failure: %w", err)
	}

	list := MatchHistoryList{Matches: []MatchHistoryRow{}}

	for _, m := range storeMatches {
		list.Matches = append(
			list.Matches,
			MatchHistoryRow{
				MatchID:     m.MatchID,
				Champion:    m.ChampionID,
				Summoners:   m.SummonerIDs,
				Kills:       m.Kills,
				Deaths:      m.Deaths,
				Assists:     m.Assists,
				CS:          m.CreepScore,
				CSPerMinute: m.CreepScorePerMinute,
				RunePage:    m.Runes,
				Items:       m.Items,
				RankChange:  nil, // TODO: cannot fill these fields yet
				LPChange:    nil,
			},
		)
	}

	return list, nil
}

// GetMatchScoreboard returns the scoreboard of a match.
func (h *Handler) GetMatchScoreboard(ctx context.Context, id string) (scoreboard templ.Component, err error) {
	return nil, nil
}

// GetSearchResults returns a list of [SearchResultCard] for accounts that
// match q. If no results were found, return [SearchNotFoundCard] instead.
//
// q should be of the form name#tag, if q has no tag, region is used as the
// tag.
func (h *Handler) GetSearchResults(ctx context.Context, region riot.Region, q string) (templ.Component, error) {
	storeSearchResults, err := h.Datasource.GetStore().SearchSummoner(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(storeSearchResults) == 0 {
		var name, tag string
		if i := strings.Index(q, "#"); i != -1 {
			name = q[:i]
			if i+1 == len(q) {
				tag = string(region)
			} else {
				tag = q[i+1:]
			}
		} else {
			name = q
			tag = string(region)
		}

		return SearchNotFoundCard{
			Name:     name,
			Tag:      tag,
			Platform: string(region),
		}, nil
	}

	searchResults := []SearchResultLink{}

	for _, r := range storeSearchResults {
		rank, err := h.Datasource.GetStore().GetRank(ctx, r.Puuid, time.Now(), true)
		if err != nil {
			return nil, fmt.Errorf("getting rank for %s#%s: %w", r.Name, r.Tagline, err)
		}

		row := SearchResultLink{
			Region: region,
			PUUID:  r.Puuid,
			Name:   r.Name,
			Tag:    r.Tagline,
			Rank:   rank.Detail,
		}

		searchResults = append(searchResults, row)
	}

	v := []templ.Component{}
	for _, r := range searchResults {
		v = append(v, r)
	}

	return templ.Join(v...), nil
}

// GetCurrentWeek returns the start of the day, 7 days ago.
func GetCurrentWeek() time.Time {
	now := time.Now().In(time.UTC)
	y, m, d := now.Date()
	startOfDay := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	return startOfDay.Add(-24 * 6 * time.Hour)
}

func validatePUUID(problems map[string]string, puuid riot.PUUID) {
	if len(puuid) != 78 {
		problems["puuid"] = "puuid is invalid"
	}
}
