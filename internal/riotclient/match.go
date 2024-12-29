package riotclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Get a list of match ids by puuid (ONLY RANKED SOLOQ, queueId 420)
// TODO: make this more general and have options to specify ranked games only
func (s *Client) GetMatchIdsByPuuid(ctx context.Context, puuid string, start int, count int) ([]string, error) {
	u := fmt.Sprintf(defaultAmerBaseURL+"/lol/match/v5/matches/by-puuid/%s/ids", puuid)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	soloqQueueID := "420"
	q := req.URL.Query()
	q.Add("queue", soloqQueueID)
	q.Add("start", strconv.Itoa(start))
	q.Add("count", strconv.Itoa(count))
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Accept", "application/json")
	body, err := execute(ctx, req)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var ids []string
	err = json.NewDecoder(body).Decode(&ids)
	if err != nil {
		return nil, fmt.Errorf("riot: json error (%v)", err)
	}

	return ids, nil
}

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
	// EndOfGameResult indicates if the game ended in termination.
	EndOfGameResult string `json:"endOfGameResult"`

	// GameCreation is a Unix timestamp, in milliseconds, incidating the
	// time of game creation on the server, i.e., the loading screen.
	GameCreation int64 `json:"gameCreation"`

	// GameDuration is the duration of a game.
	// Prior to patch 11.20, this field returns the game length in
	// milliseconds calculated from gameEndTimestamp - gameStartTimestamp.
	// Post patch 11.20, this field returns the max timePlayed of any
	// participant in the game in seconds, which makes the behavior of this
	// field consistent with that of match-v4. The best way to handling the
	// change in this field is to treat the value as milliseconds if the
	// gameEndTimestamp field isn't in the response and to treat the value
	// as seconds if gameEndTimestamp is in the response.
	GameDuration int64 `json:"gameDuration"`

	// GameEndTimestamp is a Unix timestamp for when match ends on the game
	// server. This timestamp can occasionally be significantly longer than
	// when the match "ends". The most reliable way of determining the
	// timestamp for the end of the match would be to add the max time
	// played of any participant to the gameStartTimestamp. This field was
	// added to match-v5 in patch 11.20 on Oct 5th, 2021.
	GameEndTimestamp int64 `json:"gameEndTimestamp"`

	GameId int64 `json:"gameId"`

	// Refer to the Game Constants Documentation
	GameMode string `json:"gameMode"`

	GameName string `json:"gameName"`

	// GameStartTimestamp is a Unix timestamp, in milliseconds, indicating the
	// time of game start on the server.
	GameStartTimestamp int64 `json:"gameStartTimestamp"`

	GameType string `json:"gameType"`

	GameVersion string `json:"gameVersion"`

	// Refer to the Game Constants Documentation
	MapId int `json:"mapId"`

	Participants []*MatchParticipant `json:"participants"`

	PlatformId string `json:"platformId"`

	// Refer to the Game Constants Documentation
	QueueId int `json:"queueId"`

	Teams []*MatchTeam `json:"teams"`

	TournamentCode string `json:"tournamentCode"`
}

