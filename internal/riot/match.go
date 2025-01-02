package riot

import (
	"context"

	"github.com/rank1zen/yujin/internal"
	"github.com/rank1zen/yujin/internal/riotclient"
)

func makeItems(item0, item1, item2, item3, item4, item5, item6 int) internal.ItemIDs {
	ids := internal.ItemIDs{}

	for i, item := range []int{item0, item1, item2, item3, item4, item5, item6} {
		if item == 0 {
			ids[i] = nil
		} else {
			id := internal.ItemID(item)
			ids[i] = &id
		}
	}

	return ids
}

func makeRunes(perks *riotclient.MatchPerks) internal.Runes {
	return internal.Runes{
		PrimaryTree:     internal.RuneID(perks.Styles[0].Style),
		PrimaryKeystone: internal.RuneID(perks.Styles[0].Selections[0].Perk),
		PrimaryA:        internal.RuneID(perks.Styles[0].Selections[1].Perk),
		PrimaryB:        internal.RuneID(perks.Styles[0].Selections[2].Perk),
		PrimaryC:        internal.RuneID(perks.Styles[0].Selections[3].Perk),
		SecondaryTree:   internal.RuneID(perks.Styles[1].Selections[0].Perk),
		SecondaryA:      internal.RuneID(perks.Styles[1].Selections[1].Perk),
		SecondaryB:      internal.RuneID(perks.Styles[1].Selections[2].Perk),
		MiniOffense:     internal.RuneID(perks.StatPerks.Offense),
		MiniFlex:        internal.RuneID(perks.StatPerks.Flex),
		MiniDefense:     internal.RuneID(perks.StatPerks.Defense),
	}
}

func (r *Riot) GetMatch(ctx context.Context, id internal.MatchID) (internal.RiotMatch, error) {
	m, err := r.client.GetMatch(ctx, id.String())
	if err != nil {
		return internal.RiotMatch{}, err
	}

	return internal.RiotMatch{
		ID:              internal.MatchID(m.Metadata.MatchId),
		DataVersion:     m.Metadata.DataVersion,
		Patch:           internal.GameVersion(m.Info.GameVersion),
		CreateTimestamp: riotUnixToDate(m.Info.GameCreation),
		StartTimestamp:  riotUnixToDate(m.Info.GameStartTimestamp),
		EndTimestamp:    riotUnixToDate(m.Info.GameEndTimestamp),
		Duration:        riotDurationToInterval(int(m.Info.GameDuration)),
		EndOfGameResult: m.Info.EndOfGameResult,
		// TODO: Participants:    nil,
	}, nil
}
