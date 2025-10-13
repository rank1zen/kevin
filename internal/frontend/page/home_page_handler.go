package page

import (
	"errors"
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
)

type HomePageHandler frontend.Handler

func (h *HomePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	region := frontend.StrToRiotRegion(r.FormValue("region"))

	v := HomePageData{
		Region: region,
	}

	c := HomePage(r.Context(), v)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, errors.New("templ render"))
		return
	}
}
