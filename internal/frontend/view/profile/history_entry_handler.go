package profile

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/service"
)

type HistoryEntryHandler service.Service

func (h *HistoryEntryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := service.GetMatchlistRequest{}

	switch r.Header.Get("Content-type") {
	default:
		if region := r.FormValue("region"); region != "" {
			region := frontend.StrToRiotRegion(region)
			req.Region = &region
		}

		req.PUUID = riot.PUUID(r.FormValue("puuid"))

		if start, err := time.Parse(time.RFC3339, r.FormValue("startTs")); err == nil {
			req.StartTS = &start
		}

		if end, err := time.Parse(time.RFC3339, r.FormValue("endTs")); err == nil {
			req.EndTS = &end
		}
	}

	storeMatches, err := (*service.MatchService)(h).GetMatchlist(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, fmt.Errorf("service failure: %w", err))
		return
	}

	v := HistoryEntryData{
		Date:      time.Now(), // TODO: GetMatchHistory should return this data
		Matchlist: []HistoryCardData{},
	}

	for _, match := range storeMatches {
		// HACK: very hacky; please integrate into internal
		kda := float32(match.Kills+match.Assists) / float32(match.Deaths)

		v.Matchlist = append(v.Matchlist, HistoryCardData{
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