// NOTE: there are a lot more fields here
type MatchParticipant struct {
	Assists                        int         `json:"assists"`
	BaronKills                     int         `json:"baronKills"`
	BountyLevel                    int         `json:"bountyLevel"`
	ChampExperience                int         `json:"champExperience"`
	ChampLevel                     int         `json:"champLevel"`
	ChampionID                     int         `json:"championId"`
	ChampionName                   string      `json:"championName"`
	ChampionTransform              int         `json:"championTransform"`
	ConsumablesPurchased           int         `json:"consumablesPurchased"`
	DamageDealtToBuildings         int         `json:"damageDealtToBuildings"`
	DamageDealtToObjectives        int         `json:"damageDealtToObjectives"`
	DamageDealtToTurrets           int         `json:"damageDealtToTurrets"`
	DamageSelfMitigated            int         `json:"damageSelfMitigated"`
	Deaths                         int         `json:"deaths"`
	DetectorWardsPlaced            int         `json:"detectorWardsPlaced"`
	DoubleKills                    int         `json:"doubleKills"`
	DragonKills                    int         `json:"dragonKills"`
	FirstBloodAssist               bool        `json:"firstBloodAssist"`
	FirstBloodKill                 bool        `json:"firstBloodKill"`
	FirstTowerAssist               bool        `json:"firstTowerAssist"`
	FirstTowerKill                 bool        `json:"firstTowerKill"`
	GameEndedInEarlySurrender      bool        `json:"gameEndedInEarlySurrender"`
	GameEndedInSurrender           bool        `json:"gameEndedInSurrender"`
	GoldEarned                     int         `json:"goldEarned"`
	GoldSpent                      int         `json:"goldSpent"`
	IndividualPosition             string      `json:"individualPosition"`
	InhibitorKills                 int         `json:"inhibitorKills"`
	InhibitorTakedowns             int         `json:"inhibitorTakedowns"`
	InhibitorsLost                 int         `json:"inhibitorsLost"`
	Item0                          int         `json:"item0"`
	Item1                          int         `json:"item1"`
	Item2                          int         `json:"item2"`
	Item3                          int         `json:"item3"`
	Item4                          int         `json:"item4"`
	Item5                          int         `json:"item5"`
	Item6                          int         `json:"item6"`
	ItemsPurchased                 int         `json:"itemsPurchased"`
	KillingSprees                  int         `json:"killingSprees"`
	Kills                          int         `json:"kills"`
	Lane                           string      `json:"lane"`
	LargestCriticalStrike          int         `json:"largestCriticalStrike"`
	LargestKillingSpree            int         `json:"largestKillingSpree"`
	LargestMultiKill               int         `json:"largestMultiKill"`
	LongestTimeSpentLiving         int         `json:"longestTimeSpentLiving"`
	MagicDamageDealt               int         `json:"magicDamageDealt"`
	MagicDamageDealtToChampions    int         `json:"magicDamageDealtToChampions"`
	MagicDamageTaken               int         `json:"magicDamageTaken"`
	NeutralMinionsKilled           int         `json:"neutralMinionsKilled"`
	NexusKills                     int         `json:"nexusKills"`
	NexusLost                      int         `json:"nexusLost"`
	NexusTakedowns                 int         `json:"nexusTakedowns"`
	ObjectivesStolen               int         `json:"objectivesStolen"`
	ObjectivesStolenAssists        int         `json:"objectivesStolenAssists"`
	ParticipantID                  int         `json:"participantId"`
	PentaKills                     int         `json:"pentaKills"`
	Perks                          *MatchPerks `json:"perks"`
	PhysicalDamageDealt            int         `json:"physicalDamageDealt"`
	PhysicalDamageDealtToChampions int         `json:"physicalDamageDealtToChampions"`
	PhysicalDamageTaken            int         `json:"physicalDamageTaken"`
	ProfileIcon                    int         `json:"profileIcon"`
	PUUID                          string      `json:"puuid"`
	QuadraKills                    int         `json:"quadraKills"`
	RiotIDGameName                 string      `json:"riotIdGameName"`
	RiotIDName                     string      `json:"riotIdName"`
	RiotIDTagline                  string      `json:"riotIdTagline"`
	Role                           string      `json:"role"`
	SightWardsBoughtInGame         int         `json:"sightWardsBoughtInGame"`
	Spell1Casts                    int         `json:"spell1Casts"`
	Spell2Casts                    int         `json:"spell2Casts"`
	Spell3Casts                    int         `json:"spell3Casts"`
	Spell4Casts                    int         `json:"spell4Casts"`
	Summoner1Casts                 int         `json:"summoner1Casts"`
	Summoner1ID                    int         `json:"summoner1Id"`
	Summoner2Casts                 int         `json:"summoner2Casts"`
	Summoner2ID                    int         `json:"summoner2Id"`
	SummonerID                     string      `json:"summonerId"`
	SummonerLevel                  int         `json:"summonerLevel"`
	SummonerName                   string      `json:"summonerName"`
	TeamEarlySurrendered           bool        `json:"teamEarlySurrendered"`
	TeamID                         int         `json:"teamId"`
	TeamPosition                   string      `json:"teamPosition"`
	TimeCCingOthers                int         `json:"timeCCingOthers"`
	TimePlayed                     int         `json:"timePlayed"`
	TotalDamageDealt               int         `json:"totalDamageDealt"`
	TotalDamageDealtToChampions    int         `json:"totalDamageDealtToChampions"`
	TotalDamageShieldedOnTeammates int         `json:"totalDamageShieldedOnTeammates"`
	TotalDamageTaken               int         `json:"totalDamageTaken"`
	TotalHeal                      int         `json:"totalHeal"`
	TotalHealsOnTeammates          int         `json:"totalHealsOnTeammates"`
	TotalMinionsKilled             int         `json:"totalMinionsKilled"`
	TotalTimeCCDealt               int         `json:"totalTimeCCDealt"`
	TotalTimeSpentDead             int         `json:"totalTimeSpentDead"`
	TotalUnitsHealed               int         `json:"totalUnitsHealed"`
	TripleKills                    int         `json:"tripleKills"`
	TrueDamageDealt                int         `json:"trueDamageDealt"`
	TrueDamageDealtToChampions     int         `json:"trueDamageDealtToChampions"`
	TrueDamageTaken                int         `json:"trueDamageTaken"`
	TurretKills                    int         `json:"turretKills"`
	TurretTakedowns                int         `json:"turretTakedowns"`
	TurretsLost                    int         `json:"turretsLost"`
	UnrealKills                    int         `json:"unrealKills"`
	VisionScore                    int         `json:"visionScore"`
	VisionWardsBoughtInGame        int         `json:"visionWardsBoughtInGame"`
	WardsKilled                    int         `json:"wardsKilled"`
	WardsPlaced                    int         `json:"wardsPlaced"`
	Win                            bool        `json:"win"`
}

