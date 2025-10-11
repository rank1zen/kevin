package page

import (
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/riot"
)

type HomePageHandler frontend.Handler

func (h *HomePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	v := HomePageData{
		Region: riot.RegionNA1,
	}

	if err := HomePage(ctx, v).Render(ctx, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
