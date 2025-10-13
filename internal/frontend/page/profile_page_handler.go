package page

import (
	"errors"
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
)

type ProfilePageHandler frontend.Handler

func (h *ProfilePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	region := frontend.StrToRiotRegion(r.FormValue("region"))
	name, tag := frontend.ParseRiotID(r.PathValue("riotID"))

	storeProfile, err := h.Datasource.GetProfileDetailByRiotID(r.Context(), region, name, tag)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, errors.New("storage failure"))
		return
	}

	data := ProfilePageData{
		PUUID:          storeProfile.PUUID,
		Region:         region,
		Name:           name,
		Tag:            tag,
		HistoryEntryCh: nil,
		RankCardCh:     nil,
		ChampionListCh: nil,
	}

	c := ProfilePage(r.Context(), data)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, errors.New("templ render"))
		return
	}
}
