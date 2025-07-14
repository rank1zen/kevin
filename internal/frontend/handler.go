package frontend

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

// Handler provides the API for server operations.
type Handler struct {
	Datasource *internal.Datasource
}

func NewHandler(datasource *internal.Datasource) *Handler {
	return &Handler{datasource}
}

// UpdateSummoner
func (h *Handler) UpdateSummoner(ctx context.Context, region riot.Region, name, tag string) error {
	puuid, err := h.Datasource.GetPUUID(ctx, name, tag)
	if err != nil {
		return err
	}

	if err := h.Datasource.UpdateSummoner(ctx, region, puuid); err != nil {
		return err
	}

	return nil
}

// GetLiveMatch the live match view if summoner is in a game in the region. If
// no such game is found, return a view indicating such.
func (h *Handler) GetLiveMatch(ctx context.Context, region riot.Region, puuid riot.PUUID) (view templ.Component, err error) {
	match, err := h.Datasource.GetLiveMatch(ctx, region, puuid)
	if err != nil {
		if errors.Is(err, internal.ErrNoLiveMatch) {
			return NoLiveMatchModalWindow{}, err
		}

		return LiveMatchModalWindow{}, err
	}

	redSide := []LiveMatchSummonerCard{}
	blueSide := []LiveMatchSummonerCard{}

	for _, p := range match.Participants {
		card := LiveMatchSummonerCard{
			Champion:  p.ChampionID,
			Summoners: p.SummonersIDs,
			RunePage:  p.Runes,
			Name:      "Doublelift",
			Tag:       "",
			Rank:      &internal.RankDetail{},
			TeamID:    p.TeamID,
		}
		if p.TeamID == 100 {
			redSide = append(redSide, card)
		} else {
			blueSide = append(blueSide, card)
		}
	}

	window := LiveMatchModalWindow{
		AverageRank: &internal.RankDetail{},
		StartTime:   match.Date,
		RedSide:     redSide,
		BlueSide:    blueSide,
	}

	return window, nil
}

// GetHomePage returns [HomePage].
func (h *Handler) GetHomePage(ctx context.Context, region riot.Region) (templ.Component, error) {
	page := HomePage{Region: region}
	return page, nil
}

// GetSummonerPage returns [SummonerPage] if summoner exists in store,
// otherwise, it will complete a update for summoner, then return
// [SummonerPage]. If no summoner with name#tag exists, return
// [NoSummonerPage].
func (h *Handler) GetSummonerPage(ctx context.Context, region riot.Region, name, tag string) (templ.Component, error) {
	puuid, err := h.Datasource.GetPUUID(ctx, name, tag)
	if err != nil {
		if errors.Is(err, internal.ErrSummonerDoesNotExist) {
			return NoSummonerPage{
				Region: region,
				Name:   name,
				Tag:    tag,
			}, nil
		}

		return nil, fmt.Errorf("get puuid: %w", err)
	}

	summoner, err := h.Datasource.GetStore().GetSummoner(ctx, puuid)
	if err != nil {
		if errors.Is(err, internal.ErrSummonerNotFound) {
			fromCtx(ctx).Info("first time visit", "puuid", puuid)

			if err := h.Datasource.UpdateSummoner(ctx, region, puuid); err != nil {
				return nil, fmt.Errorf("getting puuid: %w", err)
			}
		}
	}

	// try again
	summoner, err = h.Datasource.GetStore().GetSummoner(ctx, puuid)
	if err != nil {
		return nil, fmt.Errorf("get summoner: %w", err)
	}

	rank, err := h.Datasource.GetStore().GetRank(ctx, puuid, time.Now(), true)
	if err != nil {
		return nil, fmt.Errorf("get rank: %w", err)
	}

	page := SummonerPage{
		Region:      region,
		PUUID:       puuid,
		Name:        summoner.Name,
		Tag:         summoner.Tagline,
		Rank:        rank.Detail,
		LastUpdated: rank.EffectiveDate,
	}

	return page, nil
}

// GetSummonerChampions returns [SummonerChampionList] consisting of stats for
// the last week of games.
func (h *Handler) ZGetSummonerChampions(ctx context.Context, region riot.Region, puuid riot.PUUID) (templ.Component, error) {
	champions, err := h.Datasource.GetStore().GetChampions(ctx, puuid, time.Now().Add(-72*time.Hour), time.Now())
	if err != nil {
		return nil, err
	}

	list := SummonerChampionList{Champions: []ChampionPopover{}}

	for _, champion := range champions {
		list.Champions = append(
			list.Champions,
			ChampionPopover{
				Champion:             int(champion.Champion),
				KillParticipation:    champion.KillParticipation,
				CS:                   champion.CreepScore,
				CSPerMinute:          champion.CreepScorePerMinute,
				DamageDealt:          champion.DamageDealt,
				DamageTaken:          champion.DamageTaken,
				DamageDeltaEnemy:     champion.DamageDeltaEnemy,
				DamagePercentageTeam: champion.DamagePercentageTeam,
				GoldEarned:           champion.GoldEarned,
				GoldDeltaEnemy:       champion.GoldDeltaEnemy,
				GoldPercentageTeam:   champion.GoldPercentageTeam,
				VisionScore:          champion.VisionScore,
				PinkWardsBought:      champion.PinkWardsBought,
			},
		)
	}

	return list, nil
}

