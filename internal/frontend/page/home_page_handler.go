package page

import (
	"fmt"
	"net/http"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/frontend"
)

type HomePageHandler internal.Datasource

func (h *HomePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	region := frontend.StrToRiotRegion(r.FormValue("region"))

	v := HomePageData{
		Region: region,
	}

	c := HomePage(r.Context(), v)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, fmt.Errorf("failed to render home page template: %w", err))
		return
	}
}
