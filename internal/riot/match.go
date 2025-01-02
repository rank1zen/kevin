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

func makeSumms(summ0, summ1 int) internal.SummsIDs {
	ids := internal.SummsIDs{}

	for i, summ := range []int{summ0, summ1} {
		// maybe assert summ is zero
		ids[i] = internal.SummsID(summ)
	}

	return ids
}

func (r *Riot) GetMatch(ctx context.Context, id internal.MatchID) (internal.RiotMatch, error) {
	m, err := r.client.GetMatch(ctx, id.String())
	if err != nil {
		return internal.RiotMatch{}, err
	}

	participants := internal.MatchParticipantList{}

	for i := range len(m.Info.Participants) {
		p := m.Info.Participants[i]

		participants[i] = internal.RiotMatchParticipant{
			ID:                             internal.ParticipantID(p.ParticipantID),
			Puuid:                          internal.PUUID(p.PUUID),
			Match:                          internal.MatchID(m.Metadata.MatchId),
			Team:                           internal.TeamID(p.TeamID),
			Summoner:                       internal.SummonerID(p.SummonerID),
			Patch:                          internal.GameVersion(m.Info.GameVersion),
			SummonerLevel:                  p.SummonerLevel,
			SummonerName:                   p.SummonerName,
			RiotIDGameName:                 p.RiotIDGameName,
			RiotIDName:                     p.RiotIDName,
			RiotIDTagline:                  p.RiotIDTagline,
			ChampionLevel:                  p.ChampLevel,
			ChampionID:                     internal.ChampionID(p.ChampionID),
			ChampionName:                   p.ChampionName,
			GameEndedInEarlySurrender:      p.GameEndedInEarlySurrender,
			GameEndedInSurrender:           p.GameEndedInSurrender,
			Role:                           p.Role,
			TeamEarlySurrendered:           p.TeamEarlySurrendered,
			TeamPosition:                   p.TeamPosition,
			TimePlayed:                     p.TimePlayed,
			Win:                            p.Win,
			Items:                          makeItems(p.Item0, p.Item1, p.Item2, p.Item3, p.Item4, p.Item5, p.Item6),
			Runes:                          makeRunes(p.Perks),
			Summoners:                      makeSumms(p.Summoner1ID, p.Summoner2ID),
			Assists:                        p.Assists,
			DamageDealtToBuildings:         p.DamageDealtToBuildings,
			DamageDealtToObjectives:        p.DamageDealtToObjectives,
			DamageDealtToTurrets:           p.DamageDealtToTurrets,
			DamageSelfMitigated:            p.DamageSelfMitigated,
			Deaths:                         p.Deaths,
			DetectorWardsPlaced:            p.DetectorWardsPlaced,
			FirstBloodAssist:               p.FirstBloodAssist,
			FirstBloodKill:                 p.FirstBloodKill,
			FirstTowerAssist:               p.FirstTowerAssist,
			FirstTowerKill:                 p.FirstTowerKill,
			GoldEarned:                     p.GoldEarned,
			GoldSpent:                      p.GoldSpent,
			IndividualPosition:             p.IndividualPosition,
			InhibitorKills:                 p.InhibitorKills,
			InhibitorTakedowns:             p.InhibitorTakedowns,
			InhibitorsLost:                 p.InhibitorsLost,
			Kills:                          p.Kills,
			MagicDamageDealt:               p.MagicDamageDealt,
			MagicDamageDealtToChampions:    p.MagicDamageDealtToChampions,
			MagicDamageTaken:               p.MagicDamageTaken,
			PhysicalDamageDealt:            p.PhysicalDamageDealt,
			PhysicalDamageDealtToChampions: p.PhysicalDamageDealtToChampions,
			PhysicalDamageTaken:            p.PhysicalDamageTaken,
			SightWardsBoughtInGame:         p.SightWardsBoughtInGame,
			TotalDamageDealt:               p.TotalDamageDealt,
			TotalDamageDealtToChampions:    p.TotalDamageDealtToChampions,
			TotalDamageShieldedOnTeammates: p.TotalDamageShieldedOnTeammates,
			TotalDamageTaken:               p.TotalDamageTaken,
			TotalHeal:                      p.TotalHeal,
			TotalHealsOnTeammates:          p.TotalHealsOnTeammates,
			TotalMinionsKilled:             p.TotalMinionsKilled,
			TrueDamageDealt:                p.TrueDamageDealt,
			TrueDamageDealtToChampions:     p.TrueDamageDealtToChampions,
			TrueDamageTaken:                p.TrueDamageTaken,
			VisionScore:                    p.VisionScore,
			VisionWardsBoughtInGame:        p.VisionWardsBoughtInGame,
			WardsKilled:                    p.WardsKilled,
			WardsPlaced:                    p.WardsPlaced,
			NeutralMinionsKilled:           p.NeutralMinionsKilled,
		}
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
		Participants:    participants,
	}, nil
}
