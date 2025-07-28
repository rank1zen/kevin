package frontend

import (
	"context"
	"time"

	"github.com/a-h/templ"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/component"
	"github.com/rank1zen/kevin/internal/riot"
)

type View interface {
	ToTempl(context.Context) templ.Component
}

type LiveMatch struct {
	Date time.Time

	internal.LiveMatch
}

func (m LiveMatch) ToTempl(ctx context.Context) templ.Component {
	c := component.LiveMatchModalLayout{
		AverageRank: &internal.RankDetail{},
		StartTime:   time.Time{},
		BlueSide:    [5]component.LiveMatchRowLayout{},
		RedSide:     [5]component.LiveMatchRowLayout{},
	}

	// TODO: currently cannot compute
	c.AverageRank = nil

	c.StartTime = m.Date

	for i, p := range m.GetTeamParticipants(100) {
		card := component.LiveMatchRowLayout{
			MatchID: p.MatchID,
			ChampionWidget: component.ChampionWidget{
				ChampionSprite: component.ChampionSprite{
					ChampionID: p.ChampionID,
					Size:       component.TextSize2XL,
				},
				SummonerD: &component.SummonerSprite{
					SummonerID: p.SummonersIDs[0],
				},
				SummonerF: &component.SummonerSprite{
					SummonerID: p.SummonersIDs[1],
				},
			},
			RuneWidget: component.RuneWidget{RunePage: p.Runes},
			TeamID:     p.TeamID,
			PUUID:      riot.PUUID(p.PUUID),
			// Name:       name, TODO: WIP
			// Tag:        tag,
			Rank: nil,
		}

		c.BlueSide[i] = card
	}

	return c
}

type LiveMatchNotFound struct {
}

func (m LiveMatchNotFound) ToTempl(ctx context.Context) templ.Component {
	c := component.NoLiveMatchModalWindow{}
	return c
}

type HomePage struct {}

func (m HomePage) ToTempl(ctx context.Context) templ.Component {
	c := component.HomePage{}
	return c
}

type MatchDetail struct {
}

func (m MatchDetail) ToTempl(ctx context.Context) templ.Component {
	c := component.MatchDetailLayout{}

	c.Date = match.Date

	c.Duration = match.Duration

	for i, p := range match.Participants {
		row := component.MatchDetailPlayerRowLayout{
			PUUID: "",
			Name:  "",
			Tag:   "",
			ChampionWidget: component.ChampionWidget{
				ChampionSprite: component.ChampionSprite{},
				ChampionLevel:  0,
				SummonerD:      &component.SummonerSprite{},
				SummonerF:      &component.SummonerSprite{},
			},
			KDAWidget:       component.KDAWidget{},
			CSWidget:        component.CSWidget{},
			RuneWidget:      component.RuneWidget{},
			ItemWidget:      component.ItemWidget{},
			RankDeltaWidget: component.RankDeltaWidget{},
		}

		if p.TeamID == 100 {
			c.BlueSide[i%5] = row
		} else {
			c.RedSide[i%5] = row
		}
	}

	return c
}

type MatchHistory struct {

}

func (m MatchHistory) ToTempl(ctx context.Context) templ.Component {
	c := component.MatchHistoryList{}

	c.Matches = []component.MatchHistoryRowLayout{}

	for _, m := range matches {
		c.Matches = append(
			c.Matches,
			component.MatchHistoryRowLayout{
				MatchID: m.MatchID,
				ChampionWidget: component.ChampionWidget{
					ChampionSprite: component.ChampionSprite{
						ChampionID: m.ChampionID,
						Size:       component.TextSize2XL,
					},
					ChampionLevel: m.ChampionLevel,
					SummonerD: &component.SummonerSprite{
						SummonerID: m.SummonerIDs[0],
					},
					SummonerF: &component.SummonerSprite{
						SummonerID: m.SummonerIDs[1],
					},
				},
				KDAWidget: component.KDAWidget{
					Kills:          m.Kills,
					Deaths:         m.Deaths,
					Assists:        m.Assists,
					KilLDeathRatio: (float32(m.Kills) + float32(m.Assists)) / float32(m.Deaths),
				},
				CSWidget: component.CSWidget{
					CS:          m.CreepScore,
					CSPerMinute: m.CreepScorePerMinute,
				},
				RuneWidget: component.RuneWidget{
					RunePage: m.Runes,
				},
				ItemWidget: component.ItemWidget{
					Items:       [7]component.Tooltip{},
					VisionScore: 0,
					ItemHistory: component.Popover{},
				},
				RankChange: nil,
				LPChange:   nil,
				Win:        m.Win,
			},
		)
	}

	return c
}

