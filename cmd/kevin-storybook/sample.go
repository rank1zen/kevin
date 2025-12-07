package main

import (
	"time"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/frontend/page"
	"github.com/rank1zen/kevin/internal/frontend/view/livematch/teamsection"
	"github.com/rank1zen/kevin/internal/frontend/view/profile"
	"github.com/rank1zen/kevin/internal/frontend/view/search"
	"github.com/rank1zen/kevin/internal/frontend/view/shared"
	"github.com/rank1zen/kevin/internal/riot"
)

var sampleRunePage = internal.RunePage{
	PrimaryTree:     8100,
	PrimaryKeystone: 8112,
	PrimaryA:        8143,
	PrimaryB:        8140,
	PrimaryC:        8135,
	SecondaryTree:   8200,
	SecondaryA:      8234,
	SecondaryB:      8237,
	MiniOffense:     5005,
	MiniFlex:        5008,
	MiniDefense:     5011,
}

var sampleItems = [7]int{6698, 0, 3176, 0, 3134, 3814, 3364}

var sampleSpells = [2]int{4, 6}

var sampleRank = internal.Rank{
	Tier:     riot.TierGrandmaster,
	Division: riot.Division1,
	LP:       923,
}

var (
	rankCardData = profile.RankCardData{
		RankDetail: &internal.RankDetail{
			Wins:   41,
			Losses: 67,
			Rank:   sampleRank,
		},
	}

	historyCardData = profile.HistoryCardData{
		MatchID:        "",
		ChampionID:     16,
		ChampionLevel:  16,
		SummonerIDs:    sampleSpells,
		Kills:          3,
		Deaths:         4,
		Assists:        5,
		KillDeathRatio: 3.4111,
		CS:             41,
		CSPerMinute:    6.7,
		RunePage:       sampleRunePage,
		Items:          sampleItems,
		VisionScore:    41,
		RankChange:     nil,
		LPChange:       nil,
		Win:            false,
	}

	championListData = profile.ChampionListData{
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
	}

	redSideParticipantData = profile.MatchParticipantCardData{
		MatchID:        "1",
		Name:           "Joe",
		Tag:            "NA1",
		PUUID:          "1",
		ChampionID:     41,
		ChampionLevel:  18,
		SummonerIDs:    sampleSpells,
		Kills:          41,
		Deaths:         42,
		Assists:        67,
		KillDeathRatio: 41.67,
		CS:             410,
		CSPerMinute:    4.1,
		RunePage:       sampleRunePage,
		Items:          sampleItems,
		VisionScore:    23,
		Rank:           nil,
	}

	blueSideParticipantData1 = profile.MatchParticipantCardData{
		MatchID:        "1",
		Name:           "Joe",
		Tag:            "NA1",
		PUUID:          "1",
		ChampionID:     41,
		ChampionLevel:  18,
		SummonerIDs:    sampleSpells,
		Kills:          41,
		Deaths:         42,
		Assists:        67,
		KillDeathRatio: 41.67,
		CS:             410,
		CSPerMinute:    4.1,
		RunePage:       sampleRunePage,
		Items:          sampleItems,
		VisionScore:    23,
		Rank:           nil,
	}

	blueSideParticipantData2 = profile.MatchParticipantCardData{
		MatchID:        "1",
		Name:           "Bartholomew Montgomery",
		Tag:            "NA1",
		PUUID:          "1",
		ChampionID:     42,
		ChampionLevel:  18,
		SummonerIDs:    sampleSpells,
		Kills:          42,
		Deaths:         42,
		Assists:        67,
		KillDeathRatio: 42.67,
		CS:             420,
		CSPerMinute:    4.2,
		RunePage:       sampleRunePage,
		Items:          sampleItems,
		VisionScore:    42,
		Rank:           nil,
	}

	matchDetailBoxData = profile.MatchDetailBoxData{
		Date:     time.Date(2041, time.June, 7, 6, 7, 6, 0, time.UTC),
		Duration: 41 * 67 * time.Second,
		RedSide: profile.MatchTeamListData{
			Participants: []profile.MatchParticipantCardData{
				redSideParticipantData,
			},
		},
		BlueSide: profile.MatchTeamListData{
			Participants: []profile.MatchParticipantCardData{
				blueSideParticipantData1,
				blueSideParticipantData2,
			},
		},
	}

	partialHistoryEntryData = profile.PartialHistoryEntryData{
		PUUID: "bacon-egg-and-cheese",
	}

	blueSideLiveMatchParticipantData = profile.LiveMatchParticipantCardData{
		MatchID:        "1",
		Name:           "Bartholomew Montgomery",
		Tag:            "NA1",
		PUUID:          "1",
		ChampionID:     42,
		ChampionLevel:  18,
		SummonerIDs:    sampleSpells,
		Kills:          42,
		Deaths:         42,
		Assists:        67,
		KillDeathRatio: 42.67,
		CS:             420,
		CSPerMinute:    4.2,
		RunePage:       sampleRunePage,
		Rank: &internal.Rank{
			Tier:     riot.TierGrandmaster,
			Division: riot.Division1,
			LP:       900,
		},
	}

	redSideLiveMatchParticipantData = profile.LiveMatchParticipantCardData{
		MatchID:        "1",
		Name:           "Bartholomew Montgomery",
		Tag:            "NA1",
		PUUID:          "1",
		ChampionID:     42,
		ChampionLevel:  18,
		SummonerIDs:    sampleSpells,
		Kills:          42,
		Deaths:         42,
		Assists:        67,
		KillDeathRatio: 42.67,
		CS:             420,
		CSPerMinute:    4.2,
		RunePage:       sampleRunePage,
		Rank:           nil,
	}

	liveMatchScoreboardData = profile.LiveMatchScoreboardData{
		BlueSideParticipants: []profile.LiveMatchParticipantCardData{
			blueSideLiveMatchParticipantData,
		},
		RedSideParticipants: []profile.LiveMatchParticipantCardData{
			redSideLiveMatchParticipantData,
		},
	}

	profilePageData = page.ProfilePageData{
		PUUID:  "bacon-egg-and-cheese",
		Region: riot.RegionNA1,
		Name:   "Bacon",
		Tag:    "41",
	}

	profileLiveMatchPageData = page.ProfileLiveMatchPageData{
		PUUID:  "bacon-egg-and-cheese",
		Region: riot.RegionNA1,
		Name:   "Bacon",
		Tag:    "41",
	}

	sharedHeaderData = shared.HeaderData{
		Region: riot.RegionNA1,
	}

	searchResultCardData = search.ResultCardData{
		Name: "Bartholomew Montgomery",
		Tag:  "NA1",
		Rank: &internal.Rank{
			Tier:     riot.TierBronze,
			Division: riot.Division4,
			LP:       67,
		},
	}

	searchResultMenuCardData1 = search.ResultCardData{
		Name: "Bartholomew Montgomery",
		Tag:  "NA1",
		Rank: &internal.Rank{
			Tier:     riot.TierBronze,
			Division: riot.Division4,
			LP:       67,
		},
	}

	searchResultMenuCardData2 = search.ResultCardData{
		Name: "Bartholomew",
		Tag:  "NA1",
		Rank: &internal.Rank{
			Tier:     riot.TierBronze,
			Division: riot.Division4,
			LP:       67,
		},
	}

	searchResultMenuData = search.ResultMenuData{
		Name: "T1 OK GOOD YES",
		Tag:  "NA1",
		ProfileResults: []search.ResultCardData{
			searchResultMenuCardData1,
			searchResultMenuCardData2,
		},
	}
)

