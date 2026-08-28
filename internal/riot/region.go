package riot

const (
	RegionBR1  = "BR1"
	RegionEUN1 = "EUN1"
	RegionEUW1 = "EUW1"
	RegionJP1  = "JP1"
	RegionKR   = "KR"
	RegionLA1  = "LA1"
	RegionLA2  = "LA2"
	RegionME1  = "ME1"
	RegionNA1  = "NA1"
	RegionOC1  = "OC1"
	RegionRU   = "RU"
	RegionSEA  = "SG2"
	RegionTR1  = "TR1"
	RegionTW2  = "TW2"
	RegionVN2  = "VN2"
)

var regionToHost = map[string]string{
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

const (
	continentAmericas = "AMERICAS"
	continentAsia     = "ASIA"
	continentEurope   = "EUROPE"
	continentSEA      = "SEA"
)

var continentToHost = map[string]string{
	continentAmericas: "https://americas.api.riotgames.com",
	continentAsia:     "https://asia.api.riotgames.com",
	continentEurope:   "https://europe.api.riotgames.com",
	continentSEA:      "https://sea.api.riotgames.com",
}

var regionToContinent = map[string]string{
	RegionBR1:  continentAmericas,
	RegionEUN1: continentEurope,
	RegionEUW1: continentEurope,
	RegionJP1:  continentAsia,
	RegionKR:   continentAsia,
	RegionLA1:  continentAmericas,
	RegionLA2:  continentAmericas,
	RegionME1:  continentEurope,
	RegionNA1:  continentAmericas,
	RegionOC1:  continentSEA,
	RegionRU:   continentEurope,
	RegionSEA:  continentSEA,
	RegionTR1:  continentEurope,
	RegionTW2:  continentSEA,
	RegionVN2:  continentSEA,
}

func regionHost(region string) string {
	host, ok := regionToHost[region]
	if !ok {
		return regionToHost["NA1"]
	}

	return host
}

func continentHost(region string) string {
	continent := getRegionToMatchContinent(region)

	host, ok := continentToHost[continent]
	if !ok {
		return continentToHost["AMERICAS"]
	}

	return host
}

// regionToContinent is used for the match-v5 api.
//
// The AMERICAS routing value serves NA, BR, LAN and LAS. The ASIA routing
// value serves KR and JP. The EUROPE routing value serves EUNE, EUW, ME1, TR
// and RU. The SEA routing value serves OCE, SG2, TW2 and VN2.
func getRegionToMatchContinent(region string) string {
	c, ok := regionToContinent[region]
	if !ok {
		return continentAmericas
	}

	return c
}
