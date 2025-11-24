package profile

import (
	"errors"
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/service"
)

type UpdateProfileHandler service.Service

func (h *UpdateProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := service.UpdateProfileRequest{}

	switch r.Header.Get("Content-type") {
	default:
		if region := r.FormValue("region"); region != "" {
			region := frontend.StrToRiotRegion(region)
			req.Region = &region
		}

		req.Name = r.FormValue("name")
		req.Tag = r.FormValue("tag")
	}

	if err := (*service.ProfileService)(h).UpdateProfile(r.Context(), req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, errors.New("service failure"))
		return
	}

	w.WriteHeader(http.StatusOK)
}
