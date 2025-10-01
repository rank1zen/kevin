package frontend

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/view/search"
)

type SearchService struct {
	Handler *SearchHandler
}

func (s *SearchService) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /search", s.serveSearchResults)
}

func (s *SearchService) serveSearchResults(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := fromCtx(ctx)

	var (
		region = r.FormValue("region")
		q      = r.FormValue("q")
	)

	payload := slog.Group("payload", "region", region, "q", q)

	riotRegion := convertStringToRiotRegion(region)

	data, err := s.Handler.GetSearchResults(ctx, riotRegion, q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed service", slog.Any("err", err), payload)
		return
	}

	if err := search.ResultMenu(ctx, *data).Render(ctx, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed rendering", "err", err)
		return
	}
}

type SearchHandler struct {
	Datasource *internal.Datasource
}

func (h *SearchHandler) GetSearchResults(ctx context.Context, region riot.Region, q string) (*search.ResultMenuData, error) {
	storeSearchResults, err := h.Datasource.Search(ctx, region, q)
	if err != nil {
		return nil, err
	}

	var name, tag string

	name, tag = GetNameTag(q)
	if tag == "" {
		tag = string(region)
	}

	data := search.ResultMenuData{
		Name:           name,
		Tag:            tag,
		ProfileResults: []search.ResultCardData{},
	}

	for _, result := range storeSearchResults {
		var rank *internal.Rank
		if result.Rank != nil && result.Rank.Detail != nil {
			rank = &result.Rank.Detail.Rank
		}

		data.ProfileResults = append(data.ProfileResults, search.ResultCardData{
			Name: result.Name,
			Tag:  result.Tagline,
			Rank: rank,
			Path: getSummonerPath(result.Name, result.Tagline),
		})
	}

	return &data, nil
}

func getSummonerPath(name, tag string) string {
	return fmt.Sprintf("/%s-%s", name, tag)
}
