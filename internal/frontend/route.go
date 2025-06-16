package frontend

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

func Routes(ds *internal.Datasource) http.Handler {
	router := http.NewServeMux()

	router.Handle(
		"GET /static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
	)

	router.HandleFunc(
		"GET /",
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			q := r.FormValue("q")
			if q == "" {
				w.WriteHeader(http.StatusOK)
				HomePage().Render(ctx, w)
				return
			}
			results, err := ds.SearchSummoner(ctx, q)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			title := fmt.Sprintf("%s - KEVIN", q)
			w.WriteHeader(http.StatusOK)
			SearchPage(title, q, results).Render(ctx, w)
		},
	)

	router.HandleFunc(
		"GET /summoner/{PUUID}",
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			puuid := r.PathValue("PUUID")

			name, _, err := ds.GetSummoner(ctx, puuid)
			if err != nil {
				if errors.Is(err, internal.ErrSummonerNotFound) {
					// fetch from riot...
					// but what is included in this I don't really care right now
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			rank, err := ds.GetRank(ctx, puuid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			title := fmt.Sprintf("%s - KEVIN", name)
			view := SummonerMatchView(puuid, name, rank)
			w.WriteHeader(http.StatusOK)
			Page(title, view).Render(ctx, w)
		},
	)

	router.HandleFunc(
		"GET /summoner/{PUUID}/matchlist",
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			puuid := r.PathValue("PUUID")
			var page int = 0
			if pageQuery := r.FormValue("page"); pageQuery != "" {
				if pageVal, err := strconv.Atoi(pageQuery); err == nil {
					page = pageVal
				}
			}

			slog.Debug(fmt.Sprintf("fetching page %d", page))
			err := ds.UpdateMatchlist(ctx, riot.PlatformNA1, puuid, page*10,10)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			matches, err := ds.GetMatchlist(ctx, puuid, page)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			SummonerMatchlist(puuid, page, matches).Render(ctx, w)
		},
	)

	router.HandleFunc(
		"GET /summoner/{PUUID}/champions",
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			puuid := r.PathValue("PUUID")

			name, _, err := ds.GetSummoner(ctx, puuid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			rank, err := ds.GetRank(ctx, puuid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			champions, err := ds.GetChampions(ctx, puuid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			title := fmt.Sprintf("%s - KEVIN", name)
			view := SummonerChampionView(puuid, name, rank, champions)
			w.WriteHeader(http.StatusOK)
			Page(title, view).Render(ctx, w)
		},
	)

	router.HandleFunc(
		"GET /match/{matchID}",
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			matchID := internal.MatchID(r.PathValue("matchID"))

			match, participants, err := ds.GetMatch(ctx, matchID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			title := fmt.Sprintf("%s - KEVIN", matchID)
			view := MatchPage(matchID, match, participants)
			w.WriteHeader(http.StatusOK)
			Page(title, view).Render(ctx, w)
		},
	)

	return router
}
