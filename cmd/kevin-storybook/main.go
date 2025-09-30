package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/rank1zen/kevin/internal/riot"
	"github.com/rank1zen/kevin/internal/view/profile"
	"github.com/rank1zen/kevin/internal/view/search"
	"github.com/rank1zen/kevin/internal/view/shared"
)

func main() {
	config := Config{}

	ctx := context.Background()

	if err := config.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}

	os.Exit(0)
}

type Config struct{}

func (c *Config) Run(ctx context.Context) error {
	http.Handle("/{$}", http.HandlerFunc(handleIndex))

	// http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(frontend.StaticAssets))))
	http.Handle("GET /static/", http.FileServer(http.FS(frontend.StaticAssets)))

	for _, route := range routes {
		http.Handle(route.Path, route)
	}

	fmt.Println("Listening on http://localhost:3001")

	return http.ListenAndServe(":3001", nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `<html>`)
	// fmt.Fprintln(w, `<head>`)
	// fmt.Fprintln(w, `<meta name="viewport" content="width=device-width, initial-scale=1"/>`)
	// fmt.Fprintln(w, `<link rel="stylesheet" href="static/output.css"/>`)
	// fmt.Fprintln(w, `</head>`)
	fmt.Fprintln(w, `<body>`)
	fmt.Fprintln(w, `<h2>wtf-storybook</h2>`)
	fmt.Fprintln(w, `<ul>`)
	for _, route := range routes {
		fmt.Fprintf(w, `<li><a href="%s">%s</a></li>`+"\n", route.Path, route.Name)
	}
	fmt.Fprintln(w, `</ul>`)
	fmt.Fprintln(w, `</body>`)
	fmt.Fprintln(w, `</html>`)
}

// Route represents a named reference to a renderer.
type Route struct {
	Name     string
	Path     string
	Renderer templ.Component
}

func (route *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := templ.Handler(Fragment(route.Renderer))
	handler.ServeHTTP(w, r)
}

var routes = []*Route{
	{
		Name: "profile.RankCard",
		Path: "/profile-rank-card",
		Renderer: profile.RankCard(context.Background(), profile.RankCardData{
			RankDetail: &internal.RankDetail{
				Wins:   41,
				Losses: 67,
				Rank: internal.Rank{
					Tier:     riot.TierChallenger,
					Division: riot.Division1,
					LP:       41,
				},
			},
		}),
	},
	{
		Name: "profile.RankCard (Unranked)",
		Path: "/profile-rank-card-unranked",
		Renderer: profile.RankCard(context.Background(), profile.RankCardData{
			RankDetail: nil,
		}),
	},
	{
		Name: "History card",
		Path: "/history-card",
		Renderer: profile.HistoryCard(context.Background(), profile.HistoryCardData{
			MatchID:        "",
			ChampionID:     16,
			ChampionLevel:  16,
			SummonerIDs:    [2]int{4, 5},
			Kills:          3,
			Deaths:         4,
			Assists:        5,
			KillDeathRatio: 3.4111,
			CS:             41,
			CSPerMinute:    6.7,
			RunePage:       internal.RunePage{},
			Items:          [7]int{},
			VisionScore:    41,
			RankChange:     nil,
			LPChange:       nil,
			Win:            false,
			Path:           "",
			Data:           "",
		}),
	},
	{
		Name: "ChampionList",
		Path: "/champion-list",
		Renderer: profile.ChampionList(context.Background(), profile.ChampionListData{
			Champions: []profile.ChampionItemData{
				{
					Champion:    41,
					GamesPlayed: 67,
					Wins:        41,
					Losses:      26,
					LPDelta:     aaa(-1),
				},
				{
					Champion:    67,
					GamesPlayed: 67,
					Wins:        41,
					Losses:      26,
					LPDelta:     aaa(23),
				},
			},
		}),
	},
	{
		Name: "shared.Header",
		Path: "/shared-header",
		Renderer: shared.Header(context.Background(), shared.HeaderData{
			Region: riot.RegionNA1,
		}),
	},
	{
		Name: "profile.MatchDetailBox",
		Path: "/profile-match-detail-box",
		Renderer: profile.MatchDetailBox(
			context.Background(),
			profile.MatchDetailBoxData{
				Date:     time.Date(2041, time.June, 7, 6, 7, 6, 0, time.UTC),
				Duration: 41 * 67 * time.Second,
				RedSide: profile.MatchTeamListData{
					Participants: []profile.MatchParticipantCardData{
						{
							MatchID:        "1",
							Name:           "Joe",
							Tag:            "NA1",
							PUUID:          "1",
							ChampionID:     41,
							ChampionLevel:  18,
							SummonerIDs:    [2]int{},
							Kills:          41,
							Deaths:         42,
							Assists:        67,
							KillDeathRatio: 41.67,
							CS:             410,
							CSPerMinute:    4.1,
							RunePage:       internal.RunePage{},
							Items:          [7]int{},
							VisionScore:    23,
							Rank:           nil,
						},
					},
				},
				BlueSide: profile.MatchTeamListData{
					Participants: []profile.MatchParticipantCardData{
						{
							MatchID:        "1",
							Name:           "Joe",
							Tag:            "NA1",
							PUUID:          "1",
							ChampionID:     41,
							ChampionLevel:  18,
							SummonerIDs:    [2]int{},
							Kills:          41,
							Deaths:         42,
							Assists:        67,
							KillDeathRatio: 41.67,
							CS:             410,
							CSPerMinute:    4.1,
							RunePage:       internal.RunePage{},
							Items:          [7]int{},
							VisionScore:    23,
							Rank:           nil,
						},
						{
							MatchID:        "1",
							Name:           "Bartholomew Montgomery",
							Tag:            "NA1",
							PUUID:          "1",
							ChampionID:     42,
							ChampionLevel:  18,
							SummonerIDs:    [2]int{},
							Kills:          42,
							Deaths:         42,
							Assists:        67,
							KillDeathRatio: 42.67,
							CS:             420,
							CSPerMinute:    4.2,
							RunePage:       internal.RunePage{},
							Items:          [7]int{},
							VisionScore:    42,
							Rank:           nil,
						},
					},
				},
			},
		),
	},
	{
		Name: "search.ResultCard",
		Path: "/search-result-card",
		Renderer: search.ResultCard(context.Background(), search.ResultCardData{
			Name: "Bartholomew Montgomery",
			Tag:  "NA1",
			Rank: &internal.Rank{
				Tier:     riot.TierBronze,
				Division: riot.Division4,
				LP:       67,
			},
		}),
	},
}

func aaa(x int) *int {
	return &x
}
