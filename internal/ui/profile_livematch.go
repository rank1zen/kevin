package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/http/response/html"
)

func (ui *ui) profileShowLiveMatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	puuid := chi.URLParam(r, "puuid")

	livematch, err := ui.api.GetLiveMatch(ctx, internal.PUUID(puuid))
	if err != nil {
		html.ServerError(w, r, ProfileLiveMatchError(), err)
		return
	}

	models := [10]ProfileLiveMatchModel{}
	for i, puuid := range livematch.IDs {
		participant := livematch.GetParticipants()[i]

		profile, err := ui.repo.GetProfile(ctx, puuid)
		if err != nil {
			html.ServerError(w, r, ProfileLiveMatchError(), err)
			return
		}

		models[i] = ProfileLiveMatchModel{
			Puuid: puuid,
			Team:  participant.TeamID,
			Date:  livematch.StartTimestamp,
			Name:  profile.Name + "#" + profile.Tagline,
		}

	}

	html.OK(w, r, ProfileLiveMatch(models[:]))
}
