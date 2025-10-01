package frontend

import (
	"net/http"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/page"
)

type IndexService struct {
	Datasource *internal.Datasource

	ProfileService ProfileService
	SearchService  SearchService
}

func (s *IndexService) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /{$}", s.handleHomePage)

	// router.HandleFunc("GET /{riotID}", s.handleProfilePage)
}

func (s *IndexService) handleHomePage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := fromCtx(ctx)

	region := r.FormValue("region")

	// payload := slog.Group("payload", "region", region)

	riotRegion := convertStringToRiotRegion(region)

	v := page.HomePageData{
		Region: riotRegion,
	}

	if err := page.HomePage(ctx, v).Render(ctx, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Debug("failed rendering", "err", err)
		return
	}
}
