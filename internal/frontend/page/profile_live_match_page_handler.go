package page

import (
	// "errors" // No longer needed directly
	"fmt"
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/service"
)

type ProfileLiveMatchPageHandler service.Service

func (h *ProfileLiveMatchPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := service.GetProfileRequest{}

	req.Region = new(riot.Region)
	*req.Region = frontend.StrToRiotRegion(r.FormValue("region"))

	req.Name, req.Tag = frontend.ParseRiotID(r.PathValue("riotID"))

	storeProfile, err := (*service.ProfileService)(h).GetProfile(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, fmt.Errorf("failed to get profile for live match page: %w", err))
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
		frontend.LogError(r, fmt.Errorf("failed to render profile live match page template: %w", err))
		return
	}
}
