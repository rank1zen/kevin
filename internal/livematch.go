package internal

import "time"

type LiveMatchParticipant struct {
	Puuid          PUUID
	Date           time.Time
	TeamID         TeamID
	SummonerID     SummonerID
	Champion       ChampionID
	Summoners      SummsIDs
	Runes          Runes
	BannedChampion *ChampionID
}

type LiveMatchParticipantList [10]LiveMatchParticipant

type LiveMatchTeamList [5]LiveMatchParticipant

type RiotLiveMatch struct {
	StartTimestamp time.Time
	Length         time.Duration
	IDs            [10]PUUID
	Participant    LiveMatchParticipantList
}

// TODO: we should further implement checking the order of the list
func (m *RiotLiveMatch) GetParticipants() LiveMatchParticipantList {
	return m.Participant
}
