package page

import (
	"log/slog"
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/riot"
)

type ProfilePageHandler frontend.Handler

func (h *ProfilePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := frontend.LoggerFromContext(ctx)

	region := riot.RegionNA1

	payload := slog.Group("payload", "region", region)

	riotID := r.PathValue("riotID")
	name, tag, err := frontend.ParseRiotID(riotID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Debug("failed to resolve riot id", "err", err, payload)
		return
	}

	storeProfile, err := h.Datasource.GetProfileDetailByRiotID(ctx, region, name, tag)

	data := ProfilePageData{
		PUUID:          storeProfile.PUUID,
		Region:         region,
		Name:           name,
		Tag:            tag,
		HistoryEntryCh: nil,
		RankCardCh:     nil,
		ChampionListCh: nil,
	}

	c := ProfilePage(ctx, data)

	if err := c.Render(ctx, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
