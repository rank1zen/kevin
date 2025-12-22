package matchdetail

import (
	"errors"
	"net/http"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/frontend/view/matchteam"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/service"
)

type MatchDetailBoxHandler service.Service

func (h *MatchDetailBoxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := service.GetMatchDetailRequest{}

	req.Region = new(riot.Region)
	*req.Region = frontend.StrToRiotRegion(r.FormValue("region"))

	req.MatchID = r.FormValue("matchID")

	storeMatch, err := (*service.MatchService)(h).GetMatchDetail(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, errors.New("storage failure"))
		return
	}

	v := MatchdetailData{
		Date:     storeMatch.Date,
		Duration: storeMatch.Duration,
		Teams:    []matchteam.MatchteamData{},
	}

	v.Teams = append(v.Teams, matchteam.MatchteamData{
		Participants: []matchteam.MatchteamPlayerData{},
	})

	v.Teams = append(v.Teams, matchteam.MatchteamData{
		Participants: []matchteam.MatchteamPlayerData{},
	})

	for _, m := range storeMatch.Participants {
		kda := float32(m.Kills+m.Assists) / float32(m.Deaths)

		var rank *internal.Rank
		if m.RankBefore != nil && m.CurrentRank.Detail != nil {
			rank = &m.CurrentRank.Detail.Rank
		}

		data := matchteam.MatchteamPlayerData{
			MatchID:        m.MatchID,
			Name:           m.Name,
			Tag:            m.Tag,
			PUUID:          m.PUUID,
			ChampionID:     m.ChampionID,
			ChampionLevel:  m.ChampionLevel,
			SummonerIDs:    m.SummonerIDs,
			Kills:          m.Kills,
			Deaths:         m.Deaths,
			Assists:        m.Assists,
			KillDeathRatio: kda,
			CS:             m.CreepScore,
			CSPerMinute:    m.CreepScorePerMinute,
			RunePage:       m.Runes,
			Items:          m.Items,
			VisionScore:    m.VisionScore,
			Rank:           rank,
		}

		if m.TeamID == 100 {
			v.Teams[0].Participants = append(v.Teams[0].Participants, data)
		} else {
			v.Teams[1].Participants = append(v.Teams[1].Participants, data)
		}
	}

	c := Matchdetail(r.Context(), v)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, errors.New("templ error"))
		return
	}
}
