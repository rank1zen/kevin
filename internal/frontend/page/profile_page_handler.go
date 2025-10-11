package page

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/page"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/view/profile"
)

type ProfilePageHandler frontend.Handler

func (h *ProfilePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := frontend.LoggerFromContext(ctx)

	region := riot.RegionNA1

	payload := slog.Group("payload", "region", region)

	riotID := r.PathValue("riotID")
	name, tag, err := parseRiotID(riotID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Debug("failed to resolve riot id", "err", err, payload)
		return
	}

	storeProfile, err := h.Datasource.GetProfileDetailByRiotID(ctx, region, name, tag)

	data := page.ProfilePageData{
		PUUID:          storeProfile.PUUID,
		Region:         region,
		Name:           name,
		Tag:            tag,
		HistoryEntryCh: make(chan profile.HistoryEntryData),
		RankCardCh:     make(chan profile.RankCardData),
		ChampionListCh: make(chan profile.ChampionListData),
	}

	c := page.ProfilePage(ctx, data)

	go func() {
		defer close(data.HistoryEntryCh)
		defer close(data.ChampionListCh)
		defer close(data.RankCardCh)

		days := GetDays(time.Now())

		for i := range len(days) - 1 {
			historyEntryData, err := h.GetMatchHistory(ctx, MatchHistoryRequest{
				Region:  region,
				PUUID:   data.PUUID,
				StartTS: days[i+1],
				EndTS:   days[i],
			})

			if err == nil {
				data.HistoryEntryCh <- *historyEntryData
			}
		}

		championListData, err := h.GetSummonerChampions(ctx, GetSummonerChampionsRequest{
			Region: region,
			PUUID:  data.PUUID,
			Week:   GetCurrentWeek(),
		})

		if err == nil {
			data.ChampionListCh <- *championListData
		}
	}()

	templ.Handler(c, templ.WithStreaming()).ServeHTTP(w, r)
}

func parseRiotID(riotID string) (name, tag string, err error) {
	index := strings.Index(riotID, "-")
	if index == -1 {
		return "", "", ErrInvalidRiotID
	}

	if index == len(riotID)-1 {
		return "", "", ErrInvalidRiotID
	}

	name = riotID[:index]
	tag = riotID[index+1:]

	if index := strings.Index(tag, "-"); index != -1 {
		return "", "", ErrInvalidRiotID
	}

	return name, tag, nil
}
