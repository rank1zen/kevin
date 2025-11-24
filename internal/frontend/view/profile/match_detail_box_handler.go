package profile

import (
	"errors"
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
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

	v := MatchDetailBoxData{
		Date:     storeMatch.Date,
		Duration: storeMatch.Duration,
		BlueSide: MatchTeamListData{
			Participants: []MatchParticipantCardData{},
		},
		RedSide: MatchTeamListData{
			Participants: []MatchParticipantCardData{},
		},
	}

	for i := range 5 {
		v.BlueSide.Participants = append(v.BlueSide.Participants, *NewMatchParticipantCardData(storeMatch.Participants[i]))
	}

	for i := range 5 {
		v.RedSide.Participants = append(v.RedSide.Participants, *NewMatchParticipantCardData(storeMatch.Participants[5+i]))
	}

	c := MatchDetailBox(r.Context(), v)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, errors.New("templ error"))
		return
	}
}
