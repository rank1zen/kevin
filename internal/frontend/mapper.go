package frontend

import (
	"encoding/json"
	"fmt"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/component"
	"github.com/rank1zen/kevin/internal/component/live"
	"github.com/rank1zen/kevin/internal/component/match"
	"github.com/rank1zen/kevin/internal/component/profile"
	"github.com/rank1zen/kevin/internal/component/search"
	"github.com/rank1zen/kevin/internal/component/shared"
	"github.com/rank1zen/kevin/internal/component/summoner"
	"github.com/rank1zen/kevin/internal/ddragon"
	"github.com/rank1zen/kevin/internal/riot"
)

type FrontendToZZZMapper struct {
	Region riot.Region

	MatchHistory []internal.SummonerMatch
}

func (mapper FrontendToZZZMapper) Map() match.History {
	ma := mapper.MatchHistory

	if len(ma) == 0 {
		history := match.History{
			Style: component.ListStyleRaised,
			Items: []component.Component{component.ComponentFunc(match.NonePlayed)},
		}

		return history
	}

	to := []component.Component{}

	for _, m := range ma {
		card := match.HistoryCard{
			ChampionWidget:   shared.NewMatchChampionWidget(m.ChampionID, m.ChampionLevel, m.SummonerIDs),
			KDAWidget:        shared.NewKDAWidget(m.Kills, m.Deaths, m.Assists),
			CSWidget:         shared.NewCSWidget(m.CreepScore, m.CreepScorePerMinute),
			RuneWidget:       shared.NewRuneWidget(m.Runes),
			ItemWidget:       shared.NewItemInventory(m.Items, m.VisionScore),
			RankChange:       shared.NewRankDeltaWidget(nil, nil, m.Win),
			OpenDetailButton: component.AccordionTrigger{},
		}

		path, data := makeGetMatchDetailRequest(mapper.Region, m.MatchID)

		accordion := component.Accordion{
			Children: card,
			ExtraChildren: component.Loader{
				Path:     path,
				Data:     string(data),
				Children: component.ComponentFunc(match.DetailSkeleton),
				Type:     component.LoaderTypeOnReveal,
			},
		}

		to = append(to, accordion)
	}

	list := match.History{
		Style: component.ListStyleRaised,
		Items: to,
	}

	return list
}

type FrontendToProfilePageMapper struct {
	Region riot.Region

	ProfileDetail internal.ProfileDetail
}

func (mapper FrontendToProfilePageMapper) Map() component.Page {
	pd := mapper.ProfileDetail

	ma := profile.Matchlist{}

	for i := range 7 {
		path, data := makeGetMatchHistoryRequest(mapper.Region, pd.PUUID, i)

		list := component.List{
			Style: component.ListStyleRaised,
			Items: []component.Component{},
		}

		list.Items = append(list.Items, component.ComponentFunc(match.HistorySkeleton))

		loader := component.Loader{
			Path:     path,
			Data:     string(data),
			Children: list,
		}

		ma = append(ma, component.Section{
			Heading: GetDay(i).Format("Monday, Jan 2"),
			Content: loader,
		})
	}

	champPath, champData := makeGetChampionList(mapper.Region, pd.PUUID)
	livePath, liveData := makeGetLiveMatch(mapper.Region, pd.PUUID)

	champLoader := component.Loader{
		Path:     champPath,
		Type:     component.LoaderTypeOnReveal,
		Data:     string(champData),
		Children: component.ComponentFunc(summoner.ChampstatSkeleton),
	}

	layout := profile.Layout{
		Bar: profile.Bar{
			Name: pd.Name,
			Tag:  pd.Tagline,
			Rank: shared.NewRankWidget(nil),
			Champions: component.Modal{
				ButtonChildren: component.Button{
					Icon: component.ViewListIcon,
				},
				PanelChildren: component.Loader{
					Path:     champPath,
					Data:     string(champData),
					Children: shared.NewLoadingModal(),
				},
			},
			LiveMatch: component.Modal{
				ButtonChildren: component.Button{
					Icon: component.OpenMenuIcon,
				},
				PanelChildren: component.Loader{
					Path:     livePath,
					Data:     string(liveData),
					Children: shared.NewLoadingModal(),
				},
			},
		},
		Matchlist: ma,
		Side: profile.Side{
			{
				Heading: "Live Match",
				Content: live.Me{},
			},
			{
				Heading: "Past Week",
				Content: champLoader,
			},
		},
	}

	if rank := pd.Rank.Detail; rank != nil {
		layout.Bar.Rank = shared.NewRankWidget(&rank.Rank)
	}

	page := component.Page{
		Title:          fmt.Sprintf("%s#%s - Kevin", pd.Name, pd.Tagline),
		HeaderChildren: shared.DefaultPageHeader(),
		Children:       layout,
		FooterChildren: shared.NewPageFooter(),
	}

	return page
}

type FrontendToSearchResultMapper struct {
	Region riot.Region

	Results []internal.SearchResult
}

func (mapper FrontendToSearchResultMapper) Map() component.Component {
	if len(mapper.Results) == 0 {
		return nil
	}

	c := component.List{
		Style: component.ListStyleFlat,
		Items: []component.Component{},
	}

	for _, r := range mapper.Results {
		card := search.ResultCard{
			Name: r.Name,
			Tag:  r.Tagline,
			Rank: shared.RankWidget{
				Rank:         nil,
				ShowTierName: true,
				Size:         component.TextSizeXS,
			},
		}

		if rank := r.Rank.Detail; rank != nil {
			card.Rank.Rank = &rank.Rank
		}

		link := component.Link{
			Href:     fmt.Sprintf("/%s-%s", r.Name, r.Tagline),
			Children: card,
		}

		c.Items = append(c.Items, link)
	}

	return c
}

