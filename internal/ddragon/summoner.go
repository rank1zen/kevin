package ddragon

type Summoner struct {
	ID         ItemID
	Name       string
	FullIcon   string
	Sprite     string
	X, Y, W, H int
}

var (
	SummonerBarrier                 = Summoner{21, "Barrier", "SummonerBarrier.png", "spell0.png", 0, 0, 48, 48}
	SummonerBoost                   = Summoner{1, "Cleanse", "SummonerBoost.png", "spell0.png", 48, 0, 48, 48}
	SummonerCherryFlash             = Summoner{2202, "Flash", "SummonerCherryFlash.png", "spell0.png", 96, 0, 48, 48}
	SummonerCherryHold              = Summoner{2201, "Flee", "SummonerCherryHold.png", "spell0.png", 144, 0, 48, 48}
	SummonerDot                     = Summoner{14, "Ignite", "SummonerDot.png", "spell0.png", 192, 0, 48, 48}
	SummonerExhaust                 = Summoner{3, "Exhaust", "SummonerExhaust.png", "spell0.png", 240, 0, 48, 48}
	SummonerFlash                   = Summoner{4, "Flash", "SummonerFlash.png", "spell0.png", 288, 0, 48, 48}
	SummonerHaste                   = Summoner{6, "Ghost", "SummonerHaste.png", "spell0.png", 336, 0, 48, 48}
	SummonerHeal                    = Summoner{7, "Heal", "SummonerHeal.png", "spell0.png", 384, 0, 48, 48}
	SummonerMana                    = Summoner{13, "Clarity", "SummonerMana.png", "spell0.png", 432, 0, 48, 48}
	SummonerPoroRecall              = Summoner{30, "To the King!", "SummonerPoroRecall.png", "spell0.png", 0, 48, 48, 48}
	SummonerPoroThrow               = Summoner{31, "Poro Toss", "SummonerPoroThrow.png", "spell0.png", 48, 48, 48, 48}
	SummonerSmite                   = Summoner{11, "Smite", "SummonerSmite.png", "spell0.png", 96, 48, 48, 48}
	SummonerSnowURFSnowballMark     = Summoner{39, "Mark", "SummonerSnowURFSnowball_Mark.png", "spell0.png", 144, 48, 48, 48}
	SummonerSnowball                = Summoner{32, "Mark", "SummonerSnowball.png", "spell0.png", 192, 48, 48, 48}
	SummonerTeleport                = Summoner{12, "Teleport", "SummonerTeleport.png", "spell0.png", 240, 48, 48, 48}
	SummonerUltBookPlaceholder      = Summoner{54, "Placeholder", "Summoner_UltBookPlaceholder.png", "spell0.png", 288, 48, 48, 48}
	SummonerUltBookSmitePlaceholder = Summoner{55, "Placeholder and Attack-Smite", "Summoner_UltBookSmitePlaceholder.png", "spell0.png", 336, 48, 48, 48}
)
