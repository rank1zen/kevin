package riotmapper

import (
	"fmt"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

// MapLiveMatch maps riot.LiveMatch to a internal.LiveMatch.
func MapLiveMatch(riotMatch *riot.LiveMatch) *internal.LiveMatch {
	livematch := internal.LiveMatch{
		ID:           fmt.Sprintf("%s_%d", riotMatch.PlatformID, riotMatch.GameID),
		Date:         convertRiotUnixToTimestamp(riotMatch.GameStartTime),
		Participants: [10]internal.LiveParticipant{},
	}

	for i := range livematch.Participants {
		riotParticipant := riotMatch.Participants[i]
		liveparticipant := internal.LiveParticipant{
			PUUID:        riot.PUUID(riotMatch.Participants[i].PUUID),
			MatchID:      livematch.ID,
			ChampionID:   riotParticipant.ChampionID,
			Runes:        newRunePage(withRiotSpectatorPerks(&riotParticipant.Perks)),
			SummonersIDs: convertRiotLiveSummonerSpells(riotParticipant),
			TeamID:       riotParticipant.TeamID,
		}

		livematch.Participants[i] = liveparticipant
	}

	return &livematch
}

func withRiotSpectatorPerks(perks *riot.LivePerks) RunePageOption {
	return func(runes *internal.RunePage) error {
		runes.PrimaryTree = perks.PerkStyle
		runes.PrimaryKeystone = perks.PerkIDs[0]
		// runes.PrimaryA = perks.PerkIDs[1]
		// runes.PrimaryB = perks.PerkIDs[2]
		// runes.PrimaryC = perks.PerkIDs[3]
		runes.SecondaryTree = perks.PerkSubStyle
		// runes.SecondaryA = perks.PerkIDs[4]
		// runes.SecondaryB = perks.PerkIDs[5]
		// runes.MiniOffense = perks.PerkIDs[6]
		// runes.MiniFlex = perks.PerkIDs[7]
		// runes.MiniDefense = perks.PerkIDs[8]
		return nil
	}
}
