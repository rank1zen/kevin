package page

import (
	"fmt"
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/frontend/view/ladder"
	"github.com/rank1zen/kevin/internal/service"
)

type LadderPageHandler service.Service

func (h *LadderPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := service.GetLeaderboardRequest{}

	req.Region = r.FormValue("region")

	result, err := (*service.LeaderboardService)(h).GetLeaderboard(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, fmt.Errorf("failed to get leaderboard for ladder page: %w", err))
		return
	}

	data := LadderPageData{
		Region:  result.Region,
		Entries: []ladder.RowData{},
	}

	for i, entry := range result.Entries {
		data.Entries = append(data.Entries, ladder.RowData{
			Place:         i,
			Name:          entry.Name,
			Tag:           entry.Tag,
			Wins:          entry.Wins,
			Losses:        entry.Losses,
			WinPercentage: 50,
			Tier:          entry.Tier,
			Division:      entry.Division,
			LP:            entry.LP,
		})
	}

	c := LadderPage(r.Context(), data)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, fmt.Errorf("failed to render page template: %w", err))
		return
	}
}
