package internal

import "time"

// RiotMatch is strictly a ranked, soloq, 5v5, match on summoners rift.
type RiotMatch struct {
	ID              MatchID
	DataVersion     string
	Patch           GameVersion
	CreateTimestamp time.Time
	StartTimestamp  time.Time
	EndTimestamp    time.Time
	Duration        time.Duration
	EndOfGameResult string
	Participants    MatchParticipantList
}

type MatchParticipantList [10]RiotMatchParticipant

type MatchTeamList [5]RiotMatchParticipant

func (m *RiotMatch) GetParticipants() MatchParticipantList {
	return m.Participants
}

func (m *RiotMatch) GetTeams() [2]RiotMatchTeam {
	return [2]RiotMatchTeam{}
}

type RiotMatchParticipant struct {
	ID       ParticipantID
	Puuid    PUUID
	Match    MatchID
	Team     TeamID
	Summoner SummonerID
	Patch    GameVersion

	SummonerLevel             int
	SummonerName              string
	RiotIDGameName            string
	RiotIDName                string
	RiotIDTagline             string
	ChampionLevel             int
	ChampionID                ChampionID
	ChampionName              string
	GameEndedInEarlySurrender bool
	GameEndedInSurrender      bool
	Items                     ItemIDs
	Runes                     Runes
	Role                      string
	Summoners                 SummsIDs
	TeamEarlySurrendered      bool
	TeamPosition              string
	TimePlayed                int
	Win                       bool

	// Post game stats

	Assists                        int
	DamageDealtToBuildings         int
	DamageDealtToObjectives        int
	DamageDealtToTurrets           int
	DamageSelfMitigated            int
	Deaths                         int
	DetectorWardsPlaced            int
	FirstBloodAssist               bool
	FirstBloodKill                 bool
	FirstTowerAssist               bool
	FirstTowerKill                 bool
	GoldEarned                     int
	GoldSpent                      int
	IndividualPosition             string
	InhibitorKills                 int
	InhibitorTakedowns             int
	InhibitorsLost                 int
	Kills                          int
	MagicDamageDealt               int
	MagicDamageDealtToChampions    int
	MagicDamageTaken               int
	PhysicalDamageDealt            int
	PhysicalDamageDealtToChampions int
	PhysicalDamageTaken            int
	SightWardsBoughtInGame         int
	TotalDamageDealt               int
	TotalDamageDealtToChampions    int
	TotalDamageShieldedOnTeammates int
	TotalDamageTaken               int
	TotalHeal                      int
	TotalHealsOnTeammates          int
	TotalMinionsKilled             int
	TrueDamageDealt                int
	TrueDamageDealtToChampions     int
	TrueDamageTaken                int
	VisionScore                    int
	VisionWardsBoughtInGame        int
	WardsKilled                    int
	WardsPlaced                    int
	NeutralMinionsKilled           int
}

type RiotMatchTeam struct {
	ID      TeamID
	MatchID MatchID
	Win     bool
}
