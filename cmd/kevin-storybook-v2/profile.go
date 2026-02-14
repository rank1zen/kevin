package main

import (
	"net/http"
	"time"

	"github.com/rank1zen/kevin/internal/profile/web/page/overview_page"
	"github.com/rank1zen/kevin/internal/profile/web/partial/champion"
	"github.com/rank1zen/kevin/internal/profile/web/partial/history_entry"
	"github.com/rank1zen/kevin/internal/profile/web/partial/match_detail"
	"github.com/rank1zen/kevin/internal/profile/web/partial/rank_card"
)

func registerProfileRoutes(router *http.ServeMux) {
	// Pages
	router.HandleFunc("GET /profile/{riotID}/{$}", overviewHandler)

	// Partials
	router.HandleFunc("GET /profile/web/partial/rank_card/{$}", rankCardHandler)
	router.HandleFunc("GET /profile/web/partial/history_entry/{$}", historyEntryHandler)
	router.HandleFunc("GET /profile/web/partial/champion/{$}", championListHandler)
	router.HandleFunc("GET /profile/web/partial/match_detail/{$}", matchDetailBoxHandler)
	router.HandleFunc("POST /partial/profile.UpdateProfile", updateProfileHandler)
}

func overviewHandler(w http.ResponseWriter, r *http.Request) {
	data := &overview_page.IndexData{
		PUUID:  "1234",
		Region: "NA1",
		Name:   "orrange",
		Tag:    "NA1",
	}

	overview_page.Index(r.Context(), data).Render(r.Context(), w)
}

func rankCardHandler(w http.ResponseWriter, r *http.Request) {
	data := rank_card.RankCardData{
		Region:        "NA1",
		TierDivision:  "Diamond I",
		LP:            75,
		Win:           41,
		Loss:          30,
		WinPercentage: 65,
		Unranked:      false,
	}

	rank_card.RankCard(r.Context(), data).Render(r.Context(), w)
}

func historyEntryHandler(w http.ResponseWriter, r *http.Request) {
	data := &history_entry.IndexData{
		Date: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		Matchlist: []history_entry.CardData{
			{
				Win:               true,
				Date:              "Sat Feb 14 00:09",
				Duration:          "19:30",
				Rank:              "G4",
				LPChange:          "+21",
				ChampionLevel:     "19",
				Kills:             "3",
				Deaths:            "4",
				Assists:           "9",
				KillDeathRatio:    "0.75",
				CS:                "141",
				CSPerMinute:       "3.7",
				VisionScore:       "13",
				ChampionImagePath: "https://static.bigbrain.gg/assets/lol/riot_static/15.24.1/img/champion/Yasuo.webp",
				SummonerSpellImagePaths: [2]string{
					"https://static.bigbrain.gg/assets/lol/riot_static/15.24.1/img/spell/SummonerFlash.webp",
					"https://static.bigbrain.gg/assets/lol/riot_static/15.24.1/img/spell/SummonerDot.webp",
				},
				RuneKeystoneImagePath:      "https://wiki.leagueoflegends.com/en-us/images/Lethal_Tempo_rune.png?264c8",
				RuneSecondaryTreeImagePath: "https://lol.qq.com/act/a20170926preseason/img/runeBuilder/assets/Domination/icon-d.png",
				ItemImagePaths:             [7]string{},
			},
		},
		NextEntryLoaderData: &history_entry.LoaderData{},
	}

	history_entry.Index(r.Context(), data).Render(r.Context(), w)
}

func championListHandler(w http.ResponseWriter, r *http.Request) {
	data := champion.ChampionListData{
		Champions: []champion.ChampionItemData{},
	}

	champion.ChampionList(r.Context(), data).Render(r.Context(), w)
}

func matchDetailBoxHandler(w http.ResponseWriter, r *http.Request) {
	data := match_detail.MatchDetailBoxData{}

	match_detail.MatchDetailBox(r.Context(), data).Render(r.Context(), w)
}

func updateProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}
