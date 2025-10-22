package page

import (
	"errors"
	"fmt"
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
		frontend.LogError(r, fmt.Errorf("storage error: %w", err))
		return
	}

	data := ProfilePageData{
		PUUID:  storeProfile.PUUID,
		Region: region,
		Name:   storeProfile.Name,
		Tag:    storeProfile.Tagline,
	}

	c := ProfilePage(r.Context(), data)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, errors.New("templ render"))
		return
	}
}
