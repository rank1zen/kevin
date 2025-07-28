package frontend

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

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
	// Datasource handles the business logic. A nil value indicates that
	// Handler will use the zero value [internal.Datasource].
	Datasource *internal.Datasource
}

// CheckHealth does a simple request to ensure systems are working.
func (h *Handler) CheckHealth(ctx context.Context) error {
	ds := h.Datasource
	if ds == nil {
		ds = &internal.Datasource{}
	}

	_, err := ds.GetPUUID(ctx, "orrange", "NA1")
	return err
}

// UpdateSummoner syncs a summoner and their rank with riot.
func (h *Handler) UpdateSummoner(ctx context.Context, region riot.Region, name, tag string) error {
	ds := h.Datasource
	if ds == nil {
		ds = &internal.Datasource{}
	}

	puuid, err := ds.GetPUUID(ctx, name, tag)
	if err != nil {
		return err
	}

	if err := ds.UpdateSummoner(ctx, region, puuid); err != nil {
		return err
	}

	return nil
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

// GetLiveMatch returns [LiveMatch] if the summoner is in game, otherwise it
// returns [LiveMatchNotFound].
func (h *Handler) GetLiveMatch(ctx context.Context, req GetLiveMatchRequest) (View, error) {
	ds := h.Datasource
	if ds == nil {
		ds = &internal.Datasource{}
	}

	match, err := ds.GetLiveMatch(ctx, req.Region, req.PUUID)
	if err != nil {
		if errors.Is(err, internal.ErrNoLiveMatch) {
			v := LiveMatchNotFound{}
			return v, err
		}

		return nil, err
	}

	v := LiveMatch{}
	return v, nil
}

// GetHomePage returns [HomePage].
func (h *Handler) GetHomePage(ctx context.Context, region riot.Region) (View, error) {
	v := HomePage{}
	return v, nil
}

// GetSummonerPage returns [SummonerPage] if summoner exists in store,
// otherwise, it will complete a update for summoner, then return the page. If
// no summoner with name#tag exists, return [SummonerPageNotFound].
func (h *Handler) GetSummonerPage(ctx context.Context, region riot.Region, name, tag string) (View, error) {
	ds := h.Datasource
	if ds == nil {
		ds = &internal.Datasource{}
	}

	puuid, err := ds.GetPUUID(ctx, name, tag)
	if err != nil {
		if errors.Is(err, internal.ErrSummonerDoesNotExist) {
			v := SummonerPageNotFound{
				Region: region,
				Name:   name,
				Tag:    tag,
			}

			return v, nil
		}

		return nil, fmt.Errorf("get puuid: %w", err)
	}

	summoner, err := ds.GetStore().GetSummoner(ctx, puuid)
	if err != nil {
		if errors.Is(err, internal.ErrSummonerNotFound) {
			fromCtx(ctx).Info("first time visit", "puuid", puuid)

			if err := ds.UpdateSummoner(ctx, region, puuid); err != nil {
				return nil, fmt.Errorf("getting puuid: %w", err)
			}
		}
	}

	summoner, err = ds.GetStore().GetSummoner(ctx, puuid)
	if err != nil {
		return nil, fmt.Errorf("get summoner: %w", err)
	}

	rank, err := ds.GetStore().GetRank(ctx, puuid, time.Now(), true)
	if err != nil {
		return nil, fmt.Errorf("get rank: %w", err)
	}

	v := SummonerPage{
		Region:              region,
		PUUID:               puuid,
		Name:                name,
		Tag:                 tag,
		LastUpdated:         rank.EffectiveDate,
		Rank:                &rank.Detail.Rank,
	}

	return v, nil
}

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

// GetSummonerChampions returns [SummonerChampion]. The method will fetch all
// games played in the specified interval.
func (h *Handler) ZGetSummonerChampions(ctx context.Context, req ZGetSummonerChampionsRequest) (View, error) {
	ds := h.Datasource
	if ds == nil {
		ds = &internal.Datasource{}
	}

	start := req.Week
	end := start.Add(7 * 24 * time.Hour)

	err := ds.ZUpdateMatchHistory(ctx, req.Region, req.PUUID, start, end)
	if err != nil {
		return nil, fmt.Errorf("updating matchlist failed: %w", err)
	}

	storeChampions, err := ds.GetStore().GetChampions(ctx, req.PUUID, start, end)
	if err != nil {
		return nil, err
	}

	v := SummonerChampion{}
	return v, nil
}

// GetMatchHistory returns [MatchHistory], the matches played on date to date
// + 24 hours. The method will fetch riot first to ensure all matches played on
// date are in store.
func (h *Handler) GetMatchHistory(ctx context.Context, req MatchHistoryRequest) (View, error) {
	ds := h.Datasource
	if ds == nil {
		ds = &internal.Datasource{}
	}

	// TODO: consider daylight savings
	end := req.Date.Add(24 * time.Hour)

	if err := ds.ZUpdateMatchHistory(ctx, req.Region, req.PUUID, req.Date, end); err != nil {
		return nil, err
	}

	storeMatches, err := ds.GetStore().GetZMatches(ctx, req.PUUID, req.Date, end)
	if err != nil {
		return nil, fmt.Errorf("storage failure: %w", err)
	}

	v := MatchHistory{}

	return v, nil
}

// GetSearchResults returns a list of [SearchResultCard] for accounts that
// match q. If no results were found, return [SearchNotFoundCard] instead.
//
// q should be of the form name#tag, if q has no tag, region is used as the
// tag.
func (h *Handler) GetSearchResults(ctx context.Context, region riot.Region, q string) (View, error) {
	ds := h.Datasource
	if ds == nil {
		ds = &internal.Datasource{}
	}

	storeSearchResults, err := ds.GetStore().SearchSummoner(ctx, q)
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

		v := SearchResultNotFound{}

		return v, nil
	}

	v := SearchResult{}

	return v, nil
}

type GetMatchDetailsRequest struct {
	MatchID string
}

// GetMatchDetails returns [MatchDetail].
func (h *Handler) GetMatchDetails(ctx context.Context, req GetMatchDetailsRequest) (View, error) {
	match, err := h.Datasource.GetStore().GetMatch(ctx, riot.PUUID(req.MatchID))
	if err != nil {
		return nil, err
	}

	v := MatchDetail{}

	return v, nil
}

// GetCurrentWeek returns the start of the day, 7 days ago. Currently returns
// UTC time.
func GetCurrentWeek() time.Time {
	now := time.Now().In(time.UTC)
	y, m, d := now.Date()
	startOfDay := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	return startOfDay.Add(-24 * 6 * time.Hour)
}

// GetDay returns the start of the day, offset days ago. Currently returns UTC
// time.
func GetDay(offset int) time.Time {
	now := time.Now().In(time.UTC)
	y, m, d := now.Date()
	startOfDay := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	return startOfDay.Add(time.Duration(-24*offset) * time.Hour)
}

func ComputeFraction(wins, losses int) float32 {
	return float32(wins) / float32(losses)
}

func RoundToNearestInt(x float32) int {
	return int(x)
}

func validatePUUID(problems map[string]string, puuid riot.PUUID) {
	if len(puuid) != 78 {
		problems["puuid"] = "puuid is invalid"
	}
}