type MatchPerks struct {
	StatPerks *MatchPerkStats   `json:"statPerks"`
	Styles    []*MatchPerkStyle `json:"styles"`
}

type MatchPerkStats struct {
	Defense int `json:"defense"`
	Flex    int `json:"flex"`
	Offense int `json:"offense"`
}

type MatchPerkStyle struct {
	Description string                     `json:"description"`
	Selections  []*MatchPerkStyleSelection `json:"selections"`
	Style       int                        `json:"style"`
}

type MatchPerkStyleSelection struct {
	Perk int `json:"perk"`
	Var1 int `json:"var1"`
	Var2 int `json:"var2"`
	Var3 int `json:"var3"`
}

type MatchTeam struct {
	Bans       []*MatchBan      `json:"bans"`
	Objectives *MatchObjectives `json:"objectives"`
	TeamId     int              `json:"teamId"`
	Win        bool             `json:"win"`
}

type MatchBan struct {
	ChampionId int `json:"championId"`
	PickTurn   int `json:"pickTurn"`
}

type MatchObjectives struct {
	Baron      *MatchObjective `json:"baron"`
	Champion   *MatchObjective `json:"champion"`
	Dragon     *MatchObjective `json:"dragon"`
	Horde      *MatchObjective `json:"horde"`
	Inhibitor  *MatchObjective `json:"inhibitor"`
	RiftHerald *MatchObjective `json:"riftHerald"`
	Tower      *MatchObjective `json:"tower"`
}

type MatchObjective struct {
	First bool `json:"first"`
	Kills int  `json:"kills"`
}

// Get a match by match id
func (c *Client) GetMatch(ctx context.Context, matchID string) (*Match, error) {
	u := fmt.Sprintf(defaultAmerBaseURL+"/lol/match/v5/matches/%s", matchID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	body, err := execute(ctx, req)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var m Match
	err = json.NewDecoder(body).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("riot: json error (%v)", err)
	}

	return &m, nil
}

// TODO: these are todos

type MatchTimeline struct {
	Metadata *MetadataTimeLine
	Info     *InfoTimeLine
}

type MetadataTimeLine struct{}

type InfoTimeLine struct {
	EndOfGameResult string
	FrameInterval   int64
	GameId          int64
	Participants    []*ParticipantTimeLine
	Frames          []*FramesTimeLine
}

