// Package routes defines ALL pages, partials, and misc things for the application.
package route

import (
	"net/http"

	"github.com/rank1zen/kevin/internal/profile"
	"github.com/rank1zen/kevin/internal/profile/web/page/overview_page"
	"github.com/rank1zen/kevin/internal/profile/web/partial/champion"
	"github.com/rank1zen/kevin/internal/profile/web/partial/history_entry"
	"github.com/rank1zen/kevin/internal/profile/web/partial/match_detail"
	"github.com/rank1zen/kevin/internal/profile/web/partial/rank_card"
	"github.com/rank1zen/kevin/internal/profile/web/partial/update"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/web/page/not_found_page"
)

func Router(
	riotClient *riot.Client,
	profileService *profile.ProfileService,
) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /ready/{$}", func(w http.ResponseWriter, r *http.Request) {
		if _, err := riotClient.Account.GetAccountByRiotID(r.Context(), riot.RegionNA1, "orrange", "NA1"); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	profileRoutes(router, profileService)

	router.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	router.Handle("/", not_found_page.NewHandler())

	return router
}

func profileRoutes(router *http.ServeMux, profileService *profile.ProfileService) {
	// Pages
	router.Handle("GET /profile/{riotID}/{$}", overview_page.NewHandler(profileService))

	// Partials
	router.Handle("GET /partial/rank_card.RankCard/{$}", rank_card.NewHandler(profileService))
	router.Handle("GET /partial/profile.HistoryEntry", history_entry.NewHandler(profileService))
	router.Handle("GET /partial/profile.ChampionList", champion.NewHandler(*profileService))
	router.Handle("GET /partial/profile.MatchDetailBox", match_detail.NewHandler(*profileService))
	router.Handle("POST /partial/profile.UpdateProfile", update.NewHandler(*profileService))
}
