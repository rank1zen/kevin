package db

import (
	"time"
)

type MatchKillEvent struct {
	Timestamp time.Duration `db:"timestamp"`
	PosX      int           `db:"pos_x"`
	PosY      int           `db:"pos_y"`
	Bounty    int           `db:"bounty"`
	Shutdown  int           `db:"shutdown_bounty"`
	KillerID  int           `db:"killer_id"`
	VictimID  int           `db:"victim_id"`
}

type MatchKillEventList []*MatchKillEvent

type MatchItemEvent struct {
	Timestamp     time.Duration `db:"timestamp"`
	ParticipantID int
	SkillSlot     int
}

type MatchItemEventList []*MatchItemEvent

type MatchSkillEvent struct {
	Timestamp     time.Duration `db:"timestamp"`
	ParticipantID int
	ItemID        int
}

type MatchSkillEventList []*MatchSkillEvent

type SummonerMatchItemEvent struct {
	Timestamp     time.Duration
	ParticipantId int
	ItemId        int
}


type SummonerMatchKillEvent struct{}


type SummonerMatchSkillEvent struct {
	Timestamp     time.Duration `db:"timestamp"`
	ParticipantId int           `db:"participant_id"`
	SkillSlot     int           `db:"skill_slot"`
}

type SummonerMatchSkillEventList []*SummonerMatchSkillEvent