type ParticipantTimeLine struct{}

type FramesTimeLine struct {
	ParticipantFrames map[string]*ParticipantFrame `json:"participantFrames"`
	Events            []*MatchEvent                `json:"events"`
	Timestamp         int                          `json:"timestamp"`
}

type MatchPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// ParticipantFrame contains information about a participant in a game at a single timestamp
type ParticipantFrame struct {
	Position            *MatchPosition `json:"position"`
	TotalGold           int            `json:"totalGold"`
	TeamScore           int            `json:"teamScore"`
	ParticipantID       int            `json:"participantId"`
	Level               int            `json:"level"`
	CurrentGold         int            `json:"currentGold"`
	MinionsKilled       int            `json:"minionsKilled"`
	DominionScore       int            `json:"dominionScore"`
	XP                  int            `json:"xp"`
	JungleMinionsKilled int            `json:"jungleMinionsKilled"`
}

// MatchEventType is the type of an event
type MatchEventType string

// All legal value for match event types
const (
	MatchEventTypeChampionKill     MatchEventType = "CHAMPION_KILL"
	MatchEventTypeWardPlaced       MatchEventType = "WARD_PLACED"
	MatchEventTypeWardKill         MatchEventType = "WARD_KILL"
	MatchEventTypeBuildingKill     MatchEventType = "BUILDING_KILL"
	MatchEventTypeEliteMonsterKill MatchEventType = "ELITE_MONSTER_KILL"
	MatchEventTypeItemPurchased    MatchEventType = "ITEM_PURCHASED"
	MatchEventTypeItemSold         MatchEventType = "ITEM_SOLD"
	MatchEventTypeItemDestroyed    MatchEventType = "ITEM_DESTROYED"
	MatchEventTypeItemUndo         MatchEventType = "ITEM_UNDO"
	MatchEventTypeSkillLevelUp     MatchEventType = "SKILL_LEVEL_UP"
	MatchEventTypeAscendedEvent    MatchEventType = "ASCENDED_EVENT"
	MatchEventTypeCapturePoint     MatchEventType = "CAPTURE_POINT"
	MatchEventTypePoroKingSummon   MatchEventType = "PORO_KING_SUMMON"
)

// MatchEventTypes is a list of all available match events
var MatchEventTypes = []MatchEventType{
	MatchEventTypeChampionKill,
	MatchEventTypeWardPlaced,
	MatchEventTypeWardKill,
	MatchEventTypeBuildingKill,
	MatchEventTypeEliteMonsterKill,
	MatchEventTypeItemPurchased,
	MatchEventTypeItemSold,
	MatchEventTypeItemDestroyed,
	MatchEventTypeItemUndo,
	MatchEventTypeSkillLevelUp,
	MatchEventTypeAscendedEvent,
	MatchEventTypeCapturePoint,
	MatchEventTypePoroKingSummon,
}

// MatchEvent is an event in a match at a certain timestamp
type MatchEvent struct {
	Type                    *MatchEventType `json:"type"`
	Position                *MatchPosition  `json:"position"`
	LevelUpType             string          `json:"levelUpType"`
	AscendedType            string          `json:"ascendedType"`
	TowerType               string          `json:"towerType"`
	EventType               string          `json:"eventType"`
	PointCaptured           string          `json:"pointCaptured"`
	WardType                string          `json:"wardType"`
	MonsterType             string          `json:"monsterType"`
	BuildingType            string          `json:"buildingType"`
	LaneType                string          `json:"laneType"`
	MonsterSubType          string          `json:"monsterSubType"`
	AssistingParticipantIDs []int           `json:"assistingParticipantIds"`
	Timestamp               int             `json:"timestamp"`
	AfterID                 int             `json:"afterId"`
	VictimID                int             `json:"victimId"`
	SkillSlot               int             `json:"skillSlot"`
	ItemID                  int             `json:"itemId"`
	ParticipantID           int             `json:"participantId"`
	TeamID                  int             `json:"teamId"`
	CreatorID               int             `json:"creatorId"`
	KillerID                int             `json:"killerId"`
	BeforeID                int             `json:"beforeId"`
}
