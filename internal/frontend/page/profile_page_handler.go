package page

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/riot"
)

type ProfilePageHandler internal.Datasource

func (h *ProfilePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := internal.GetProfileRequest{}

	req.Region = new(riot.Region)
	*req.Region = frontend.StrToRiotRegion(r.FormValue("region"))

	req.Name, req.Tag = frontend.ParseRiotID(r.PathValue("riotID"))

	storeProfile, err := (*internal.ProfileService)(h).GetProfile(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, fmt.Errorf("storage error: %w", err))
		return
	}

	data := ProfilePageData{
		PUUID:  storeProfile.PUUID,
		Region: *req.Region,
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
