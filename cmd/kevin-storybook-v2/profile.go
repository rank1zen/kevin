package main

import (
	"net/http"

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
		PUUID:                "1234",
		Region:               "NA1",
		Name:                 "orrange",
		Tag:                  "NA1",
		ProfileIconImagePath: "https://ddragon.leagueoflegends.com/cdn/13.14.1/img/profileicon/1003.png",
	}

	overview_page.Index(r.Context(), data).Render(r.Context(), w)
}

func rankCardHandler(w http.ResponseWriter, r *http.Request) {
	data := &rank_card.IndexData{
		Region:        "NA1",
		TierDivision:  "Diamond I",
		LP:            75,
		Win:           41,
		Loss:          30,
		WinPercentage: 65,
		Unranked:      false,
		RankValues:    []float64{
			100.0,
			101.0,
			104.0,
			107.0,
			110.0,
			58.0,
			70.0,
			30.0,
			20.0,
			10.0,
		},
		RankDates:     []string{
			"Jan 01",
			"Jan 02",
			"Jan 03",
			"Jan 04",
			"Jan 05",
			"Jan 06",
			"Jan 07",
			"Jan 08",
			"Jan 09",
			"Jan 10",
		},
	}

	rank_card.Index(r.Context(), data).Render(r.Context(), w)
}

func historyEntryHandler(w http.ResponseWriter, r *http.Request) {
	data := &history_entry.IndexData{
		Date: "Fri Feb 14",
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
		MatchDetailLoaderData: []match_detail.LoaderData{
			{},
		},
		NextEntryLoaderData: &history_entry.LoaderData{},
	}

	history_entry.Index(r.Context(), data).Render(r.Context(), w)
}

func championListHandler(w http.ResponseWriter, r *http.Request) {
	data := &champion.IndexData{
		ChampionData: []champion.RowData{
			{
				ChampionImagePath: "https://static.bigbrain.gg/assets/lol/riot_static/15.24.1/img/champion/Karthus.webp",
				KDA:               "4.5 / 5.7 / 2.3",
				WinRate:           "75% (10)",
			},
			{
				ChampionImagePath: "https://static.bigbrain.gg/assets/lol/riot_static/15.24.1/img/champion/Garen.webp",
				KDA:               "3.8 / 4.2 / 1.9",
				WinRate:           "68% (8)",
			},
		},
		AverageKDA: "4.2 / 4.5 / 2.1",
		WinRate:    "71% (18)",
	}

	champion.Index(r.Context(), data).Render(r.Context(), w)
}

func matchDetailBoxHandler(w http.ResponseWriter, r *http.Request) {
	data := &match_detail.IndexData{
		Teams: [2][]match_detail.ScoreCardData{
			{
				{
					Name:                       "FLY Vicla",
					Rank:                       "C1",
					ChampionImagePath:          "https://static.bigbrain.gg/assets/lol/riot_static/15.24.1/img/champion/Karthus.webp",
					ChampionLevel:              "15",
					SummonerSpellImagePaths:    [2]string{},
					Kills:                      "4",
					Deaths:                     "4",
					Assists:                    "4",
					CS:                         "123",
					CSPerMinute:                "9.8",
					RuneKeystoneImagePath:      "",
					RuneSecondaryTreeImagePath: "",
					ItemImagePaths:             [7]string{},
					VisionScore:                "12",
				},
				{
					Name:                       "C9 Summit",
					Rank:                       "C1",
					ChampionImagePath:          "https://static.bigbrain.gg/assets/lol/riot_static/15.24.1/img/champion/Darius.webp",
					ChampionLevel:              "15",
					SummonerSpellImagePaths:    [2]string{},
					Kills:                      "4",
					Deaths:                     "4",
					Assists:                    "4",
					CS:                         "123",
					CSPerMinute:                "9.8",
					RuneKeystoneImagePath:      "",
					RuneSecondaryTreeImagePath: "",
					ItemImagePaths:             [7]string{},
					VisionScore:                "12",
				},
			},
			{
				{
					Name:                       "T1 Faker",
					Rank:                       "Challenger 1409LP",
					ChampionImagePath:          "https://static.bigbrain.gg/assets/lol/riot_static/15.24.1/img/champion/Vayne.webp",
					ChampionLevel:              "15",
					SummonerSpellImagePaths:    [2]string{},
					Kills:                      "4",
					Deaths:                     "4",
					Assists:                    "4",
					CS:                         "123",
					CSPerMinute:                "9.8",
					RuneKeystoneImagePath:      "",
					RuneSecondaryTreeImagePath: "",
					ItemImagePaths:             [7]string{},
					VisionScore:                "12",
				},
			},
		},
	}

	match_detail.Index(r.Context(), data).Render(r.Context(), w)
}

func updateProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}
