package internal

import (
	"log/slog"
	"net/http"
	"strconv"
)

func FetchRoutes(ds *Datasource) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc(
		"GET /",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("this is the fetcher"))
			w.WriteHeader(http.StatusOK)
		},
	)

	router.HandleFunc(
		"GET /summoner",
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			platform := r.FormValue("platform")
			puuid := r.FormValue("puuid")

			err := ds.RecordSummoner(ctx, platform, puuid)
			if err != nil {
				slog.Debug("could not fetch summoner", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		},
	)

	router.HandleFunc(
		"GET /matchlist",
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			platform := r.FormValue("platform")
			puuid := r.FormValue("puuid")
			var page int = 0
			if pageQuery := r.FormValue("page"); pageQuery != "" {
				if pageVal, err := strconv.Atoi(pageQuery); err != nil {
					page = pageVal
				}
			}

			err := ds.UpdateMatchlist(ctx, platform, puuid, page*10, 10)
			if err != nil {
				slog.Debug("could not fetch matches", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		},
	)

	router.HandleFunc(
		"GET /timeline",
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			platform := r.FormValue("platform")
			id := r.FormValue("id")

			err := ds.RecordMatchTimeline(ctx, platform, id)
			if err != nil {
				slog.Debug("could not fetch timeline", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		},
	)

	return router
}
