package ui

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/http/response/html"
)

func (ui *ui) profileShowLiveMatch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	puuid := internal.PUUID(chi.URLParam(r, "puuid"))

	livematch, err := ui.api.GetLiveMatch(ctx, puuid)
	if err != nil {
		html.ServerError(w, r, profileLiveMatchError(), err)
		return
	}

	models := [10]profileLiveMatchModel{}
	for i, puuid := range livematch.IDs {
		participant := livematch.GetParticipants()[i]

		profile, err := ui.repo.GetProfile(ctx, puuid)
		if err != nil {
			html.ServerError(w, r, profileLiveMatchError(), err)
			return
		}

		models[i] = profileLiveMatchModel{
			Puuid:          puuid,
			TeamID:         participant.TeamID,
			Date:           livematch.StartTimestamp,
			Name:           profile.Name,
			Tagline:        profile.Tagline,
			Champion:       participant.Champion,
			Runes:          participant.Runes,
			Summoners:      participant.Summoners,
			Rank:           profile.Rank,
			BannedChampion: participant.BannedChampion,
		}
	}

	html.OK(w, r, profileLiveMatchPartial(models))
}
