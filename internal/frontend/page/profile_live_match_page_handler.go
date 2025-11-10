package page

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/riot"
)

type ProfileLiveMatchPageHandler internal.Datasource

func (h *ProfileLiveMatchPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	data := ProfileLiveMatchPageData{
		PUUID:  storeProfile.PUUID,
		Region: *req.Region, // FIXME: should use a region value returned by service.
		Name:   storeProfile.Name,
		Tag:    storeProfile.Tagline,
	}

	c := ProfileLiveMatchPage(r.Context(), data)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, errors.New("templ render"))
		return
	}
}
