package ladder

// Leaderboard is the top 500 players in a region
type Leaderboard struct {
	Region string

	Entries []LeaderboardEntry
}

// LeaderboardEntry is an entry for a player in a leaderboard.
type LeaderboardEntry struct {
	Region string

	PUUID string

	Name, Tag string

	Tier string

	Division string

	LP int

	Wins, Losses int
}

type LeaderboardFilter struct {
	Start int
	Count int
}
