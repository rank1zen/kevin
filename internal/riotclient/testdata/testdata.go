package testdata

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"

	"github.com/rank1zen/yujin/internal/riotclient"
)

//go:embed match/*.json
var matchJsonFiles embed.FS

func init() {
	files, err := fs.Glob(matchJsonFiles, "match/*.json")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		data, err := matchJsonFiles.ReadFile(file)
		if err != nil {
			log.Fatalf("reading file: %v", err)
		}

		var match riotclient.Match

		err = json.Unmarshal(data, &match)
		if err != nil {
			log.Fatalf("scanning: %v", err)
		}

		matchMap[match.Metadata.MatchId] = &match
	}
}

// Get match by match id
var matchMap = map[string]*riotclient.Match{}

func GetMatch(matchID string) *riotclient.Match {
	if matchMap[matchID] == nil {
		log.Fatalf("match does not exist: %v", matchID)
	}

	return matchMap[matchID]
}
