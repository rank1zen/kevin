package ui

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/http/response/html"
)

func (ui *ui) profileRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	puuid := internal.PUUID(chi.URLParam(r, "puuid"))

	profile, err := ui.api.GetProfile(ctx, puuid)
	if err != nil {
		html.ServerError(w, r, nil, nil)
		return
	}

	err = ui.repo.UpdateProfile(ctx, profile)
	if err != nil {
		html.ServerError(w, r, nil, fmt.Errorf("updating summary: %w", err))
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusAccepted)
}
