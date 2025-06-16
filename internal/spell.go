package internal

import "time"

type Spell int

type SpellEvent struct {
	PUUID           int
	InGameTimestamp time.Duration
	SpellSlot       int
}

type SkillProgression [18]SpellEvent

func GetSpellSprite(id Spell) (u string, x, y, w, h int) {
	sprite := spellSprite[id]
	x = (sprite.Index % 10) * 24
	y = (sprite.Index / 10) * 24
	return sprite.Sprite, x, y, 24, 24
}

var spellSprite = map[Spell]struct {
	Sprite string
	Index  int
}{
	21:   {"tiny_spell0.png", 0},
	1:    {"tiny_spell0.png", 1},
	2202: {"tiny_spell0.png", 2},
	2201: {"tiny_spell0.png", 3},
	14:   {"tiny_spell0.png", 4},
	3:    {"tiny_spell0.png", 5},
	4:    {"tiny_spell0.png", 6},
	6:    {"tiny_spell0.png", 7},
	7:    {"tiny_spell0.png", 8},
	13:   {"tiny_spell0.png", 9},
	30:   {"tiny_spell0.png", 10},
	31:   {"tiny_spell0.png", 11},
	11:   {"tiny_spell0.png", 12},
	39:   {"tiny_spell0.png", 13},
	32:   {"tiny_spell0.png", 14},
	12:   {"tiny_spell0.png", 15},
	54:   {"tiny_spell0.png", 16},
	55:   {"tiny_spell0.png", 17},
}
