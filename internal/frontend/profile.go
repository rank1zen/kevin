package frontend

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/component"
	"github.com/rank1zen/kevin/internal/component/view"
)

type ProfileHandler struct {
	Datasource *internal.Datasource
}

// GetMatchHistory returns matches played between the given timestamps. The
// method will fetch riot first to ensure all matches played on date are in
// store.
func (h *ProfileHandler) GetMatchHistory(ctx context.Context, req MatchHistoryRequest) (view.HistoryList, error) {
	start := req.StartTS
	end := req.EndTS

	storeMatches, err := h.Datasource.GetMatchHistory(ctx, req.Region, req.PUUID, start, end)
	if err != nil {
		return view.HistoryList{}, fmt.Errorf("storage failure: %w", err)
	}

	v := view.HistoryList{
		MatchHistory: []struct {
			internal.SummonerMatch
			Path string
			Data string
		}{},
	}

	for _, match := range storeMatches {
		path, data := makeGetMatchDetailRequest(req.Region, match.MatchID)

		v.MatchHistory = append(v.MatchHistory, struct {
			internal.SummonerMatch
			Path string
			Data string
		}{
			SummonerMatch: match,
			Path:          path,
			Data:          string(data),
		})
	}

	return v, nil
}

// GetMatchDetail returns [MatchDetail].
func (h *ProfileHandler) GetMatchDetail(ctx context.Context, req MatchDetailRequest) (component.Component, error) {
	matchDetail, err := h.Datasource.GetMatchDetail(ctx, req.Region, req.MatchID)
	if err != nil {
		return nil, err
	}

	mapper := FrontendToMatchDetailMapper{
		MatchDetail: matchDetail,
	}

	c := mapper.Map()

	return c, nil
}

// GetLiveMatch returns a live match view depending on whether the summoner is
// on game.
func (h *ProfileHandler) GetLiveMatch(ctx context.Context, req LiveMatchRequest) (component.Component, error) {
	match, err := h.Datasource.GetLiveMatch(ctx, req.Region, req.PUUID)
	if err != nil {
		if errors.Is(err, internal.ErrNoLiveMatch) {
			c := view.LiveMatchNotFound{}
			return c, nil
		}

		return nil, err
	}

	c := MapLiveMatch(match)

	return c, nil
}

// UpdateSummoner syncs a summoner and their rank with riot.
func (h *ProfileHandler) UpdateSummoner(ctx context.Context, req UpdateSummonerRequest) error {
	if err := h.Datasource.UpdateProfileByRiotID(ctx, req.Region, req.Name, req.Tag); err != nil {
		return err
	}

	return nil
}

// GetSummonerChampions returns The method will fetch all
// games played in the specified interval.
func (h *ProfileHandler) GetSummonerChampions(ctx context.Context, req GetSummonerChampionsRequest) (component.Component, error) {
	start := req.Week
	end := start.Add(7 * 24 * time.Hour)

	storeChampions, err := h.Datasource.GetSummonerChampions(ctx, req.Region, req.PUUID, start, end)
	if err != nil {
		return nil, err
	}

	v := view.ChampionSection{Champions: storeChampions}

	return v, nil
}
