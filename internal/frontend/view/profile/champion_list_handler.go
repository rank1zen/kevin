package profile

import (
	"errors"
	"net/http"
	"time"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/riot"
)

type ChampionListHandler frontend.Handler

func (h *ChampionListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := frontend.GetSummonerChampionsRequest{}

	req.Region = new(riot.Region)
	*req.Region = frontend.StrToRiotRegion(r.FormValue("region"))

	req.PUUID = riot.PUUID(r.FormValue("puuid"))

	if start, err := time.Parse(time.RFC3339, r.FormValue("start")); err == nil {
		req.StartTS = &start
	}

	if end, err := time.Parse(time.RFC3339, r.FormValue("end")); err == nil {
		req.EndTS = &end
	}

	profileService := (*frontend.ProfileService)(h)

	storeChamps, err := profileService.GetSummonerChampions(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, errors.New("storage failure"))
		return
	}

	v := &ChampionListData{
		Champions: []ChampionItemData{},
	}

	for _, champ := range storeChamps {
		v.Champions = append(v.Champions, ChampionItemData{
			Champion:    int(champ.Champion),
			GamesPlayed: champ.GamesPlayed,
			Wins:        champ.Wins,
			Losses:      champ.Losses,
			LPDelta:     nil,
		})
	}

	c := ChampionList(r.Context(), *v)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, errors.New("templ error"))
		return
	}
}
