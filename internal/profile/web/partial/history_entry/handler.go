package history_entry

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/frontend/view/historycard"
	"github.com/rank1zen/kevin/internal/profile"
)

type Handler struct {
	service profile.ProfileService
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := &profile.GetMatchlistRequest{}

	switch r.Header.Get("Content-type") {
	default:
		req.PUUID = r.FormValue("puuid")

		if start, err := time.Parse(time.RFC3339, r.FormValue("startTs")); err == nil {
			req.StartTS = &start
		}

		if end, err := time.Parse(time.RFC3339, r.FormValue("endTs")); err == nil {
			req.EndTS = &end
		}
	}

	storeMatches, err := h.service.GetMatchlist(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, fmt.Errorf("service failure: %w", err))
		return
	}

	v := HistoryEntryData{
		Date:      time.Now(), // TODO: GetMatchHistory should return this data
		Matchlist: []historycard.HistorycardData{},
	}

	for _, match := range storeMatches {
		// HACK: very hacky; please integrate into internal
		kda := float32(match.Kills+match.Assists) / float32(match.Deaths)

		v.Matchlist = append(v.Matchlist, historycard.HistorycardData{
			ChampionID:     match.ChampionID,
			ChampionLevel:  match.ChampionLevel,
			SummonerIDs:    match.SummonerIDs,
			Kills:          match.Kills,
			Deaths:         match.Deaths,
			Assists:        match.Assists,
			KillDeathRatio: kda,
			CS:             match.CreepScore,
			CSPerMinute:    match.CreepScorePerMinute,
			RunePage:       match.Runes,
			Items:          match.Items,
			VisionScore:    match.VisionScore,
			RankChange:     nil,
			LPChange:       nil,
			Win:            match.Win,
		})
	}

	c := HistoryEntry(r.Context(), v)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, fmt.Errorf("templ error: %w", err))
		return
	}
}
