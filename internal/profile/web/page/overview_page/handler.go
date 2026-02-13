package overview_page

import (
	// "errors" // No longer needed directly
	"fmt"
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/profile"
	"github.com/rank1zen/kevin/internal/riot"
)

type Handler struct {
	service profile.ProfileService
}

func NewHandler(service profile.ProfileService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := profile.GetProfileRequest{}

	req.Region = r.FormValue("region")

	req.Name, req.Tag = frontend.ParseRiotID(r.PathValue("riotID"))

	storeProfile, err := h.service.GetProfile(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, fmt.Errorf("failed to get profile for profile page: %w", err))
		return
	}

	data := ProfilepData{
		PUUID:  riot.PUUID(storeProfile.PUUID),
		Region: riot.Region(req.Region),
		Name:   storeProfile.Name,
		Tag:    storeProfile.Tagline,
	}

	c := Profilep(r.Context(), data)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, fmt.Errorf("failed to render profile page template: %w", err))
		return
	}
}
