package internal

import (
	"time"

	"github.com/rank1zen/kevin/internal/riot"
)

type SkillEvent struct {
	MatchID         string
	PUUID           int
	InGameTimestamp time.Duration
	SpellSlot       int
}

type SkillProgression [18]*int

func NewSkillProgression(opts ...SkillProgressionOption) (m SkillProgression) {
	for _, f := range opts {
		f(&m)
	}
	return m
}

type SkillProgressionOption func(*SkillProgression) error

func WithRiotTimeline(timeline *riot.Timeline, puuid string) SkillProgressionOption {
	// TODO: implement
	return func(m *SkillProgression) error {
		return nil
	}
}
