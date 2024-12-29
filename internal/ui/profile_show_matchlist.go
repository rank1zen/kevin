package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/http/request"
	"github.com/rank1zen/yujin/internal/http/response/html"
)

func (ui *ui) profileShowMatchList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	puuid := internal.PUUID(chi.URLParam(r, "puuid"))
	page := request.QueryIntParam(r, "page", 0)

	newIDs, err := ui.api.GetMatchList(ctx, puuid)
	if err != nil {
		html.ServerError(w, r, profileMatchListError(), err)
		return
	}

	for _, id := range newIDs {
		ui.api.GetMatch(ctx, internal.MatchID(id))
	}
	matches, err := ui.repo.GetMatchList(ctx, puuid, page, true)
	if err != nil {
		html.ServerError(w, r, profileMatchListError(), err)
		return
	}

	models := make([]profileMatchModel, len(matches))

	for i, match := range matches {
		models[i] = profileMatchModel{
			MatchID: match.Match,
			TeamID:  match.Team,
			Patch:   match.Patch,
			Name:    match.RiotIDGameName + "#" + match.RiotIDTagline,
			Kills:   match.Kills,
			Deaths:  match.Deaths,
			Assists: match.Assists,
		}
	}

	html.OK(w, r, profileMatchListPartial(models))
}