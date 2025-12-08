package page

import (
	"fmt"
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/frontend/view/livematch/teamsection"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/service"
)

type LivematchPageHandler service.Service

func (h *LivematchPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := service.GetLiveMatchByIDRequest{}

	req.Region = new(riot.Region)
	*req.Region = frontend.StrToRiotRegion(r.FormValue("region"))

	req.MatchID = r.PathValue("matchID")

	result, err := (*service.LiveMatchService)(h).GetLiveMatchByID(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, fmt.Errorf("failed to get profile for live match page: %w", err))
		return
	}

	data := LivematchPageData{
		Region:  result.Region,
		MatchID: result.ID,
		Teams:   []teamsection.TeamsectionData{},
	}

	data.Teams = append(data.Teams, teamsection.TeamsectionData{
		TeamName:     "Blue Side",
		Participants: []teamsection.ParticipantCardData{},
	})

	data.Teams = append(data.Teams, teamsection.TeamsectionData{
		TeamName:     "Red Side",
		Participants: []teamsection.ParticipantCardData{},
	})

	for _, participant := range result.Participants {
		cardData := teamsection.ParticipantCardData{
			MatchID:           participant.MatchID,
			Name:              "TODO",
			Tag:               "NOT IMPLEMENTED",
			PUUID:             participant.PUUID,
			ChampionID:        participant.ChampionID,
			SummonerIDs:       participant.SummonersIDs,
			RunePage:          participant.Runes,
			Rank:              nil,
			BannedChampionID:  0,
			WinsInLast20Games: 0,
		}
		if participant.TeamID == 100 {
			data.Teams[0].Participants = append(data.Teams[0].Participants, cardData)
		} else {
			data.Teams[1].Participants = append(data.Teams[1].Participants, cardData)
		}
	}

	c := LivematchPage(r.Context(), data)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, fmt.Errorf("failed to render live match page template: %w", err))
		return
	}
}
