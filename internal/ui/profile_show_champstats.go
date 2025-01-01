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

	matches := make([]profileChampStatsModel, len(champs.List))

	for i, champion := range champs.List {
		matches[i] = profileChampStatsModel{
			Puuid:             champion.Puuid,
			Champion:          champion.Champion,
			GamesPlayed:       champion.GamesPlayed,
			Wins:              champion.Wins,
			Losses:            champion.Losses,
			Kills:             champion.Kills,
			Deaths:            champion.Deaths,
			Assists:           champion.Assists,
			WinPercentage:     champion.WinPercentage,
			LpDelta:           champion.LpDelta,
			KillParticipation: champion.KillParticipation,
			CreepScore:        champion.CreepScore,
		}
	}

	html.OK(w, r, profileChampStatsPartial(matches))
}
