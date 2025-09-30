package frontend

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/page"
	"github.com/rank1zen/kevin/internal/view/profile"
)

type IndexService struct {
	Datasource *internal.Datasource

	ProfileService ProfileService
	SearchService  SearchService
}

func (s *IndexService) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /{$}", s.handleHomePage)

	router.HandleFunc("GET /{riotID}", s.handleProfilePage)
}

func (s *IndexService) handleHomePage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := fromCtx(ctx)

	region := r.FormValue("region")

	payload := slog.Group("payload", "region", region)

	riotRegion := convertStringToRiotRegion(region)

	v := page.HomePageData{
		Region: riotRegion,
	}

	if err := page.HomePage(ctx, v).Render(ctx, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed rendering", "err", err)
		return
	}
}

func (s *IndexService) handleSumonerPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := fromCtx(ctx)

	region := r.FormValue("region")

	payload := slog.Group("payload", "region", region)

	riotID := r.PathValue("riotID")
	name, tag, err := ParseRiotID(riotID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Debug("failed to resolve riot id", "err", err, payload)
		return
	}

	riotRegion := convertStringToRiotRegion(region)

	data, err := s.handler.GetSummonerPage(ctx, riotRegion, name, tag, time.UTC)
	if err != nil {
		if errors.Is(err, internal.ErrSummonerNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Info("summoner is not found", "name", name, "tag", tag, payload)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed service", "err", err, payload)
		return
	}

	championListCh := make(chan profile.ChampionListData)
	data.ChampionListCh = championListCh

	rankCardCh := make(chan profile.RankCardData)
	data.RankCardCh = rankCardCh

	historyEntryCh := make(chan profile.HistoryEntryData)
	data.HistoryEntryCh = historyEntryCh
	go func() {
		defer close(historyEntryCh)
		defer close(championListCh)
		defer close(rankCardCh)

		profileHandler := ProfileHandler{Datasource: s.handler.Datasource}
		days := GetDays(time.Now())
		for i := range len(days) - 1 {
			historyEntryData, err := profileHandler.GetMatchHistory(ctx, MatchHistoryRequest{
				Region:  riotRegion,
				PUUID:   data.PUUID,
				StartTS: days[i+1],
				EndTS:   days[i],
			})
			if err == nil {
				historyEntryCh <- *historyEntryData
			}
		}
	}()

	component := page.ProfilePage(ctx, *data)

	templ.Handler(component, templ.WithStreaming()).ServeHTTP(w, r)
}
