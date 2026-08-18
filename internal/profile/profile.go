package profile

// Profile is a summoner's profile. PUUID is unique and immutable. Name +
// Tagline is unique but mutable.
type Profile struct {
	PUUID         string
	Name, Tagline string

	// Rank is the summoners most recent rank record in store, meaning it could be
	// out-of-date.
	Rank RankStatus
}
