// Package routes defines ALL pages, partials, and misc things for the application.
package route

import (
	"net/http"

	"github.com/rank1zen/kevin/internal/profile"
	"github.com/rank1zen/kevin/internal/riot"
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

	return router
}

func profileRoutes(router *http.ServeMux, profileService *profile.ProfileService) {
}
