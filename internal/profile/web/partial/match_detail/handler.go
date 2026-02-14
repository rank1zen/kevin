package match_detail

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
	panic("not yet implemented")
}
