package champion

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/profile"
)

type Handler struct {
	service profile.ProfileService
}

func NewHandler(service profile.ProfileService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := profile.GetSummonerChampionsRequest{}

	req.Region = r.FormValue("region")

	req.PUUID = r.FormValue("puuid")

	if start, err := time.Parse(time.RFC3339, r.FormValue("start")); err == nil {
		req.StartTS = &start
	}

	if end, err := time.Parse(time.RFC3339, r.FormValue("end")); err == nil {
		req.EndTS = &end
	}

	storeChamps, err := h.service.GetSummonerChampions(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	v := &IndexData{
		ChampionData: []RowData{},
	}

	for _, champ := range storeChamps {
		v.ChampionData = append(v.ChampionData, RowData{
			ChampionImagePath: "",
			KDA:               fmt.Sprintf("%d/%d/%d", champ.AverageKillsPerGame, champ.AverageDeathsPerGame, champ.AverageAssistsPerGame),
			WinRate:           "",
		})
	}

	c := Index(r.Context(), v)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, errors.New("templ error"))
		return
	}

}
