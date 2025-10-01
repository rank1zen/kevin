package profile

import (
	"net/http"
)

func MatchDetailBoxPartial() http.Handler {
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

	return http.HandlerFunc(handler)
}