// GetSummonerChampions returns a [ChampionsModal].
//
// Deprecated: Not using this.
func (h *Handler) GetSummonerChampions(ctx context.Context, region riot.Region, puuid riot.PUUID) (templ.Component, error) {
	champions, err := h.Datasource.GetStore().GetChampions(ctx, puuid, time.Now(), time.Now())
	if err != nil {
		return nil, err
	}

	modal := ChampionsModal{
		Cards: []SummonerChampionCard{},
	}

	for _, champion := range champions {
		modal.Cards = append(modal.Cards, SummonerChampionCard{
			Champion:             int(champion.Champion),
			Kills:                champion.Kills,
			Deaths:               champion.Deaths,
			Assists:              champion.Assists,
			KillParticipation:    champion.KillParticipation,
			CS:                   champion.CreepScore,
			CSPerMinute:          champion.CreepScorePerMinute,
			DamageDealt:          champion.DamageDealt,
			DamageTaken:          champion.DamageTaken,
			DamageDeltaEnemy:     champion.DamageDeltaEnemy,
			DamagePercentageTeam: champion.DamagePercentageTeam,
			GoldEarned:           champion.GoldEarned,
			GoldDeltaEnemy:       champion.GoldDeltaEnemy,
			GoldPercentageTeam:   champion.GoldPercentageTeam,
			VisionScore:          champion.VisionScore,
			PinkWardsBought:      champion.PinkWardsBought,
		})
	}

	return modal, nil
}

// GetSummonerMatchHistory returns a [MatchHistoryList], being the matches
// played on date to date + 24 hours. The method will fetch riot first to
// ensure all matches played on date are in store.
func (h *Handler) GetSummonerMatchHistory(ctx context.Context, region riot.Region, puuid riot.PUUID, date int64) (templ.Component, error) {
	ts := time.Unix(date, 0)

	if err := h.Datasource.ZUpdateMatchHistory(ctx, region, puuid, ts, ts.Add(24*time.Hour)); err != nil {
		return nil, err
	}

	storeMatches, err := h.Datasource.GetStore().GetZMatches(ctx, puuid, ts, ts.Add(24*time.Hour))
	if err != nil {
		return nil, fmt.Errorf("storage failure: %w", err)
	}

	list := MatchHistoryList{Matches: []MatchHistoryRow{}}

	for _, m := range storeMatches {
		list.Matches = append(
			list.Matches,
			MatchHistoryRow{
				MatchID:     m.MatchID,
				Champion:    m.ChampionID,
				Summoners:   m.SummonerIDs,
				Kills:       m.Kills,
				Deaths:      m.Deaths,
				Assists:     m.Assists,
				CS:          m.CreepScore,
				CSPerMinute: m.CreepScorePerMinute,
				RunePage:    m.Runes,
				Items:       m.Items,
				Rank:        &internal.RankDetail{},
				LPChange:    new(int),
			},
		)
	}

	return list, nil
}

// GetMatchScoreboard returns the scoreboard of a match.
func (h *Handler) GetMatchScoreboard(ctx context.Context, id string) (scoreboard templ.Component, err error) {
	return nil, nil
}

// GetSearchResults returns a list of [SearchResultCard] for accounts that
// match q. If no results were found, return [SearchNotFoundCard] instead.
//
// q should be of the form name#tag, if q has no tag, region is used as the
// tag.
func (h *Handler) GetSearchResults(ctx context.Context, region riot.Region, q string) (templ.Component, error) {
	storeSearchResults, err := h.Datasource.GetStore().SearchSummoner(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(storeSearchResults) == 0 {
		var name, tag string
		if i := strings.Index(q, "#"); i != -1 {
			name = q[:i]
			if i+1 == len(q) {
				tag = string(region)
			} else {
				tag = q[i+1:]
			}
		} else {
			name = q
			tag = string(region)
		}

		return SearchNotFoundCard{
			Name:     name,
			Tag:      tag,
			Platform: string(region),
		}, nil
	}

	searchResults := []SearchResultLink{}

	for _, r := range storeSearchResults {
		rank, err := h.Datasource.GetStore().GetRank(ctx, r.Puuid, time.Now(), true)
		if err != nil {
			return nil, fmt.Errorf("getting rank for %s#%s: %w", r.Name, r.Tagline, err)
		}

		row := SearchResultLink{
			Region: region,
			PUUID:  r.Puuid,
			Name:   r.Name,
			Tag:    r.Tagline,
			Rank:   rank.Detail,
		}

		searchResults = append(searchResults, row)
	}

	v := []templ.Component{}
	for _, r := range searchResults {
		v = append(v, r)
	}

	return templ.Join(v...), nil
}
