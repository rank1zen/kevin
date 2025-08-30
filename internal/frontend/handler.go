package frontend

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/component"
	"github.com/rank1zen/kevin/internal/component/shared"
	"github.com/rank1zen/kevin/internal/component/view"
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

type UpdateSummonerRequest struct {
	Region riot.Region `json:"region"`
	Name   string      `json:"name"`
	Tag    string      `json:"tag"`
}

func (r UpdateSummonerRequest) Validate() (problems map[string]string) {
	return nil
}

// UpdateSummoner syncs a summoner and their rank with riot.
func (h *Handler) UpdateSummoner(ctx context.Context, req UpdateSummonerRequest) error {
	if err := h.Datasource.UpdateProfileByRiotID(ctx, req.Region, req.Name, req.Tag); err != nil {
		return err
	}

	return nil
}

type LiveMatchRequest struct {
	Region riot.Region `json:"region"`
	PUUID  riot.PUUID  `json:"puuid"`
}

func (r LiveMatchRequest) Validate() (problems map[string]string) {
	problems = make(map[string]string)
	validatePUUID(problems, r.PUUID)
	return problems
}

// GetLiveMatch returns a live match view depending on whether the summoner is
// on game.
func (h *Handler) GetLiveMatch(ctx context.Context, req LiveMatchRequest) (component.Component, error) {
	match, err := h.Datasource.GetLiveMatch(ctx, req.Region, req.PUUID)
	if err != nil {
		if errors.Is(err, internal.ErrNoLiveMatch) {
			c := view.LiveMatchNotFound{}
			return c, nil
		}

		return nil, err
	}

	c := MapLiveMatch(match)

	return c, nil
}

// GetHomePage returns the home page.
func (h *Handler) GetHomePage(ctx context.Context, region riot.Region) (component.Component, error) {
	v := shared.NewHomePage()
	return v, nil
}

// GetSummonerPage returns a summoner's profile page if summoner exists in
// store, otherwise, it will complete a update for summoner, then return the
// page. If no summoner with name#tag exists, return a does not exist page.
func (h *Handler) GetSummonerPage(ctx context.Context, region riot.Region, name, tag string) (component.Component, error) {
	detail, err := h.Datasource.GetProfileDetailByRiotID(ctx, region, name, tag)
	if err != nil {
		return nil, err
	}

	mapper := FrontendToProfilePageMapper{
		Region:        region,
		ProfileDetail: detail,
	}

	c := mapper.Map()

	return c, nil
}

type GetSummonerChampionsRequest struct {
	Region riot.Region `json:"region"`
	PUUID  riot.PUUID  `json:"puuid"`
	Week   time.Time   `json:"week"`
}

func (r GetSummonerChampionsRequest) Validate() (problems map[string]string) {
	problems = make(map[string]string)

	validatePUUID(problems, r.PUUID)

	if r.Week.Hour() != 0 || r.Week.Minute() != 0 || r.Week.Second() != 0 {
		problems["date"] = "date needs to be start of day"
	}

	return problems
}

// GetSummonerChampions returns The method will fetch all
// games played in the specified interval.
func (h *Handler) GetSummonerChampions(ctx context.Context, req GetSummonerChampionsRequest) (component.Component, error) {
	start := req.Week
	end := start.Add(7 * 24 * time.Hour)

	storeChampions, err := h.Datasource.GetSummonerChampions(ctx, req.Region, req.PUUID, start, end)
	if err != nil {
		return nil, err
	}

	mapper := FrontendToSummonerChampstatMapper{
		Champions: storeChampions,
	}

	c := mapper.Map()

	return c, nil
}

// GetMatchHistory returns [MatchHistory], the matches played on date to date
// + 24 hours. The method will fetch riot first to ensure all matches played on
// date are in store.
func (h *Handler) GetMatchHistory(ctx context.Context, req MatchHistoryRequest) (component.Component, error) {
	// TODO: consider daylight savings
	end := req.Date.Add(24 * time.Hour)

	storeMatches, err := h.Datasource.GetMatchHistory(ctx, req.Region, req.PUUID, req.Date, end)
	if err != nil {
		return nil, fmt.Errorf("storage failure: %w", err)
	}

	mapper := FrontendToHistoryMapper{
		Region:       req.Region,
		MatchHistory: storeMatches,
	}

	c := mapper.Map()

	return c, nil
}

// GetSearchResults returns a list of [SearchResultCard] for accounts that
// match q. If no results were found, return [SearchNotFoundCard] instead. If
// no tag was provided in q, region will be used.
func (h *Handler) GetSearchResults(ctx context.Context, region riot.Region, q string) (component.Component, error) {
	storeSearchResults, err := h.Datasource.Search(ctx, region, q)
	if err != nil {
		return nil, err
	}

	if len(storeSearchResults) == 0 {
		name, tag := GetNameTag(q)
		if tag == "" {
			tag = string(region)
		}

		v := NewSearchNotFoundCard(region, name, tag)

		return v, nil
	}

	mapper := FrontendToSearchResultMapper{
		Region:  region,
		Results: storeSearchResults,
	}

	c := mapper.Map()

	return c, nil
}

type MatchDetailRequest struct {
	Region  riot.Region `json:"region"`
	MatchID string
}

func (r MatchDetailRequest) Validate() (problems map[string]string) {
	return nil
}

// GetMatchDetail returns [MatchDetail].
func (h *Handler) GetMatchDetail(ctx context.Context, req MatchDetailRequest) (component.Component, error) {
	matchDetail, err := h.Datasource.GetMatchDetail(ctx, req.Region, req.MatchID)
	if err != nil {
		return nil, err
	}

	mapper := FrontendToMatchDetailMapper{
		MatchDetail: matchDetail,
	}

	c := mapper.Map()

	return c, nil
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

// GetNameTag extracts the name and tag from query, which should be of the form
// name#tag.
func GetNameTag(query string) (name, tag string) {
	if i := strings.Index(query, "#"); i != -1 {
		name = query[:i]
		if i+1 == len(query) {
			return name, ""
		}

		tag = query[i+1:]
		return name, tag
	}

	return query, ""
}

func RoundToNearestInt(x float32) int {
	return int(x)
}

func validatePUUID(problems map[string]string, puuid riot.PUUID) {
	if len(puuid) != 78 {
		problems["puuid"] = "puuid is invalid"
	}
}
