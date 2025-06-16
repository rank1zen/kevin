package internal

import (
	"github.com/rank1zen/kevin/internal/riot"
)

type Rune int

func GetRuneIcon(id Rune) (u string) {
	return runeIcon[id]
}

type RunePage struct {
	PrimaryTree     Rune
	PrimaryKeystone Rune
	PrimaryA        Rune
	PrimaryB        Rune
	PrimaryC        Rune
	SecondaryTree   Rune
	SecondaryA      Rune
	SecondaryB      Rune
	MiniOffense     Rune
	MiniFlex        Rune
	MiniDefense     Rune
}

func NewRunePage(opts ...func(*RunePage)) (runes RunePage) {
	for _, f := range opts {
		f(&runes)
	}
	return runes
}

func RiotPerksToRunePage(perks *riot.Perks) func(*RunePage) {
	return func(runes *RunePage) {
		runes.PrimaryTree = Rune(perks.Styles[0].Style)
		runes.PrimaryKeystone = Rune(perks.Styles[0].Selections[0].Perk)
		runes.PrimaryA = Rune(perks.Styles[0].Selections[1].Perk)
		runes.PrimaryB = Rune(perks.Styles[0].Selections[2].Perk)
		runes.PrimaryC = Rune(perks.Styles[0].Selections[3].Perk)
		runes.SecondaryTree = Rune(perks.Styles[1].Style)
		runes.SecondaryA = Rune(perks.Styles[1].Selections[0].Perk)
		runes.SecondaryB = Rune(perks.Styles[1].Selections[1].Perk)
		runes.MiniOffense = Rune(perks.StatPerks.Offense)
		runes.MiniFlex = Rune(perks.StatPerks.Flex)
		runes.MiniDefense = Rune(perks.StatPerks.Defense)
	}
}

func makeRunePageFromRiot(perks *riot.Perks) RunePage {
	return RunePage{
		PrimaryTree:     Rune(perks.Styles[0].Style),
		PrimaryKeystone: Rune(perks.Styles[0].Selections[0].Perk),
		PrimaryA:        Rune(perks.Styles[0].Selections[1].Perk),
		PrimaryB:        Rune(perks.Styles[0].Selections[2].Perk),
		PrimaryC:        Rune(perks.Styles[0].Selections[3].Perk),
		SecondaryTree:   Rune(perks.Styles[1].Style),
		SecondaryA:      Rune(perks.Styles[1].Selections[0].Perk),
		SecondaryB:      Rune(perks.Styles[1].Selections[1].Perk),
		MiniOffense:     Rune(perks.StatPerks.Offense),
		MiniFlex:        Rune(perks.StatPerks.Flex),
		MiniDefense:     Rune(perks.StatPerks.Defense),
	}
}

func makeRunePage(runes [11]Rune) RunePage {
	return RunePage{
		PrimaryTree:     runes[0],
		PrimaryKeystone: runes[1],
		PrimaryA:        runes[2],
		PrimaryB:        runes[3],
		PrimaryC:        runes[4],
		SecondaryTree:   runes[5],
		SecondaryA:      runes[6],
		SecondaryB:      runes[7],
		MiniOffense:     runes[8],
		MiniFlex:        runes[9],
		MiniDefense:     runes[10],
	}
}

