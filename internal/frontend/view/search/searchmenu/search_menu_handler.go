package searchmenu

import (
	"fmt"
	"net/http"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/frontend/view/search"
	"github.com/rank1zen/kevin/internal/service"
)

// SearchMenuHandler returns a SearchMenu component.
//
// NOTE: currently only supports form value for query.
type SearchMenuHandler service.Service

func (h *SearchMenuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := service.SearchProfileRequest{}

	region := frontend.StrToRiotRegion(r.FormValue("region"))
	req.Region = &region

	req.Query = r.FormValue("q")

	results, err := (*service.SearchService)(h).SearchProfile(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, fmt.Errorf("failed to search profile: %w", err))
		return
	}

	v := SearchMenuData{
		Name:           "",
		Tag:            "",
		ProfileResults: []search.ResultCardData{},
	}

	for _, result := range results {
		var rank *internal.Rank
		if r := result.Rank; r != nil {
			if rr := r.Detail; rr != nil {
				rank = &internal.Rank{
					Tier:     rr.Rank.Tier,
					Division: rr.Rank.Division,
					LP:       rr.Rank.LP,
				}
			}
		}
		v.ProfileResults = append(v.ProfileResults, search.ResultCardData{
			Name: result.Name,
			Tag:  result.Tagline,
			Rank: rank,
			Path: "",
		})
	}

	c := SearchMenu(r.Context(), v)

	if err := c.Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		frontend.LogError(r, fmt.Errorf("failed to render home page template: %w", err))
		return
	}
}
