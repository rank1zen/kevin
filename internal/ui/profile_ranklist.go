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

	m, err := ui.repo.GetRankList(ctx, puuid)
	if err != nil {
		html.ServerError(w, r, ProfileRankListError(), err)
		return
	}

	html.OK(w, r, ProfileRankList(m))
}
