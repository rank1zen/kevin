package profile

import (
	"errors"
	"net/http"
	"time"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/riot"
)

type HistoryEntryHandler frontend.Handler

func (h *HistoryEntryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := frontend.GetMatchHistoryRequest{}

	switch r.Header.Get("Content-type") {
	default:
		req.Region = new(riot.Region)
		*req.Region = frontend.StrToRiotRegion(r.FormValue("region"))

		req.PUUID = riot.PUUID(r.FormValue("puuid"))

		if start, err := time.Parse(time.RFC3339, r.FormValue("start")); err == nil {
			req.StartTS = &start
		}

		if end, err := time.Parse(time.RFC3339, r.FormValue("end")); err == nil {
			req.EndTS = &end
		}
	}

	profileService := (*frontend.ProfileService)(h)

	storeMatches, err := profileService.GetMatchHistory(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, errors.New("service failure"))
		return
	}

	v := &HistoryEntryData{
		Date:      *req.StartTS,
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

	c := HistoryEntry(r.Context(), *v)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, errors.New("templ error"))
		return
	}
}
