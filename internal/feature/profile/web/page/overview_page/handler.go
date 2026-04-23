package overview_page

import (
	"fmt"
	"net/http"

	"github.com/rank1zen/kevin/internal/feature/profile"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/web"
	"github.com/rank1zen/kevin/internal/web/page/server_error_page"
)

type Handler struct {
	service *profile.ProfileService
}

func NewHandler(service *profile.ProfileService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := profile.GetProfileRequest{}

	req.Region = r.FormValue("region")

	req.Name, req.Tag = web.ParseRiotID(r.PathValue("riotID"))

	storeProfile, err := h.service.GetProfile(r.Context(), req)
	if err != nil {
		server_error_page.Render(w, r, err, "failed to get profile for profile page")
		return
	}

	data := &IndexData{
		PUUID:  riot.PUUID(storeProfile.PUUID),
		Region: riot.Region(req.Region),
		Name:   storeProfile.Name,
		Tag:    storeProfile.Tagline,
	}

	c := Index(r.Context(), data)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		web.LogError(r, fmt.Errorf("failed to render profile page template: %w", err))
		return
	}
}
