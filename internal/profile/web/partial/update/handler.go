package update

import (
	"net/http"

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
	req := &profile.UpdateProfileRequest{}

	switch r.Header.Get("Content-type") {
	default:
		req.Region = r.FormValue("region")
		req.Name = r.FormValue("name")
		req.Tag = r.FormValue("tag")
	}

	if err := h.service.UpdateProfile(r.Context(), req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
