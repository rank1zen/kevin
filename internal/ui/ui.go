package ui

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rank1zen/yujin/internal"
)

type ui struct {
	repo internal.Repository
	api  internal.RiotClient
}

func Routes(repo internal.Repository, api internal.RiotClient) *chi.Mux {
	base := chi.NewMux()

	base.Use(addRequestID, logMeta)
	base.Use(middleware.NoCache)

	base.NotFound(notFound)

	handler := ui{repo: repo, api: api}

	base.Get("/profile/{puuid}", handler.profileShow)
	base.Post("/profile/{puuid}/update", handler.profileRefresh)

	base.Get("/partials/profile/{puuid}/matchlist", handler.profileShowMatchList)
	base.Get("/partials/profile/{puuid}/champstats", handler.profileShowChampStats)
	base.Get("/partials/profile/{puuid}/ranklist", handler.profileShowRankList)
	base.Get("/partials/profile/{puuid}/live", handler.profileShowLiveMatch)

	return base
}
