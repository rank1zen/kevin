package riot

import (
	"net/url"
)

const (
	PlatformBR1  = "BR1"
	PlatformEUN1 = "EUN1"
	PlatformEUW1 = "EUW1"
	PlatformJP1  = "JP1"
	PlatformKR   = "KR"
	PlatformLA1  = "LA1"
	PlatformLA2  = "LA2"
	PlatformNA1  = "NA1"
	PlatformOC1  = "OC1"
	PlatformTR1  = "TR1"
	PlatformRU   = "RU"
	PlatformPH2  = "PH2"
	PlatformSG2  = "SG2"
	PlatformTH2  = "TH2"
	PlatformTW2  = "TW2"
	PlatformVN2  = "VN2"
)

func platformHost(platform string) *url.URL {
	// TODO: fill the rest in
	switch platform {
	case PlatformBR1:
		u, _ := url.ParseRequestURI("https://br1.api.riotgames.com")
		return u
	case PlatformNA1:
		u, _ := url.ParseRequestURI("https://na1.api.riotgames.com")
		return u
	default:
		u, _ := url.ParseRequestURI("https://na1.api.riotgames.com")
		return u
	}
}

func PlatformToRegion(platform string) string {
	switch platform {
	case PlatformNA1:
		return RegionAmericas
	default:
		return RegionAmericas
	}
}
