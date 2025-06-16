package riot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Match struct {
	Info     *MatchInfo     `json:"info"`
	Metadata *MatchMetadata `json:"metadata"`
}

type MatchMetadata struct {
	DataVersion  string   `json:"dataVersion"`
	MatchId      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type MatchInfo struct {
	EndOfGameResult    string         `json:"endOfGameResult"`
	GameCreation       int64          `json:"gameCreation"`
	GameDuration       int64          `json:"gameDuration"`
	GameEndTimestamp   int64          `json:"gameEndTimestamp"`
	GameId             int64          `json:"gameId"`
	GameMode           string         `json:"gameMode"`
	GameName           string         `json:"gameName"`
	GameStartTimestamp int64          `json:"gameStartTimestamp"`
	GameType           string         `json:"gameType"`
	GameVersion        string         `json:"gameVersion"`
	MapId              int            `json:"mapId"`
	Participants       []*Participant `json:"participants"`
	PlatformId         string         `json:"platformId"`
	QueueId            int            `json:"queueId"`
	Teams              []*Team        `json:"teams"`
	TournamentCode     string         `json:"tournamentCode"`
}

type Participant struct {
	Assists                        int    `json:"assists"`
	BaronKills                     int    `json:"baronKills"`
	BountyLevel                    int    `json:"bountyLevel"`
	ChampExperience                int    `json:"champExperience"`
	ChampLevel                     int    `json:"champLevel"`
	ChampionId                     int    `json:"championId"`
	ChampionName                   string `json:"championName"`
	ChampionTransform              int    `json:"championTransform"`
	ConsumablesPurchased           int    `json:"consumablesPurchased"`
	DamageDealtToBuildings         int    `json:"damageDealtToBuildings"`
	DamageDealtToObjectives        int    `json:"damageDealtToObjectives"`
	DamageDealtToTurrets           int    `json:"damageDealtToTurrets"`
	DamageSelfMitigated            int    `json:"damageSelfMitigated"`
	Deaths                         int    `json:"deaths"`
	DetectorWardsPlaced            int    `json:"detectorWardsPlaced"`
	DoubleKills                    int    `json:"doubleKills"`
	DragonKills                    int    `json:"dragonKills"`
	FirstBloodAssist               bool   `json:"firstBloodAssist"`
	FirstBloodKill                 bool   `json:"firstBloodKill"`
	FirstTowerAssist               bool   `json:"firstTowerAssist"`
	FirstTowerKill                 bool   `json:"firstTowerKill"`
	GameEndedInEarlySurrender      bool   `json:"gameEndedInEarlySurrender"`
	GameEndedInSurrender           bool   `json:"gameEndedInSurrender"`
	GoldEarned                     int    `json:"goldEarned"`
	GoldSpent                      int    `json:"goldSpent"`
	IndividualPosition             string `json:"individualPosition"`
	InhibitorKills                 int    `json:"inhibitorKills"`
	InhibitorTakedowns             int    `json:"inhibitorTakedowns"`
	InhibitorsLost                 int    `json:"inhibitorsLost"`
	Item0                          int    `json:"item0"`
	Item1                          int    `json:"item1"`
	Item2                          int    `json:"item2"`
	Item3                          int    `json:"item3"`
	Item4                          int    `json:"item4"`
	Item5                          int    `json:"item5"`
	Item6                          int    `json:"item6"`
	ItemsPurchased                 int    `json:"itemsPurchased"`
	KillingSprees                  int    `json:"killingSprees"`
	Kills                          int    `json:"kills"`
	Lane                           string `json:"lane"`
	LargestCriticalStrike          int    `json:"largestCriticalStrike"`
	LargestKillingSpree            int    `json:"largestKillingSpree"`
	LargestMultiKill               int    `json:"largestMultiKill"`
	LongestTimeSpentLiving         int    `json:"longestTimeSpentLiving"`
	MagicDamageDealt               int    `json:"magicDamageDealt"`
	MagicDamageDealtToChampions    int    `json:"magicDamageDealtToChampions"`
	MagicDamageTaken               int    `json:"magicDamageTaken"`
	NeutralMinionsKilled           int    `json:"neutralMinionsKilled"`
	NexusKills                     int    `json:"nexusKills"`
	NexusLost                      int    `json:"nexusLost"`
	NexusTakedowns                 int    `json:"nexusTakedowns"`
	ObjectivesStolen               int    `json:"objectivesStolen"`
	ObjectivesStolenAssists        int    `json:"objectivesStolenAssists"`
	ParticipantId                  int    `json:"participantId"`
	PentaKills                     int    `json:"pentaKills"`
	Perks                          *Perks `json:"perks"`
	PhysicalDamageDealt            int    `json:"physicalDamageDealt"`
	PhysicalDamageDealtToChampions int    `json:"physicalDamageDealtToChampions"`
	PhysicalDamageTaken            int    `json:"physicalDamageTaken"`
	ProfileIcon                    int    `json:"profileIcon"`
	PUUID                          string `json:"puuid"`
	QuadraKills                    int    `json:"quadraKills"`
	RiotIdGameName                 string `json:"riotIdGameName"`
	RiotIdName                     string `json:"riotIdName"`
	RiotIdTagline                  string `json:"riotIdTagline"`
	Role                           string `json:"role"`
	SightWardsBoughtInGame         int    `json:"sightWardsBoughtInGame"`
	Spell1Casts                    int    `json:"spell1Casts"`
	Spell2Casts                    int    `json:"spell2Casts"`
	Spell3Casts                    int    `json:"spell3Casts"`
	Spell4Casts                    int    `json:"spell4Casts"`
	Summoner1Casts                 int    `json:"summoner1Casts"`
	Summoner1Id                    int    `json:"summoner1Id"`
	Summoner2Casts                 int    `json:"summoner2Casts"`
	Summoner2Id                    int    `json:"summoner2Id"`
	SummonerId                     string `json:"summonerId"`
	SummonerLevel                  int    `json:"summonerLevel"`
	SummonerName                   string `json:"summonerName"`
	TeamEarlySurrendered           bool   `json:"teamEarlySurrendered"`
	TeamId                         int    `json:"teamId"`
	TeamPosition                   string `json:"teamPosition"`
	TimeCCingOthers                int    `json:"timeCCingOthers"`
	TimePlayed                     int    `json:"timePlayed"`
	TotalDamageDealt               int    `json:"totalDamageDealt"`
	TotalDamageDealtToChampions    int    `json:"totalDamageDealtToChampions"`
	TotalDamageShieldedOnTeammates int    `json:"totalDamageShieldedOnTeammates"`
	TotalDamageTaken               int    `json:"totalDamageTaken"`
	TotalHeal                      int    `json:"totalHeal"`
	TotalHealsOnTeammates          int    `json:"totalHealsOnTeammates"`
	TotalMinionsKilled             int    `json:"totalMinionsKilled"`
	TotalTimeCCDealt               int    `json:"totalTimeCCDealt"`
	TotalTimeSpentDead             int    `json:"totalTimeSpentDead"`
	TotalUnitsHealed               int    `json:"totalUnitsHealed"`
	TripleKills                    int    `json:"tripleKills"`
	TrueDamageDealt                int    `json:"trueDamageDealt"`
	TrueDamageDealtToChampions     int    `json:"trueDamageDealtToChampions"`
	TrueDamageTaken                int    `json:"trueDamageTaken"`
	TurretKills                    int    `json:"turretKills"`
	TurretTakedowns                int    `json:"turretTakedowns"`
	TurretsLost                    int    `json:"turretsLost"`
	UnrealKills                    int    `json:"unrealKills"`
	VisionScore                    int    `json:"visionScore"`
	VisionWardsBoughtInGame        int    `json:"visionWardsBoughtInGame"`
	WardsKilled                    int    `json:"wardsKilled"`
	WardsPlaced                    int    `json:"wardsPlaced"`
	Win                            bool   `json:"win"`
}

type Perks struct {
	StatPerks *PerkStats   `json:"statPerks"`
	Styles    []*PerkStyle `json:"styles"`
}

type PerkStats struct {
	Defense int `json:"defense"`
	Flex    int `json:"flex"`
	Offense int `json:"offense"`
}

type PerkStyle struct {
	Description string                `json:"description"`
	Selections  []*PerkStyleSelection `json:"selections"`
	Style       int                   `json:"style"`
}

type PerkStyleSelection struct {
	Perk int `json:"perk"`
	Var1 int `json:"var1"`
	Var2 int `json:"var2"`
	Var3 int `json:"var3"`
}

type Team struct {
	Bans       []*Ban      `json:"bans"`
	Objectives *Objectives `json:"objectives"`
	TeamId     int         `json:"teamId"`
	Win        bool        `json:"win"`
}

type Ban struct {
	ChampionId int `json:"championId"`
	PickTurn   int `json:"pickTurn"`
}

type Objectives struct {
	Baron      *Objective `json:"baron"`
	Champion   *Objective `json:"champion"`
	Dragon     *Objective `json:"dragon"`
	Horde      *Objective `json:"horde"`
	Inhibitor  *Objective `json:"inhibitor"`
	RiftHerald *Objective `json:"riftHerald"`
	Tower      *Objective `json:"tower"`
}

type Objective struct {
	First bool `json:"first"`
	Kills int  `json:"kills"`
}

// GetMatch returns a match by match id.
//
// Riot API docs: https://developer.riotgames.com/apis#match-v5/GET_getMatch
//
// GET /lol/match/v5/matches/{matchId}
func (c *Client) GetMatch(ctx context.Context, region, matchId string) (*Match, error) {
	u := regionHost(region)
	path := fmt.Sprintf("/lol/match/v5/matches/%s", matchId)
	u = u.JoinPath(path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Riot-Token", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	defer resp.Body.Close()
	var match Match
	err = json.NewDecoder(resp.Body).Decode(&match)
	return &match, err
}

// GetMatchIDsByPUUID returns a list of match ids by puuid.
//
// Riot API docs: https://developer.riotgames.com/apis#match-v5/GET_getMatchIdsByPUUID
//
// GET /lol/match/v5/matches/by-puuid/{puuid}/ids
func (c *Client) GetMatchIDsByPUUID(ctx context.Context, region, puuid, query string) (ids []string, err error) {
	u := regionHost(region)
	path := fmt.Sprintf("/lol/match/v5/matches/by-puuid/%s/ids", puuid)
	u = u.JoinPath(path)

	values, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}
	u.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Riot-Token", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&ids)
	return ids, err
}

type Timeline struct {
	Metadata MetadataTimeline `json:"metadata"`
	Info     InfoTimeline     `json:"info"`
}

type MetadataTimeline struct {
	DataVersion  string   `json:"dataVersion"`
	MatchID      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type InfoTimeline struct {
	EndOfGameResult string                `json:"endOfGameResult"`
	FrameInterval   int                   `json:"frameInterval"`
	GameID          int64                 `json:"gameId"`
	Participants    []ParticipantTimeline `json:"participants"`
	Frames          []FramesTimeline      `json:"frames"`
}

type ParticipantTimeline struct {
	ParticipantID int    `json:"participantId"`
	Puuid         string `json:"puuid"`
}

type FramesTimeline struct {
	Events            map[string]any           `json:"events"`
	ParticipantFrames map[int]ParticipantFrame `json:"participantFrames"`
	Timestamp         int                      `json:"timestamp"`
}

type EventsTimeline struct {
	RealTimestamp int64  `json:"realTimestamp"`
	Timestamp     int    `json:"timestamp"`
	Type          string `json:"type"`
}

type TLEventItemPurchased struct {
	ItemID        int    `json:"itemId"`
	ParticipantID int    `json:"participantId"`
	Timestamp     int    `json:"timestamp"`
	Type          string `json:"type"`
}

type TLEventSkillLevelUp struct {
	LevelUpType   string `json:"levelUpType"`
	ParticipantID int    `json:"participantId"`
	SkillSlot     int    `json:"skillSlot"`
	Timestamp     int    `json:"timestamp"`
	Type          string `json:"type"`
}

type ParticipantFrame struct {
	ChampionStats            ChampionStats `json:"championStats"`
	CurrentGold              int           `json:"currentGold"`
	GoldPerSecond            int           `json:"goldPerSecond"`
	JungleMinionsKilled      int           `json:"jungleMinionsKilled"`
	Level                    int           `json:"level"`
	MinionsKilled            int           `json:"minionsKilled"`
	ParticipantID            int           `json:"participantId"`
	Position                 Position      `json:"position"`
	DamageStats              DamageStats   `json:"damageStats"`
	TimeEnemySpentControlled int           `json:"timeEnemySpentControlled"`
	TotalGold                int           `json:"totalGold"`
	Xp                       int           `json:"xp"`
}

type ChampionStats struct {
	AbilityHaste         int `json:"abilityHaste"`
	AbilityPower         int `json:"abilityPower"`
	Armor                int `json:"armor"`
	ArmorPen             int `json:"armorPen"`
	ArmorPenPercent      int `json:"armorPenPercent"`
	AttackDamage         int `json:"attackDamage"`
	AttackSpeed          int `json:"attackSpeed"`
	BonusArmorPenPercent int `json:"bonusArmorPenPercent"`
	BonusMagicPenPercent int `json:"bonusMagicPenPercent"`
	CcReduction          int `json:"ccReduction"`
	CooldownReduction    int `json:"cooldownReduction"`
	Health               int `json:"health"`
	HealthMax            int `json:"healthMax"`
	HealthRegen          int `json:"healthRegen"`
	Lifesteal            int `json:"lifesteal"`
	MagicPen             int `json:"magicPen"`
	MagicPenPercent      int `json:"magicPenPercent"`
	MagicResist          int `json:"magicResist"`
	MovementSpeed        int `json:"movementSpeed"`
	Omnivamp             int `json:"omnivamp"`
	PhysicalVamp         int `json:"physicalVamp"`
	Power                int `json:"power"`
	PowerMax             int `json:"powerMax"`
	PowerRegen           int `json:"powerRegen"`
	SpellVamp            int `json:"spellVamp"`
}

type DamageStats struct {
	MagicDamageDone               int `json:"magicDamageDone"`
	MagicDamageDoneToChampions    int `json:"magicDamageDoneToChampions"`
	MagicDamageTaken              int `json:"magicDamageTaken"`
	PhysicalDamageDone            int `json:"physicalDamageDone"`
	PhysicalDamageDoneToChampions int `json:"physicalDamageDoneToChampions"`
	PhysicalDamageTaken           int `json:"physicalDamageTaken"`
	TotalDamageDone               int `json:"totalDamageDone"`
	TotalDamageDoneToChampions    int `json:"totalDamageDoneToChampions"`
	TotalDamageTaken              int `json:"totalDamageTaken"`
	TrueDamageDone                int `json:"trueDamageDone"`
	TrueDamageDoneToChampions     int `json:"trueDamageDoneToChampions"`
	TrueDamageTaken               int `json:"trueDamageTaken"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

const (
	EventBuildingKill         = "BUILDING_KILL"
	EventChampionKill         = "CHAMPION_KILL"
	EventChampionSpecialKill  = "CHAMPION_SPECIAL_KILL"
	EventGameEnd              = "GAME_END"
	EventItemDestroyed        = "ITEM_DESTROYED"
	EventItemPurchased        = "ITEM_PURCHASED"
	EventLevelUp              = "LEVEL_UP"
	EventPauseEnd             = "PAUSE_END"
	EventSkillLevelUp         = "SKILL_LEVEL_UP"
	EventTurretPlateDestroyed = "TURRET_PLATE_DESTROYED"
	EventWardKill             = "WARD_KILL"
	EventWardPlaced           = "WARD_PLACED"
)

// GetMatchTimeline returns the timeline of a match.
//
// Riot API docs: https://developer.riotgames.com/apis#match-v5/GET_getTimeline
func (c *Client) GetMatchTimeline(ctx context.Context, region, id string) (*Timeline, error) {
	u := regionHost(region)
	path := fmt.Sprintf("/lol/match/v5/matches/%s/timeline", id)
	u = u.JoinPath(path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Riot-Token", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	defer resp.Body.Close()
	var timeline Timeline
	err = json.NewDecoder(resp.Body).Decode(&timeline)
	return &timeline, err
}
