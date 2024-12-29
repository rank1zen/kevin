package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/http/response/html"
)

func (ui *ui) profileShow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	puuid := internal.PUUID(chi.URLParam(r, "puuid"))

	exists, err := ui.repo.CheckProfileExists(ctx, puuid)
	if err != nil {
		html.ServerError(w, r, ServerError(), err)
		return
	}

	if !exists {
		html.BadRequest(w, r, ProfileDoesNotExist(), err)
		return
	}

	profile, err := ui.repo.GetProfile(ctx, puuid)
	if err != nil {
		html.ServerError(w, r, ServerError(), err)
		return
	}

	m := ProfileModel{
		Puuid:   profile.Puuid,
		Name:    profile.Name,
		Tagline: profile.Tagline,
		Rank:    profile.Rank,
	}

	html.OK(w, r, ProfilePage(m))

	return
}
