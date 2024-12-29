package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/http/response/html"
)

func (ui *ui) profileShowChampStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	puuid := chi.URLParam(r, "puuid")

	champs, err := ui.repo.GetChampionList(ctx, internal.PUUID(puuid))
	if err != nil {

		return
	}

	matches := make([]ProfileChampStatsModel, len(champs))

	for i, champion := range champs {
		matches[i] = ProfileChampStatsModel{
			Puuid: champion.Puuid,
		}
	}

	html.OK(w, r, ProfileChampStatsRows(matches))
}