var runeIcon = map[Rune]string{
	8100: "perk-images/Styles/7200_Domination.png",
	8112: "perk-images/Styles/Domination/Electrocute/Electrocute.png",
	8128: "perk-images/Styles/Domination/DarkHarvest/DarkHarvest.png",
	9923: "perk-images/Styles/Domination/HailOfBlades/HailOfBlades.png",
	8126: "perk-images/Styles/Domination/CheapShot/CheapShot.png",
	8139: "perk-images/Styles/Domination/TasteOfBlood/GreenTerror_TasteOfBlood.png",
	8143: "perk-images/Styles/Domination/SuddenImpact/SuddenImpact.png",
	8136: "perk-images/Styles/Domination/ZombieWard/ZombieWard.png",
	8120: "perk-images/Styles/Domination/GhostPoro/GhostPoro.png",
	8138: "perk-images/Styles/Domination/EyeballCollection/EyeballCollection.png",
	8135: "perk-images/Styles/Domination/TreasureHunter/TreasureHunter.png",
	8105: "perk-images/Styles/Domination/RelentlessHunter/RelentlessHunter.png",
	8106: "perk-images/Styles/Domination/UltimateHunter/UltimateHunter.png",

	8300: "perk-images/Styles/7203_Whimsy.png",
	8351: "perk-images/Styles/Inspiration/GlacialAugment/GlacialAugment.png",
	8360: "perk-images/Styles/Inspiration/UnsealedSpellbook/UnsealedSpellbook.png",
	8369: "perk-images/Styles/Inspiration/FirstStrike/FirstStrike.png",
	8306: "perk-images/Styles/Inspiration/HextechFlashtraption/HextechFlashtraption.png",
	8304: "perk-images/Styles/Inspiration/MagicalFootwear/MagicalFootwear.png",
	8321: "perk-images/Styles/Inspiration/CashBack/CashBack2.png",
	8313: "perk-images/Styles/Inspiration/PerfectTiming/AlchemistCabinet.png",
	8352: "perk-images/Styles/Inspiration/TimeWarpTonic/TimeWarpTonic.png",
	8345: "perk-images/Styles/Inspiration/BiscuitDelivery/BiscuitDelivery.png",
	8347: "perk-images/Styles/Inspiration/CosmicInsight/CosmicInsight.png",
	8316: "perk-images/Styles/Inspiration/JackOfAllTrades/JackofAllTrades2.png",

	8000: "perk-images/Styles/7201_Precision.png",
	8005: "perk-images/Styles/Precision/PressTheAttack/PressTheAttack.png",
	8008: "perk-images/Styles/Precision/LethalTempo/LethalTempoTemp.png",
	8021: "perk-images/Styles/Precision/FleetFootwork/FleetFootwork.png",
	8010: "perk-images/Styles/Precision/Conqueror/Conqueror.png",
	9101: "perk-images/Styles/Precision/AbsorbLife/AbsorbLife.png",
	9111: "perk-images/Styles/Precision/Triumph.png",
	8009: "perk-images/Styles/Precision/PresenceOfMind/PresenceOfMind.png",
	9104: "perk-images/Styles/Precision/LegendAlacrity/LegendAlacrity.png",
	9105: "perk-images/Styles/Precision/LegendHaste/LegendHaste.png",
	9103: "perk-images/Styles/Precision/LegendBloodline/LegendBloodline.png",
	8014: "perk-images/Styles/Precision/CoupDeGrace/CoupDeGrace.png",
	8017: "perk-images/Styles/Precision/CutDown/CutDown.png",

	8400: "perk-images/Styles/7204_Resolve.png",
	8410: "perk-images/Styles/Resolve/ApproachVelocity/ApproachVelocity.png",
	8437: "perk-images/Styles/Resolve/GraspOfTheUndying/GraspOfTheUndying.png",
	8439: "perk-images/Styles/Resolve/VeteranAftershock/VeteranAftershock.png",
	8465: "perk-images/Styles/Resolve/Guardian/Guardian.png",
	8446: "perk-images/Styles/Resolve/Demolish/Demolish.png",
	8463: "perk-images/Styles/Resolve/FontOfLife/FontOfLife.png",
	8401: "perk-images/Styles/Resolve/MirrorShell/MirrorShell.png",
	8429: "perk-images/Styles/Resolve/Conditioning/Conditioning.png",
	8444: "perk-images/Styles/Resolve/SecondWind/SecondWind.png",
	8473: "perk-images/Styles/Resolve/BonePlating/BonePlating.png",
	8451: "perk-images/Styles/Resolve/Overgrowth/Overgrowth.png",
	8453: "perk-images/Styles/Resolve/Revitalize/Revitalize.png",

	8200: "perk-images/Styles/7202_Sorcery.png",
	8242: "perk-images/Styles/Sorcery/Unflinching/Unflinching.png",
	8214: "perk-images/Styles/Sorcery/SummonAery/SummonAery.png",
	8229: "perk-images/Styles/Sorcery/ArcaneComet/ArcaneComet.png",
	8230: "perk-images/Styles/Sorcery/PhaseRush/PhaseRush.png",
	8224: "perk-images/Styles/Sorcery/NullifyingOrb/Pokeshield.png",
	8226: "perk-images/Styles/Sorcery/ManaflowBand/ManaflowBand.png",
	8275: "perk-images/Styles/Sorcery/NimbusCloak/6361.png",
	8210: "perk-images/Styles/Sorcery/Transcendence/Transcendence.png",
	8234: "perk-images/Styles/Sorcery/Celerity/CelerityTemp.png",
	8233: "perk-images/Styles/Sorcery/AbsoluteFocus/AbsoluteFocus.png",
	8237: "perk-images/Styles/Sorcery/Scorch/Scorch.png",
	8232: "perk-images/Styles/Sorcery/Waterwalking/Waterwalking.png",
	8236: "perk-images/Styles/Sorcery/GatheringStorm/GatheringStorm.png",
	8299: "perk-images/Styles/Sorcery/LastStand/LastStand.png",
}
