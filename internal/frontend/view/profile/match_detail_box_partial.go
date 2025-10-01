package profile

import (
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
)

func MatchDetailBoxPartial() *frontend.Partial {
	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		v, err := s.Handler.GetMatchDetail(ctx, req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := MatchDetailBox(ctx, *v).Render(ctx, w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	return nil
}