type FrontendToLoaderMapper struct{}

type FrontendToMatchDetailMapper struct {
	MatchDetail internal.MatchDetail
}

func (mapper FrontendToMatchDetailMapper) Map() component.Component {
	md := mapper.MatchDetail

	var (
		blue = component.List{Style: component.ListStyleRaised, Items: []component.Component{}}
		red  = component.List{Style: component.ListStyleRaised, Items: []component.Component{}}
	)

	newParticipantCard := func(p internal.ParticipantDetail) match.ParticipantCard {
		var rank *internal.Rank = nil
		if p.CurrentRank != nil {
			if p.CurrentRank.Detail != nil {
				rank = &p.CurrentRank.Detail.Rank
			}
		}

		return match.ParticipantCard{
			ChampionWidget: shared.NewMatchChampionWidget(p.ChampionID, p.ChampionLevel, p.SummonerIDs),
			NameWidget:     shared.NewNameWidget("Name", "Tag", rank),
			KDAWidget:      shared.NewKDAWidget(p.Kills, p.Deaths, p.Assists),
			CSWidget:       shared.NewCSWidget(p.CreepScore, p.CreepScorePerMinute),
			RuneWidget:     shared.NewRuneWidget(p.Runes),
			Items:          shared.NewItemInventory(p.Items, p.VisionScore),
		}
	}

	for _, p := range getTeamParticipants(md, 100) {
		blue.Items = append(blue.Items, newParticipantCard(p))
	}

	for _, p := range getTeamParticipants(md, 200) {
		red.Items = append(red.Items, newParticipantCard(p))
	}

	scoreboard := match.Scoreboard{
		BlueSide: component.Section{
			Heading: "Blue Side",
			Content: blue,
		},
		RedSide: component.Section{
			Heading: "Red Side",
			Content: red,
		},
	}

	detail := match.Detail{
		MatchDateDuration: shared.NewMatchDateDurationWidget(md.Date, md.Duration),
		Tabs: component.TabList{
			Tabs: []component.TabTrigger{
				{Label: "Scoreboard"},
				{Label: "Performance"},
			},
		},
		Panels: component.TabPanelList{
			PanelList: []component.TabPanel{
				{Children: scoreboard},
				{Children: nil},
			},
		},
	}

	c := component.TabContainer{
		Children: detail,
	}

	return c
}

type FrontendToSummonerChampstatMapper struct {
	Champions []internal.SummonerChampion
}

func (mapper FrontendToSummonerChampstatMapper) Map() component.Component {
	ch := mapper.Champions

	list := component.List{
		Style: component.ListStyleFlat,
		Items: []component.Component{},
	}

	for i, c := range ch {
		if i > 3 {
			continue
		}

		name := ddragon.ChampionMap[int(c.Champion)].Name

		card := summoner.ChampionCard{
			ChampionWidget: shared.Champion{
				ChampionSprite: shared.NewChampionSprite(int(c.Champion), component.TextSizeLG),
			},
			ChampionName: name,
			Wins:         c.Wins,
			Losses:       c.Losses,
			LPDelta:      nil,
		}

		list.Items = append(list.Items, card)
	}

	return list
}

func NewSearchNotFoundCard(region riot.Region, name, tag string) search.NotFoundCard {
	path, data := makeUpdateSummoner(region, name, tag)

	c := search.NotFoundCard{
		Region: region,
		Name:   name,
		Tag:    tag,
		Path:   path,
		Data:   string(data),
	}

	return c
}

func makeGetMatchHistoryRequest(region riot.Region, puuid riot.PUUID, index int) (string, []byte) {
	path := "/summoner/matchlist"

	req := MatchHistoryRequest{
		Region: region,
		PUUID:  puuid,
		Date:   GetDay(index),
	}

	bytes, _ := json.Marshal(req)

	return path, bytes
}

func makeGetMatchDetailRequest(region riot.Region, matchID string) (string, []byte) {
	path := "/match"

	req := MatchDetailRequest{
		Region:  region,
		MatchID: matchID,
	}

	bytes, _ := json.Marshal(req)

	return path, bytes
}

func makeGetLiveMatch(region riot.Region, puuid riot.PUUID) (string, []byte) {
	path := "/summoner/live"

	req := LiveMatchRequest{
		Region: region,
		PUUID:  puuid,
	}

	bytes, _ := json.Marshal(req)

	return path, bytes
}

func makeGetChampionList(region riot.Region, puuid riot.PUUID) (string, []byte) {
	path := "/summoner/champions"

	req := GetSummonerChampionsRequest{
		Region: region,
		PUUID:  puuid,
		Week:   GetCurrentWeek(),
	}

	bytes, _ := json.Marshal(req)

	return path, bytes
}

func makeUpdateSummoner(region riot.Region, name, tag string) (string, []byte) {
	path := "/summoner/fetch"

	req := UpdateSummonerRequest{
		Region: region,
		Name:   name,
		Tag:    tag,
	}

	bytes, _ := json.Marshal(req)

	return path, bytes
}

func getTeamParticipants(match internal.MatchDetail, teamID int) [5]internal.ParticipantDetail {
	result := [5]internal.ParticipantDetail{}
	i := 0
	for _, p := range match.Participants {
		if p.TeamID == teamID {
			result[i] = p
			i++
		}
	}

	return result
}
