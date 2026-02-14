package history_entry

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/profile"
)

type Handler struct {
	service *profile.ProfileService
}

func NewHandler(service *profile.ProfileService) *Handler {
	return &Handler{
		service: service,
	}
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

	c := Index(r.Context(), mapper(storeMatches))

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, fmt.Errorf("templ error: %w", err))
		return
	}
}

func mapper(matches []profile.Match) *IndexData {
	mapp := &IndexData{
		Date:      time.Now(),
		Matchlist: []CardData{},
	}
	for _, match := range matches {
		kda := float32(match.Kills+match.Assists) / float32(match.Deaths)

		var lpChange string

		mapp.Matchlist = append(mapp.Matchlist, CardData{
			Kills:                      strconv.Itoa(match.Kills),
			Deaths:                     strconv.Itoa(match.Deaths),
			Assists:                    strconv.Itoa(match.Assists),
			KillDeathRatio:             strconv.FormatFloat(float64(kda), 'f', 2, 64),
			CS:                         strconv.Itoa(match.CreepScore),
			CSPerMinute:                strconv.FormatFloat(float64(match.CreepScorePerMinute), 'f', 2, 64),
			VisionScore:                strconv.Itoa(match.VisionScore),
			LPChange:                   lpChange,
			Win:                        match.Win,
			Date:                       "",
			Duration:                   "",
			Rank:                       "",
			ChampionImagePath:          "",
			ChampionLevel:              "",
			SummonerSpellImagePaths:    [2]string{},
			RuneKeystoneImagePath:      "",
			RuneSecondaryTreeImagePath: "",
			ItemImagePaths:             [7]string{},
		})
	}
	return mapp
}
