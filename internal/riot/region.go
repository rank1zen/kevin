package riot

type Region string

func (r Region) host() string {
	host, ok := regionToHost[r]
	if !ok {
		panic("region is not valid")
	}

	return host
}

func (r Region) continentHost() string {
	c := RegionToContinent(r)
	host, ok := continentToHost[c]
	if !ok {
		panic("region is not valid")
	}

	return host
}

const (
	RegionBR1 Region = "BR1"
	RegionEUN1 Region = "EUN1"
	RegionEUW1 Region = "EUW1"
	RegionJP1 Region = "JP1"
	RegionKR Region = "KR"
	RegionLA1 Region = "LA1"
	RegionLA2 Region = "LA2"
	RegionME1 Region = "ME1"
	RegionNA1 Region = "NA1"
	RegionOC1 Region = "OC1"
	RegionRU Region = "RU"
	RegionSEA Region = "SG2"
	RegionTR1 Region = "TR1"
	RegionTW2 Region = "TW2"
	RegionVN2 Region = "VN2"
)

var regionToHost = map[Region]string{
	RegionBR1:  "https://br1.api.riotgames.com",
	RegionEUN1: "https://eun1.api.riotgames.com",
	RegionEUW1: "https://euw1.api.riotgames.com",
	RegionJP1:  "https://jp1.api.riotgames.com",
	RegionKR:   "https://kr.api.riotgames.com",
	RegionLA1:  "https://la1.api.riotgames.com",
	RegionLA2:  "https://la2.api.riotgames.com",
	RegionME1:  "https://me1.api.riotgames.com",
	RegionNA1:  "https://na1.api.riotgames.com",
	RegionOC1:  "https://oc1.api.riotgames.com",
	RegionRU:   "https://ru.api.riotgames.com",
	RegionSEA:  "https://sg2.api.riotgames.com",
	RegionTR1:  "https://tr1.api.riotgames.com",
	RegionTW2:  "https://tw2.api.riotgames.com",
	RegionVN2:  "https://vn2.api.riotgames.com",
}

type Continent string

const (
	ContinentAmericas = "AMERICAS"
	ContinentAsia     = "ASIA"
	ContinentEurope   = "EUROPE"
	ContinentSea      = "SEA"
)

var continentToHost = map[Continent]string{
	ContinentAmericas: "https://americas.api.riotgames.com",
	ContinentAsia    : "https://asia.api.riotgames.com",
	ContinentEurope  : "https://europe.api.riotgames.com",
	ContinentSea     : "https://sea.api.riotgames.com",
}

func RegionToContinent(region Region) Continent {
	switch region {
	case RegionNA1:
		return ContinentAmericas
	default:
		return ContinentAmericas
	}
}
