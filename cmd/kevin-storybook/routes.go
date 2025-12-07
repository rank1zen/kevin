package main

import (
	"context"

	"github.com/a-h/templ"
	"github.com/rank1zen/kevin/internal/frontend/page"
	"github.com/rank1zen/kevin/internal/frontend/view/profile"
	"github.com/rank1zen/kevin/internal/frontend/view/search"
	"github.com/rank1zen/kevin/internal/frontend/view/shared"
)

type Route struct {
	Name     string
	Renderer templ.Component
}

var routes = []*Route{
	{
		Name:     "profile.RankCard",
		Renderer: profile.RankCard(context.Background(), rankCardData),
	},
	{
		Name:     "profile.RankCard (Unranked)",
		Renderer: profile.RankCard(context.Background(), profile.RankCardData{RankDetail: nil}),
	},
	{
		Name:     "profile.HistoryCard",
		Renderer: profile.HistoryCard(context.Background(), historyCardData),
	},
	{
		Name:     "ChampionList",
		Renderer: profile.ChampionList(context.Background(), championListData),
	},
	{
		Name:     "profile.MatchDetailBox",
		Renderer: profile.MatchDetailBox(context.Background(), matchDetailBoxData),
	},
	{
		Name:     "profile.PartialHistoryEntry",
		Renderer: profile.PartialHistoryEntry(context.Background(), partialHistoryEntryData),
	},
	{
		Name:     "profile.LiveMatchScoreboard",
		Renderer: profile.LiveMatchScoreboard(context.Background(), liveMatchScoreboardData),
	},
	{
		Name:     "page.ProfilePage",
		Renderer: page.ProfilePage(context.Background(), profilePageData),
	},
	{
		Name:     "page.ProfileLiveMatchPage",
		Renderer: page.ProfileLiveMatchPage(context.Background(), profileLiveMatchPageData),
	},
	{
		Name:     "shared.Header",
		Renderer: shared.Header(context.Background(), sharedHeaderData),
	},
	{
		Name:     "search.ResultCard",
		Renderer: search.ResultCard(context.Background(), searchResultCardData),
	},
	{
		Name:     "search.ResultMenu",
		Renderer: search.ResultMenu(context.Background(), searchResultMenuData),
	},
	{
		Name:     "page.LivematchPage",
		Renderer: page.LivematchPage(context.Background(), livematchPageData),
	},
}
