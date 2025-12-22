package rank_card

import (
	"fmt"
	"net/http"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/service"
)

// TODO: not implemented fully
type RankCardHandler service.Service

func (h *RankCardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// (*service.ProfileService)(h).GetRankHistory(r.Context(), req)

	v := RankCardData{
		Region:        "NA1",
		LP:            0,
		Win:           0,
		Loss:          4,
		Unranked:      false,
		TierDivision:  "Diamond I",
		WinPercentage: 12,
	}

	c := RankCard(r.Context(), v)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, fmt.Errorf("failed to render template: %w", err))
		return
	}
}
