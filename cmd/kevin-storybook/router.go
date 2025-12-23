package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rank1zen/kevin/internal/ddragon"
	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/frontend/partial/historyentry"
	"github.com/rank1zen/kevin/internal/frontend/partial/rank_card"
	"github.com/rank1zen/kevin/internal/frontend/view/historycard"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	dd := ddragon.New("https://ddragon.leagueoflegends.com/cdn/15.24.1")

	mux.Handle("GET /{$}", &IndexHandler{})

	mux.HandleFunc(
		"GET /partial/rank_card.RankCard/{$}",
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)

			c := rank_card.RankCard(r.Context(), rank_card.RankCardData{
				Region:        "NA1",
				TierDivision:  "Diamond I",
				LP:            68,
				Win:           12,
				Loss:          4,
				WinPercentage: 23,
				Unranked:      false,
			})

			if err := c.Render(r.Context(), w); err != nil {
				w.WriteHeader(500)
			}
		},
	)

	mux.HandleFunc(
		"GET /partial/historyentry.Historyentry/{$}",
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)

			lpChange := 10

			c := historyentry.Historyentry(r.Context(), historyentry.HistoryentryData{
				Date: time.Now(),
				Matchlist: []historycard.HistorycardData{
					{
						MatchID:                "NA1_123456789",
						ChampionLevel:          17,
						ChampionIconPath:       dd.GetChampionImage(41),
						Kills:                  3,
						Deaths:                 4,
						Assists:                5,
						KillDeathRatio:         6.1,
						CS:                     324,
						CSPerMinute:            8.9,
						RunePage:               sampleRunePage,
						Items:                  sampleItems,
						VisionScore:            23,
						RankChange:             nil,
						LPChange:               &lpChange,
						Win:                    true,
						SummonerSpellIconPaths: [2]string{},
					},
				},
			})

			if err := c.Render(r.Context(), w); err != nil {
				w.WriteHeader(500)
			}
		},
	)

	mux.Handle("GET /static/", http.FileServer(http.FS(frontend.StaticAssets)))

	for _, route := range routes {
		mux.HandleFunc(
			fmt.Sprintf("GET /%s/{$}", route.Name),
			func(w http.ResponseWriter, r *http.Request) {
				c := Fragment(route.Renderer)
				if err := c.Render(r.Context(), w); err != nil {
					w.WriteHeader(500)
				}
			},
		)
	}

	return mux
}
