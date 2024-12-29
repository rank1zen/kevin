package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/http/response/html"
)

func (ui *ui) profileShowRankList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	puuid := internal.PUUID(chi.URLParam(r, "puuid"))

	ranklist, err := ui.repo.GetRankList(ctx, puuid)
	if err != nil {
		html.ServerError(w, r, profileRankListError(), err)
		return
	}

	models := make([]profileRankModel, len(ranklist))
	for i := range len(ranklist) {
		models[i] = profileRankModel{}
	}

	html.OK(w, r, profileRankListPartial(models))
}
