package profile

import (
	"errors"
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
)

type UpdateProfileHandler frontend.Handler

func (h *UpdateProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	profileService := (*frontend.ProfileService)(h)

	req := frontend.UpdateProfileRequest{}

	switch r.Header.Get("Content-type") {
	default:
		if region := r.FormValue("region"); region != "" {
			region := frontend.StrToRiotRegion(region)
			req.Region = &region
		}

		req.Name = r.FormValue("name")
		req.Tag = r.FormValue("tag")
	}

	if err := profileService.UpdateProfile(r.Context(), req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, errors.New("service failure"))
		return
	}

	w.WriteHeader(http.StatusOK)
}
