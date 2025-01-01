package internal

import "time"

type Participant struct {
	Puuid             PUUID
	ID                ParticipantID
	MatchID           MatchID
	TeamID            TeamID
	StartTimestamp    time.Time
	EndTimestamp      time.Time
	Duration          time.Duration
	Patch             GameVersion
	Position          string
	Win               bool
	BannedChampion    *ChampionID
	Champion          ChampionID
	ChampionLevel     int
	Summs             SummsIDs
	Items             ItemIDs
	Runes             Runes
	Kills             int
	Deaths            int
	Assists           int
	KillParticipation float32
	CreepScore        int
	CsPerMinute       float32
	Gold              int
	GoldPercentage    float32
	GoldDelta         int
	Damage            int
	DamagePercentage  float32
	DamageDelta       int
	VisionScore       int
	PinkWards         int
}

type MatchHistory struct {
	Page    int
	Count   int
	HasMore bool
	List    []Participant
}

type MatchParticipants [10]Participant