var sampleTeamsectionParticipantCardData = teamsection.ParticipantCardData{
	PUUID:            "bacon-egg-and-cheese",
	Name:             "Bacon",
	Tag:              "41",
	MatchID:          "NA1_1234567890",
	ChampionID:       41,
	SummonerIDs:      sampleSpells,
	RunePage:         sampleRunePage,
	Rank:             &sampleRank,
	BannedChampionID: 41,
}

var livematchPageData = page.LivematchPageData{
	Region:  riot.RegionNA1,
	MatchID: "NA1_1234567890",
	Teams: []teamsection.TeamsectionData{
		{
			TeamName: "Blue Side",
			Participants: []teamsection.ParticipantCardData{
				sampleTeamsectionParticipantCardData,
				sampleTeamsectionParticipantCardData,
				sampleTeamsectionParticipantCardData,
				sampleTeamsectionParticipantCardData,
				sampleTeamsectionParticipantCardData,
			},
		},
		{
			TeamName: "Red Side",
			Participants: []teamsection.ParticipantCardData{
				sampleTeamsectionParticipantCardData,
				sampleTeamsectionParticipantCardData,
				sampleTeamsectionParticipantCardData,
				sampleTeamsectionParticipantCardData,
				sampleTeamsectionParticipantCardData,
			},
		},
	},
}
