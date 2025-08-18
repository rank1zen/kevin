// sample provides sample objects used for testing
package sample

import (
	"embed"
	"encoding/json"

	"github.com/rank1zen/kevin/internal/riot"
)

//go:embed samples
var content embed.FS

// WithSampleMatch instantiates some valid [internal.Match], usually used for
// testing.
func WithSampleMatch() riot.Match {
	matchFile, err := content.Open("samples/NA1_5304757838.json")
	if err != nil {
		panic(err)
	}

	var riotMatch riot.Match
	err = json.NewDecoder(matchFile).Decode(&riotMatch)
	if err != nil {
		panic(err)
	}

	return riotMatch
}

// Match5346312088 returns a real sample match.
//
// https://op.gg/lol/summoners/na/T1%2520OK%2520GOOD%2520YES-NA1/matches/9xQlqPbXBXYkNDjORlI3_vOXw2r6bzl277RNnqi7xqk%3D/1755120925000
func Match5346312088() riot.Match {
	matchFile, err := content.Open("samples/NA1_5346312088.json")
	if err != nil {
		panic(err)
	}

	var riotMatch riot.Match
	err = json.NewDecoder(matchFile).Decode(&riotMatch)
	if err != nil {
		panic(err)
	}

	return riotMatch
}

func Match5347748140() riot.Match {
	matchFile, err := content.Open("samples/NA1_5347748140.json")
	if err != nil {
		panic(err)
	}

	var riotMatch riot.Match
	err = json.NewDecoder(matchFile).Decode(&riotMatch)
	if err != nil {
		panic(err)
	}

	return riotMatch
}

func WithSampleLiveMatch() riot.LiveMatch {
	file, err := content.Open("samples/live_match.json")
	if err != nil {
		panic(err)
	}

	var riotMatch riot.LiveMatch
	err = json.NewDecoder(file).Decode(&riotMatch)
	if err != nil {
		panic(err)
	}

	return riotMatch
}