type SearchResult struct {

}

func (m SearchResult) ToTempl(ctx context.Context) templ.Component {
	c := []component.SearchResultLink{}

	for _, r := range storeSearchResults {
		rank, err := ds.GetStore().GetRank(ctx, r.PUUID, time.Now(), true)
		if err != nil {
			return nil, fmt.Errorf("getting rank for %s#%s: %w", r.Name, r.Tagline, err)
		}

		row := SearchResultLink{
			Region: region,
			PUUID:  r.PUUID,
			Name:   r.Name,
			Tag:    r.Tagline,
			Rank:   &rank.Detail.Rank,
		}

		searchResults = append(searchResults, row)
	}
	return c
}

type SearchResultNotFound struct {

}

func (m SearchResultNotFound) ToTempl(ctx context.Context) templ.Component {
	c := []component.SearchResultLink{}
	return c
}

type SummonerPage struct {
	// Region specifies a riot region. All results in the page belong to
	// this region.
	Region riot.Region

	PUUID riot.PUUID

	Name, Tag string

	LastUpdated time.Time

	Rank *internal.Rank
}

func (m SummonerPage) ToTempl(ctx context.Context) component.SummonerPage {
	c := component.SummonerPage{
		Region: m.Region,

		PUUID: m.PUUID,

		Name: m.Name,

		Tag: m.Tag,

		LastUpdated: m.LastUpdated,

		RankWidget: component.RankWidget{
			Rank: m.Rank,
			Size: component.TextSizeSM,
		},

		LiveMatch: component.Modal{
			ButtonChildren: component.ButtonLayout{Icon: component.ViewListIcon},
			PanelChildren: component.LiveMatchModalWindowLoader{
				Request: ZGetSummonerChampionsRequest{},
			},
		},

		Champions: component.Modal{
			ButtonChildren: component.ButtonLayout{Icon: component.ViewListIcon},
			PanelChildren:  nil,
		},

		MatchHistory: []component.Section{},
	}

	for i := range 7 {
		day := GetDay(i)

		c.MatchHistory = append(
			c.MatchHistory,
			component.Section{
				Heading: day.Format("Monday, Jan 2"),
				Content: component.MatchHistoryListLoader{
					Request: MatchHistoryRequest{
						Region: m.Region,
						PUUID:  m.PUUID,
						Date:   day,
					},
				},
			},
		)
	}

	return c
}

type SummonerPageNotFound struct {
}

func (m SummonerPageNotFound) ToTempl(ctx context.Context) templ.Component {
	return nil
}

type SummonerChampion component.ChampionModalLayout

func (m SummonerChampion) ToTempl(ctx context.Context) templ.Component {

	layout := ChampionModalLayout{
		List: ChampionModalList{
			Champions: []ChampionModalRowLayout{},
		},
	}

	for _, c := range storeChampions {
		layout.List.Champions = append(
			layout.List.Champions,
			ChampionModalRowLayout{
				ChampionWidget: ChampionWidget{ChampionSprite: ChampionSprite{ChampionID: int(c.Champion), Size: TextSize2XL}},
				GamesPlayed:    c.GamesPlayed,
				Wins:           c.Wins,
				Losses:         c.Losses,
				WinRate:        ComputeFraction(c.Wins, c.GamesPlayed),
				KDAWidget: KDAWidget{
					Kills:          int(c.AverageKillsPerGame),
					Deaths:         int(c.AverageDeathsPerGame),
					Assists:        int(c.AverageAssistsPerGame),
				},
				CSWidget: CSWidget{CS: int(c.AverageCreepScorePerGame), CSPerMinute: c.AverageCreepScorePerMinutePerGame},
			},
		)
	}
	return nil
}
