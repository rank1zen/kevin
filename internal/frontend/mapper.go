package frontend

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/component"
	"github.com/rank1zen/kevin/internal/component/search"
	"github.com/rank1zen/kevin/internal/component/shared"
	"github.com/rank1zen/kevin/internal/component/summoner"
	"github.com/rank1zen/kevin/internal/component/view"
	"github.com/rank1zen/kevin/internal/ddragon"
	"github.com/rank1zen/kevin/internal/riot"
)

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

type FrontendToSummonerChampstatMapper struct {
	Champions []internal.SummonerChampion
}

func (mapper FrontendToSummonerChampstatMapper) Map() component.Component {
	ch := mapper.Champions

	list := summoner.ChampstatList{}

	for _, c := range ch {
		name := ddragon.ChampionMap[int(c.Champion)].Name

		card := summoner.Champstat{
			ChampionIcon: shared.NewChampionSprite(int(c.Champion), component.TextSizeLG),
			ChampionName: name,
			Wins:         c.Wins,
			Losses:       c.Losses,
			LPDelta:      nil,
		}

		list = append(list, card)
	}

	return list
}

func MapLiveMatch(m internal.LiveMatch) view.LiveMatchModal {
	v := view.LiveMatchModal{
		Date:         m.Date,
		Participants: m.Participants,
	}

	return v
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

func makeGetMatchHistoryRequest(region riot.Region, puuid riot.PUUID, start, end time.Time) (string, []byte) {
	path := "/summoner/matchlist"

	req := MatchHistoryRequest{
		Region:  region,
		PUUID:   puuid,
		StartTS: start,
		EndTS:   end,
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
