package frontend

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/page"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/view/search"
	"github.com/rank1zen/kevin/internal/view/shared"
)

func Error(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(200)

	logger := fromCtx(ctx)

	logger.Debug("failed service")
	shared.Error(ctx, shared.ErrorData{
		StatusCode: 0,
		Header:     "An error has occurred",
		Message:    "",
	})
}

// Handler provides the API for server operations.
type Handler struct {
	// Datasource handles the business logic. A nil value indicates that
	// Handler will use the zero value [internal.Datasource].
	Datasource *internal.Datasource
}

func (h *Handler) GetHomePage(ctx context.Context, region riot.Region) (*page.HomePageData, error) {
	v := page.HomePageData{
		Region: region,
	}

	return &v, nil
}

// GetSummonerPage returns a summoner's profile page if summoner exists in
// store, otherwise, it will complete a update for summoner, then return the
// page. If no summoner with name#tag exists, return a does not exist page.
func (h *Handler) GetSummonerPage(ctx context.Context, region riot.Region, name, tag string, tz *time.Location) (*page.ProfilePageData, error) {
	storeProfile, err := h.Datasource.GetProfileDetailByRiotID(ctx, region, name, tag)
	if err != nil {
		return nil, err
	}

	v := page.ProfilePageData{
		PUUID:  storeProfile.PUUID,
		Region: region,
		Name:   name,
		Tag:    tag,
	}

	return &v, nil
}

func (h *Handler) GetSearchResults(ctx context.Context, region riot.Region, q string) (templ.Component, error) {
	storeSearchResults, err := h.Datasource.Search(ctx, region, q)
	if err != nil {
		return nil, err
	}

	if len(storeSearchResults) == 0 {
		name, tag := GetNameTag(q)
		if tag == "" {
			tag = string(region)
		}

		data := search.NotFoundCardData{
			Region: region,
			Name:   name,
			Tag:    tag,
			Path:   "",
			Data:   "",
		}

		return search.NotFoundCard(ctx, data), nil
	}

	data := search.ResultListData{
		Cards: []search.ResultCardData{},
	}

	for _, result := range storeSearchResults {
		data.Cards = append(data.Cards, search.ResultCardData{
			Name: result.Name,
			Tag:  result.Tagline,
			Rank: nil,
			Path: "",
		})
	}

	return search.ResultList(ctx, data), nil
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
