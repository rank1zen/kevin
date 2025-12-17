package riot

type Region string

func (r Region) host() string {
	host, ok := regionToHost[r]
	if !ok {
		panic("region is not valid")
	}

	return host
}

func regionHost(region string) string {
	host, ok := regionStrToHost[region]
	if !ok {
		return regionStrToHost["NA1"]
	}

	return host
}

func (r Region) continentHost() string {
	c := getRegionToMatchContinent(r)

	host, ok := continentToHost[c]
	if !ok {
		panic("region is not valid")
	}

	return host
}

const (
	RegionBR1  Region = "BR1"
	RegionEUN1 Region = "EUN1"
	RegionEUW1 Region = "EUW1"
	RegionJP1  Region = "JP1"
	RegionKR   Region = "KR"
	RegionLA1  Region = "LA1"
	RegionLA2  Region = "LA2"
	RegionME1  Region = "ME1"
	RegionNA1  Region = "NA1"
	RegionOC1  Region = "OC1"
	RegionRU   Region = "RU"
	RegionSEA  Region = "SG2"
	RegionTR1  Region = "TR1"
	RegionTW2  Region = "TW2"
	RegionVN2  Region = "VN2"
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

var regionStrToHost = map[string]string{
	"BR1":  "https://br1.api.riotgames.com",
	"EUN1": "https://eun1.api.riotgames.com",
	"EUW1": "https://euw1.api.riotgames.com",
	"JP1":  "https://jp1.api.riotgames.com",
	"KR":   "https://kr.api.riotgames.com",
	"LA1":  "https://la1.api.riotgames.com",
	"LA2":  "https://la2.api.riotgames.com",
	"ME1":  "https://me1.api.riotgames.com",
	"NA1":  "https://na1.api.riotgames.com",
	"OC1":  "https://oc1.api.riotgames.com",
	"RU":   "https://ru.api.riotgames.com",
	"SEA":  "https://sg2.api.riotgames.com",
	"TR1":  "https://tr1.api.riotgames.com",
	"TW2":  "https://tw2.api.riotgames.com",
	"VN2":  "https://vn2.api.riotgames.com",
}

type continent string

const (
	continentAmericas = "AMERICAS"
	continentAsia     = "ASIA"
	continentEurope   = "EUROPE"
	continentSEA      = "SEA"
)

var continentToHost = map[continent]string{
	continentAmericas: "https://americas.api.riotgames.com",
	continentAsia:     "https://asia.api.riotgames.com",
	continentEurope:   "https://europe.api.riotgames.com",
	continentSEA:      "https://sea.api.riotgames.com",
}

var regionToMatchContinent = map[Region]continent{
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

// regionToMatchContinent is used for the match-v5 api.
//
// The AMERICAS routing value serves NA, BR, LAN and LAS. The ASIA routing
// value serves KR and JP. The EUROPE routing value serves EUNE, EUW, ME1, TR
// and RU. The SEA routing value serves OCE, SG2, TW2 and VN2.
func getRegionToMatchContinent(region Region) continent {
	c, ok := regionToMatchContinent[region]
	if !ok {
		return continentAmericas
	}

	return c
}
