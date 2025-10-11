package profile

import (
	"log/slog"
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/view/profile"
)

type MatchDetailBoxHandler frontend.Handler

func (h *MatchDetailBoxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	v := &profile.HistoryEntryData{
		Date:      start,
		Matchlist: []profile.HistoryCardData{},
	}

	for _, match := range storeMatches {
		path, data := makeGetMatchDetailRequest(req.Region, match.MatchID)

		kda := float32(match.Kills+match.Assists) / float32(match.Deaths)

		v.Matchlist = append(v.Matchlist, profile.HistoryCardData{
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
			Path:           path,
			Data:           string(data),
		})
	}

	c := profile.HistoryEntry(ctx, *v)

	if err := c.Render(ctx, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
