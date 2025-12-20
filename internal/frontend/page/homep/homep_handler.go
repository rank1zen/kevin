package homep

import (
	"fmt"
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/service"
)

type HomepHandler service.Service

func (h *HomepHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	region := frontend.StrToRiotRegion(r.FormValue("region"))

	v := HomepData{
		Region: region,
	}

	c := Homep(r.Context(), v)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, fmt.Errorf("failed to render home page template: %w", err))
		return
	}
}
