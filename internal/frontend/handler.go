package frontend

import (
	"context"
	"strings"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/component"
	"github.com/rank1zen/kevin/internal/component/shared"
	"github.com/rank1zen/kevin/internal/component/view"
	"github.com/rank1zen/kevin/internal/riot"
)

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

// GetHomePage returns the home page.
func (h *Handler) GetHomePage(ctx context.Context, region riot.Region) (component.Component, error) {
	v := shared.NewHomePage()
	return v, nil
}

// GetSummonerPage returns a summoner's profile page if summoner exists in
// store, otherwise, it will complete a update for summoner, then return the
// page. If no summoner with name#tag exists, return a does not exist page.
func (h *Handler) GetSummonerPage(ctx context.Context, region riot.Region, name, tag string, tz *time.Location) (*view.ProfilePage, error) {
	detail, err := h.Datasource.GetProfileDetailByRiotID(ctx, region, name, tag)
	if err != nil {
		return nil, err
	}

	v := view.ProfilePage{
		PUUID:     detail.PUUID,
		Name:      detail.Name,
		Tag:       detail.Tagline,
		Rank:      detail.Rank,
		Requests:  []view.MatchHistoryRequest{},
		LiveMatch: view.LiveMatchRequest{},
		Champion:  view.ChampionRequest{},
		Update:    view.UpdateSummonerRequest{},
	}

	var tmp []byte
	v.LiveMatch.Path, tmp = makeGetLiveMatch(region, detail.PUUID)
	v.LiveMatch.Data = string(tmp)

	v.Champion.Path, tmp = makeGetChampionList(region, detail.PUUID)
	v.Champion.Data = string(tmp)

	v.Update.Path, tmp = makeUpdateSummoner(region, detail.Name, detail.Tagline)
	v.Update.Data = string(tmp)

	dates := GetDays(time.Now().In(tz))
	for i := range len(dates) - 1 {
		path, data := makeGetMatchHistoryRequest(region, detail.PUUID, dates[i+1], dates[i])
		v.Requests = append(v.Requests, view.MatchHistoryRequest{Date: dates[i+1], Path: path, Data: string(data)})
	}

	return &v, nil
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

// GetCurrentWeek returns the start of the day, 7 days ago. Currently returns
// UTC time.
func GetCurrentWeek() time.Time {
	now := time.Now().In(time.UTC)
	y, m, d := now.Date()
	startOfDay := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	return startOfDay.Add(-24 * 6 * time.Hour)
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
