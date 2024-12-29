package ddragon

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"

	"github.com/rank1zen/yujin/internal/logging"
)

// static data and images are hard af
// FIXME: better for this

//go:embed assets/*.json
var assets embed.FS

const (
	CdragonCdnBaseUrl = "https://cdn.communitydragon.org"
	CdragonRawBaseUrl = "https://raw.communitydragon.org"
)

type SummonerSpell struct {
	Name     string `json:"name"`
	IconPath string `json:"iconPath"`
	Id       int    `json:"id"`
}

type Perk struct {
	Name     string `json:"name"`
	IconPath string `json:"iconPath"`
	Id       int    `json:"id"`
}

var (
	IdToSpell = map[int]SummonerSpell{}
	IdToPerk  = map[int]Perk{}
)

func init() {
	f, _ := fs.Sub(assets, "assets")
	spells, err := f.Open("summoner-spells.json")
	if err != nil {
		logging.Get().Sugar().DPanicf("opening summoner spells: %v", err)
	}

	perks, err := f.Open("perks.json")
	if err != nil {
		logging.Get().Sugar().DPanicf("opening summoner spells: %v", err)
	}

	var m []SummonerSpell
	err = json.NewDecoder(spells).Decode(&m)
	if err != nil {
		logging.Get().Sugar().DPanicf("decoding: %v", err)
	}

	var p []Perk
	err = json.NewDecoder(perks).Decode(&p)
	if err != nil {
		logging.Get().Sugar().DPanicf("decoding: %v", err)
	}

	for _, spell := range m {
		IdToSpell[spell.Id] = spell
	}

	for _, perk := range p {
		IdToPerk[perk.Id] = perk
	}
}

func GetSummonerSpellUrl(spellId int) string {
	path, _ := strings.CutPrefix(IdToSpell[spellId].IconPath, "/lol-game-data/assets")
	return "https://raw.communitydragon.org/latest/plugins/rcp-be-lol-game-data/global/default" + strings.ToLower(path)
}

func GetRuneIconUrl(runeId int) string {
	path, _ := strings.CutPrefix(IdToPerk[runeId].IconPath, "/lol-game-data/assets")
	return "https://raw.communitydragon.org/latest/plugins/rcp-be-lol-game-data/global/default" + strings.ToLower(path)
}

func GetChampionIconUrl(id int) string {
	return fmt.Sprintf(CdragonCdnBaseUrl+"/14.16.1/champion/%d/square", id)
}

func GetItemIconUrl(id int) string {
	return fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/14.16.1/img/item/%d.png", id)
}

var runeTreeIcons = map[int]string{
	8000: CdragonRawBaseUrl + "/14.16/plugins/rcp-be-lol-game-data/global/default/v1/perk-images/styles/7201_precision.png",
	8100: CdragonRawBaseUrl + "/14.16/plugins/rcp-be-lol-game-data/global/default/v1/perk-images/styles/7200_domination.png",
	8200: CdragonRawBaseUrl + "/14.16/plugins/rcp-be-lol-game-data/global/default/v1/perk-images/styles/7202_sorcery.png",
	8300: CdragonRawBaseUrl + "/14.16/plugins/rcp-be-lol-game-data/global/default/v1/perk-images/styles/7203_whimsy.png",
	8400: CdragonRawBaseUrl + "/14.16/plugins/rcp-be-lol-game-data/global/default/v1/perk-images/styles/7204_resolve.png",
}

func GetRuneTreeIconUrl(id int) string {
	return runeTreeIcons[id]
}
