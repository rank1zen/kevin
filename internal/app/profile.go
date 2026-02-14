package app

import (
	"net/http"

	"github.com/rank1zen/kevin/internal/profile"
	"github.com/rank1zen/kevin/internal/profile/web/page/overview_page"
	"github.com/rank1zen/kevin/internal/profile/web/partial/champion"
	"github.com/rank1zen/kevin/internal/profile/web/partial/history_entry"
	"github.com/rank1zen/kevin/internal/profile/web/partial/match_detail"
	"github.com/rank1zen/kevin/internal/profile/web/partial/rank_card"
	"github.com/rank1zen/kevin/internal/profile/web/partial/update"
)

func ProfileRoutes() {
	router := http.NewServeMux()

	profileService := profile.NewProfileService(nil, nil)

	// Pages
	router.Handle("GET /profile/{riotID}/{$}", overview_page.NewHandler(profileService))

	// Partials
	router.Handle("GET /partial/rank_card.RankCard/{$}", rank_card.NewHandler(profileService))
	router.Handle("GET /partial/profile.HistoryEntry", history_entry.NewHandler(profileService))
	router.Handle("GET /partial/profile.ChampionList", champion.NewHandler(*profileService))
	router.Handle("GET /partial/profile.MatchDetailBox", match_detail.NewHandler(*profileService))
	router.Handle("POST /partial/profile.UpdateProfile", update.NewHandler(*profileService))
}
