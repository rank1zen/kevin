package main

import (
	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

var sampleRunePage = internal.RunePage{
	PrimaryTree:     8100,
	PrimaryKeystone: 8112,
	PrimaryA:        8143,
	PrimaryB:        8140,
	PrimaryC:        8135,
	SecondaryTree:   8200,
	SecondaryA:      8234,
	SecondaryB:      8237,
	MiniOffense:     5005,
	MiniFlex:        5008,
	MiniDefense:     5011,
}

var sampleItems = [7]int{6698, 0, 3176, 0, 3134, 3814, 3364}

var sampleSpells = [2]int{4, 5}

var sampleRank = internal.Rank{
	Tier:     riot.TierGrandmaster,
	Division: riot.Division1,
	LP:       923,
}
