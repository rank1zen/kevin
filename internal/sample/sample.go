// sample provides sample objects used for testing
package sample

import (
	"embed"
	"encoding/json"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

var (
	SummonerOrrangeNA1 = internal.Summoner{
		PUUID:   "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q",
		Name:    "orrange",
		Tagline: "NA1",
	}

	SummonerT1OKGOODYESNA1 = internal.Summoner{
		PUUID:   "44Js96gJP_XRb3GpJwHBbZjGZmW49Asc3_KehdtVKKTrq3MP8KZdeIn_27MRek9FkTD-M4_n81LNqg",
		Name:    "T1 OK GOOD YES",
		Tagline: "NA1",
	}
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
