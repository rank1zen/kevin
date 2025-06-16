package riot

import "net/url"

const (
	RegionAmericas = "AMERICAS"
	RegionAsia     = "ASIA"
	RegionEurope   = "EUROPE"
	RegionSea      = "SEA"
)

func regionHost(region string) *url.URL {
	switch region {
	case RegionAmericas:
		u, _ := url.ParseRequestURI("https://americas.api.riotgames.com")
		return u
	case RegionAsia:
		u, _ := url.ParseRequestURI("https://asia.api.riotgames.com")
		return u
	case RegionEurope:
		u, _ := url.ParseRequestURI("https://europe.api.riotgames.com")
		return u
	case RegionSea:
		u, _ := url.ParseRequestURI("https://sea.api.riotgames.com")
		return u
	default:
		u, _ := url.ParseRequestURI("https://americas.api.riotgames.com")
		return u
	}
}
