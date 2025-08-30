package riot

import (
	"context"
	"fmt"
	"net/url"

	"github.com/rank1zen/kevin/internal/riot/internal"
)

type MatchService service

// MatchListOptions are parameters to send in a [MatchService.GetMatchList]
// request.
type MatchListOptions struct {
	StartTime *int64

	EndTime *int64

	Queue *int

	Type *string

	Start int

	Count int
}

// MatchList is a list of match ids.
type MatchList []string

// GetMatchIDsByPUUID returns a list of match ids by puuid.
//
// Riot API docs: https://developer.riotgames.com/apis#match-v5/GET_getMatchIdsByPUUID
//
// GET /lol/match/v5/matches/by-puuid/{puuid}/ids
func (m *MatchService) GetMatchList(ctx context.Context, region Region, puuid string, opts MatchListOptions) (MatchList, error) {
	endpoint := fmt.Sprintf("/lol/match/v5/matches/by-puuid/%s/ids", puuid)

	query := url.Values{}

	query.Add("start", fmt.Sprintf("%d", opts.Start))
	query.Add("count", fmt.Sprintf("%d", opts.Count))

	if opts.StartTime != nil {
		query.Add("startTime", fmt.Sprintf("%d", *opts.StartTime))
	}
	if opts.EndTime != nil {
		query.Add("endTime", fmt.Sprintf("%d", *opts.EndTime))
	}
	if opts.Queue != nil {
		query.Add("queue", fmt.Sprintf("%d", *opts.Queue))
	}
	if opts.Type != nil {
		query.Add("type", *opts.Type)
	}

	req := &internal.Request{
		BaseURL:  region.continentHost(),
		Endpoint: endpoint,
		APIKey:   m.client.apiKey,
		Query:    query,
	}

	if m.client.baseURL != "" {
		req.BaseURL = m.client.baseURL
	}

	var ids MatchList
	if err := m.client.internals.DispatchRequest(ctx, req, &ids); err != nil {
		return nil, err
	}

	return ids, nil
}

type Match struct {
	Info     *MatchInfo     `json:"info"`
	Metadata *MatchMetadata `json:"metadata"`
}

type MatchMetadata struct {
	DataVersion  string   `json:"dataVersion"`
	MatchID      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type MatchInfo struct {
	EndOfGameResult    string              `json:"endOfGameResult"`
	GameCreation       int64               `json:"gameCreation"`
	GameDuration       int64               `json:"gameDuration"`
	GameEndTimestamp   int64               `json:"gameEndTimestamp"`
	GameID             int64               `json:"gameId"`
	GameMode           string              `json:"gameMode"`
	GameName           string              `json:"gameName"`
	GameStartTimestamp int64               `json:"gameStartTimestamp"`
	GameType           string              `json:"gameType"`
	GameVersion        string              `json:"gameVersion"`
	MapID              int                 `json:"mapId"`
	Participants       []*MatchParticipant `json:"participants"`
	PlatformID         string              `json:"platformId"`
	QueueID            int                 `json:"queueId"`
	Teams              []*MatchTeam        `json:"teams"`
	TournamentCode     string              `json:"tournamentCode"`
}

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
	RiotIdName                     string      `json:"riotIdName"`
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
	TeamID     int              `json:"teamId"`
	Win        bool             `json:"win"`
}

