package internal

import (
	"context"
	"encoding/json"
	"os"

	"github.com/rank1zen/kevin/internal/riot"
)

type Sample struct {
}

func NewSample() {

}

func (s *Sample) GetMatch(ctx context.Context, region, matchID string) (*riot.Match, error) {
	testdata := os.DirFS("./testdata")


	matchFile, err := testdata.Open("NA1_5304757838.json")
	if err != nil {
		return nil, err
	}

	var riotMatch riot.Match
	err = json.NewDecoder(matchFile).Decode(&riotMatch)
	if err != nil {
		return nil, err
	}

	return &riotMatch, nil
}
