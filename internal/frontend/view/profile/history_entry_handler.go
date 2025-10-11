package profile

import (
	"net/http"
	"time"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/riot"
)

type HistoryEntryHandler frontend.Handler

func (h *HistoryEntryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := frontend.LoggerFromContext(ctx)

	puuid := riot.PUUID("aaa")
	start := time.Now()
	end := start.Add(12 * time.Second)

	storeMatches, err := h.Datasource.GetMatchHistory(ctx, riot.RegionNA1, puuid, start, end)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("datasource failed", "err", err)
		return
	}

	v := &HistoryEntryData{
		Date:      start,
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

	c := HistoryEntry(ctx, *v)

	if err := c.Render(ctx, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