type MatchBan struct {
	ChampionID int `json:"championId"`
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

// GetMatch returns a match by match id.
//
// Riot API docs: https://developer.riotgames.com/apis#match-v5/GET_getMatch
//
// GET /lol/match/v5/matches/{matchId}
func (m *MatchService) GetMatch(ctx context.Context, region Region, matchID string) (*Match, error) {
	endpoint := fmt.Sprintf("/lol/match/v5/matches/%s", matchID)

	req := &internal.Request{
		BaseURL:  region.continentHost(),
		Endpoint: endpoint,
		APIKey:   m.client.apiKey,
	}

	if m.client.baseURL != "" {
		req.BaseURL = m.client.baseURL
	}

	var match Match
	if err := m.client.internals.DispatchRequest(ctx, req, &match); err != nil {
		return nil, err
	}

	// TODO: filter the matches that are played in region. Currently it
	// might return everything on the continent.

	return &match, nil
}

type TimelineEventType string

const (
	TimelineEventTypeAscendedEvent        TimelineEventType = "ASCENDED_EVENT"
	TimelineEventTypeBuildingKill         TimelineEventType = "BUILDING_KILL"
	TimelineEventTypeCapturePoint         TimelineEventType = "CAPTURE_POINT"
	TimelineEventTypeChampionKill         TimelineEventType = "CHAMPION_KILL"
	TimelineEventTypeChampionSpecialKill  TimelineEventType = "CHAMPION_SPECIAL_KILL"
	TimelineEventTypeEliteMonsterKill     TimelineEventType = "ELITE_MONSTER_KILL"
	TimelineEventTypeGameEnd              TimelineEventType = "GAME_END"
	TimelineEventTypeItemDestroyed        TimelineEventType = "ITEM_DESTROYED"
	TimelineEventTypeItemPurchased        TimelineEventType = "ITEM_PURCHASED"
	TimelineEventTypeItemSold             TimelineEventType = "ITEM_SOLD"
	TimelineEventTypeItemUndo             TimelineEventType = "ITEM_UNDO"
	TimelineEventTypeLevelUp              TimelineEventType = "LEVEL_UP"
	TimelineEventTypePauseEnd             TimelineEventType = "PAUSE_END"
	TimelineEventTypePoroKingSummon       TimelineEventType = "PORO_KING_SUMMON"
	TimelineEventTypeSkillLevelUp         TimelineEventType = "SKILL_LEVEL_UP"
	TimelineEventTypeTurretPlateDestroyed TimelineEventType = "TURRET_PLATE_DESTROYED"
	TimelineEventTypeWardKill             TimelineEventType = "WARD_KILL"
	TimelineEventTypeWardPlaced           TimelineEventType = "WARD_PLACED"
)

type Timeline struct {
	Metadata TimelineMetadata `json:"metadata"`
	Info     TimelineInfo     `json:"info"`
}

type TimelineMetadata struct {
	DataVersion  string   `json:"dataVersion"`
	MatchID      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type TimelineInfo struct {
	EndOfGameResult string                `json:"endOfGameResult"`
	FrameInterval   int                   `json:"frameInterval"`
	GameID          int64                 `json:"gameId"`
	Participants    []TimelineParticipant `json:"participants"`
	Frames          []TimelineFrames      `json:"frames"`
}

type TimelineParticipant struct {
	ParticipantID int    `json:"participantId"`
	PUUID         string `json:"puuid"`
}

type TimelineFrames struct {
	Events            map[string]any                   `json:"events"`
	ParticipantFrames map[int]TimelineParticipantFrame `json:"participantFrames"`
	Timestamp         int                              `json:"timestamp"`
}

type TimelineEvent struct {
	AfterID                 *int              `json:"afterId,omitempty"`
	AscendedType            *string           `json:"ascendedType,omitempty"`
	AssistingParticipantIDs *[]int            `json:"assistingParticipantIds,omitempty"`
	BeforeID                *int              `json:"beforeId,omitempty"`
	BuildingType            *string           `json:"buildingType,omitempty"`
	CreatorID               *int              `json:"creatorId,omitempty"`
	ItemID                  *int              `json:"itemId,omitempty"`
	KillerID                *int              `json:"killerId,omitempty"`
	LaneType                *string           `json:"laneType,omitempty"`
	LevelUpType             *string           `json:"levelUpType,omitempty"`
	MonsterSubType          *string           `json:"monsterSubType,omitempty"`
	MonsterType             *string           `json:"monsterType,omitempty"`
	ParticipantID           *int              `json:"participantId,omitempty"`
	PointCaptured           *string           `json:"pointCaptured,omitempty"`
	Position                *TimelinePosition `json:"position,omitempty"`
	RealTimestamp           int64             `json:"realTimestamp"`
	SkillSlot               *int              `json:"skillSlot,omitempty"`
	TeamID                  *int              `json:"teamId,omitempty"`
	Timestamp               int               `json:"timestamp"`
	TowerType               *string           `json:"towerType,omitempty"`
	Type                    TimelineEventType `json:"type"`
	VictimID                *int              `json:"victimId,omitempty"`
	WardType                *string           `json:"wardType,omitempty"`
}

type TimelineParticipantFrame struct {
	ChampionStats            TimelineChampionStats `json:"championStats"`
	CurrentGold              int                   `json:"currentGold"`
	GoldPerSecond            int                   `json:"goldPerSecond"`
	JungleMinionsKilled      int                   `json:"jungleMinionsKilled"`
	Level                    int                   `json:"level"`
	MinionsKilled            int                   `json:"minionsKilled"`
	ParticipantID            int                   `json:"participantId"`
	Position                 TimelinePosition      `json:"position"`
	DamageStats              TimelineDamageStats   `json:"damageStats"`
	TimeEnemySpentControlled int                   `json:"timeEnemySpentControlled"`
	TotalGold                int                   `json:"totalGold"`
	XP                       int                   `json:"xp"`
}

type TimelineChampionStats struct {
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

type TimelineDamageStats struct {
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

type TimelinePosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// GetMatchTimeline returns the timeline of a match.
//
// Riot API docs: https://developer.riotgames.com/apis#match-v5/GET_getTimeline
//
// GET /lol/match/v5/matches/{matchId}/timeline
func (m *MatchService) GetTimeline(ctx context.Context, region Region, id string) (*Timeline, error) {
	endpoint := fmt.Sprintf("/lol/match/v5/matches/%s/timeline", id)

	req := &internal.Request{
		BaseURL:  region.continentHost(),
		Endpoint: endpoint,
		APIKey:   m.client.apiKey,
	}

	var timeline Timeline
	if err := m.client.internals.DispatchRequest(ctx, req, &timeline); err != nil {
		return nil, err
	}

	return &timeline, nil
}
