// sample provides sample objects used for testing
package sample

import (
	"embed"
	"encoding/json"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

//go:embed samples
var content embed.FS

// WithSampleMatch instantiates some valid [internal.Match], usually used for
// testing.
func WithSampleMatch() internal.MatchOption {
	matchFile, err := content.Open("samples/NA1_5304757838.json")
	if err != nil {
		panic(err)
	}

	var riotMatch riot.Match
	err = json.NewDecoder(matchFile).Decode(&riotMatch)
	if err != nil {
		panic(err)
	}

	return internal.WithRiotMatch(&riotMatch)
}

func WithMatchID(id string) internal.MatchOption {
	return func(m *internal.Match) error {
		m.ID = id
		for i := range m.Participants {
			m.Participants[i].MatchID = id
		}

		return nil
	}
}
