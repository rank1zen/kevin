package profile

import (
	"log/slog"
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
)

type ChampionListHandler frontend.Handler

func (h *ChampionListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	region := r.FormValue("region")

	payload := slog.Group("payload", "region", region)

	req, err := decode[MatchHistoryRequest](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Debug("bad request", "err", err)
		return
	}

	storeMatches, err := h.Datasource.GetMatchHistory(ctx, region, puuid)

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
