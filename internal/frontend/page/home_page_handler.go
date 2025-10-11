package page

import (
	"log/slog"
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
)

type HomePageHandler frontend.Handler

func (h *HomePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	region := r.FormValue("region")

	payload := slog.Group("payload", "region", region)

	riotRegion := convertStringToRiotRegion(region)

	v := HomePageData{
		Region: riotRegion,
	}

	if err := HomePage(ctx, v).Render(ctx, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
