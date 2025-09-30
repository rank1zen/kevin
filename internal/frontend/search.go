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

	if err := search.ResultList(ctx, *data).Render(ctx, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed rendering", "err", err)
		return
	}
}

type SearchHandler struct {
	Datasource *internal.Datasource
}

func (h *SearchHandler) GetSearchResults(ctx context.Context, region riot.Region, q string) (*search.ResultListData, error) {
	storeSearchResults, err := h.Datasource.Search(ctx, region, q)
	if err != nil {
		return nil, err
	}

	data := search.ResultListData{
		Cards: []search.ResultCardData{},
	}

	for _, result := range storeSearchResults {
		var rank *internal.Rank
		if result.Rank != nil && result.Rank.Detail != nil {
			rank = &result.Rank.Detail.Rank
		}

		data.Cards = append(data.Cards, search.ResultCardData{
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
