package internal

import "time"

// Match is strictly a ranked, soloq, 5v5, match on summoners rift.
type Match struct {
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

type MatchParticipantList [10]MatchParticipant

type MatchTeamList [5]MatchParticipant

func (m *Match) GetParticipants() MatchParticipantList {
	return m.Participants
}

func (m *Match) GetTeams() [2]MatchTeam {
	return [2]MatchTeam{}
}

type MatchParticipant struct {
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

type MatchTeam struct {
	ID      TeamID
	MatchID MatchID
	Win     bool
}
