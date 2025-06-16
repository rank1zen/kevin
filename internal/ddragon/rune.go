package ddragon

type Rune struct {
	ID       int
	IconPath string
	Name     string
}

const (
	RuneTreeDominationID = 8100
	RuneTreeInspirationID = 8300
	RuneTreePrecisionID = 8000
	RuneTreeResolveID = 8400
	RuneTreeSorceryID = 8200
)

var (
	RuneTreeDomination  = Rune{8100, "perk-images/Styles/7200_Domination.png", "Domination"}
	RuneTreeInspiration = Rune{8300, "perk-images/Styles/7203_Whimsy.png", "Inspiration"}
	RuneTreePrecision   = Rune{8000, "perk-images/Styles/7201_Precision.png", "Precision"}
	RuneTreeResolve     = Rune{8400, "perk-images/Styles/7204_Resolve.png", "Resolve"}
	RuneTreeSorcery     = Rune{8200, "perk-images/Styles/7202_Sorcery.png", "Sorcery"}
)

var (
	RuneElectrocute          = Rune{8112, "perk-images/Styles/Domination/Electrocute/Electrocute.png", "Electrocute"}
	RuneDarkHarvest          = Rune{8128, "perk-images/Styles/Domination/DarkHarvest/DarkHarvest.png", "Dark Harvest"}
	RuneHailOfBlades         = Rune{9923, "perk-images/Styles/Domination/HailOfBlades/HailOfBlades.png", "Hail of Blades"}
	RuneCheapShot            = Rune{8126, "perk-images/Styles/Domination/CheapShot/CheapShot.png", "Cheap Shot"}
	RuneTasteOfBlood         = Rune{8139, "perk-images/Styles/Domination/TasteOfBlood/GreenTerror_TasteOfBlood.png", "Taste of Blood"}
	RuneSuddenImpact         = Rune{8143, "perk-images/Styles/Domination/SuddenImpact/SuddenImpact.png", "Sudden Impact"}
	RuneSixthSense           = Rune{8137, "perk-images/Styles/Domination/SixthSense/SixthSense.png", "Sixth Sense"}
	RuneGrislyMementos       = Rune{8140, "perk-images/Styles/Domination/GrislyMementos/GrislyMementos.png", "Grisly Mementos"}
	RuneDeepWard             = Rune{8141, "perk-images/Styles/Domination/DeepWard/DeepWard.png", "Deep Ward"}
	RuneTreasureHunter       = Rune{8135, "perk-images/Styles/Domination/TreasureHunter/TreasureHunter.png", "Treasure Hunter"}
	RuneRelentlessHunter     = Rune{8105, "perk-images/Styles/Domination/RelentlessHunter/RelentlessHunter.png", "Relentless Hunter"}
	RuneUltimateHunter       = Rune{8106, "perk-images/Styles/Domination/UltimateHunter/UltimateHunter.png", "Ultimate Hunter"}
	RuneGlacialAugment       = Rune{8351, "perk-images/Styles/Inspiration/GlacialAugment/GlacialAugment.png", "Glacial Augment"}
	RuneUnsealedSpellbook    = Rune{8360, "perk-images/Styles/Inspiration/UnsealedSpellbook/UnsealedSpellbook.png", "Unsealed Spellbook"}
	RuneFirstStrike          = Rune{8369, "perk-images/Styles/Inspiration/FirstStrike/FirstStrike.png", "First Strike"}
	RuneHextechFlashtraption = Rune{8306, "perk-images/Styles/Inspiration/HextechFlashtraption/HextechFlashtraption.png", "Hextech Flashtraption"}
	RuneMagicalFootwear      = Rune{8304, "perk-images/Styles/Inspiration/MagicalFootwear/MagicalFootwear.png", "Magical Footwear"}
	RuneCashBack             = Rune{8321, "perk-images/Styles/Inspiration/CashBack/CashBack2.png", "Cash Back"}
	RunePerfectTiming        = Rune{8313, "perk-images/Styles/Inspiration/PerfectTiming/AlchemistCabinet.png", "Triple Tonic"}
	RuneTimeWarpTonic        = Rune{8352, "perk-images/Styles/Inspiration/TimeWarpTonic/TimeWarpTonic.png", "Time Warp Tonic"}
	RuneBiscuitDelivery      = Rune{8345, "perk-images/Styles/Inspiration/BiscuitDelivery/BiscuitDelivery.png", "Biscuit Delivery"}
	RuneCosmicInsight        = Rune{8347, "perk-images/Styles/Inspiration/CosmicInsight/CosmicInsight.png", "Cosmic Insight"}
	RuneApproachVelocity     = Rune{8410, "perk-images/Styles/Resolve/ApproachVelocity/ApproachVelocity.png", "Approach Velocity"}
	RuneJackOfAllTrades      = Rune{8316, "perk-images/Styles/Inspiration/JackOfAllTrades/JackofAllTrades2.png", "Jack Of All Trades"}
	RunePressTheAttack       = Rune{8005, "perk-images/Styles/Precision/PressTheAttack/PressTheAttack.png", "Press the Attack"}
	RuneLethalTempo          = Rune{8008, "perk-images/Styles/Precision/LethalTempo/LethalTempoTemp.png", "Lethal Tempo"}
	RuneFleetFootwork        = Rune{8021, "perk-images/Styles/Precision/FleetFootwork/FleetFootwork.png", "Fleet Footwork"}
	RuneConqueror            = Rune{8010, "perk-images/Styles/Precision/Conqueror/Conqueror.png", "Conqueror"}
	RuneAbsorbLife           = Rune{9101, "perk-images/Styles/Precision/AbsorbLife/AbsorbLife.png", "Absorb Life"}
	RuneTriumph              = Rune{9111, "perk-images/Styles/Precision/Triumph.png", "Triumph"}
	RunePresenceOfMind       = Rune{8009, "perk-images/Styles/Precision/PresenceOfMind/PresenceOfMind.png", "Presence of Mind"}
	RuneLegendAlacrity       = Rune{9104, "perk-images/Styles/Precision/LegendAlacrity/LegendAlacrity.png", "Legend: Alacrity"}
	RuneLegendHaste          = Rune{9105, "perk-images/Styles/Precision/LegendHaste/LegendHaste.png", "Legend: Haste"}
	RuneLegendBloodline      = Rune{9103, "perk-images/Styles/Precision/LegendBloodline/LegendBloodline.png", "Legend: Bloodline"}
	RuneCoupDeGrace          = Rune{8014, "perk-images/Styles/Precision/CoupDeGrace/CoupDeGrace.png", "Coup de Grace"}
	RuneCutDown              = Rune{8017, "perk-images/Styles/Precision/CutDown/CutDown.png", "Cut Down"}
	RuneLastStand            = Rune{8299, "perk-images/Styles/Sorcery/LastStand/LastStand.png", "Last Stand"}
	RuneGraspOfTheUndying    = Rune{8437, "perk-images/Styles/Resolve/GraspOfTheUndying/GraspOfTheUndying.png", "Grasp of the Undying"}
	RuneAftershock           = Rune{8439, "perk-images/Styles/Resolve/VeteranAftershock/VeteranAftershock.png", "Aftershock"}
	RuneGuardian             = Rune{8465, "perk-images/Styles/Resolve/Guardian/Guardian.png", "Guardian"}
	RuneDemolish             = Rune{8446, "perk-images/Styles/Resolve/Demolish/Demolish.png", "Demolish"}
	RuneFontOfLife           = Rune{8463, "perk-images/Styles/Resolve/FontOfLife/FontOfLife.png", "Font of Life"}
	RuneShieldBash           = Rune{8401, "perk-images/Styles/Resolve/MirrorShell/MirrorShell.png", "Shield Bash"}
	RuneConditioning         = Rune{8429, "perk-images/Styles/Resolve/Conditioning/Conditioning.png", "Conditioning"}
	RuneSecondWind           = Rune{8444, "perk-images/Styles/Resolve/SecondWind/SecondWind.png", "Second Wind"}
	RuneBonePlating          = Rune{8473, "perk-images/Styles/Resolve/BonePlating/BonePlating.png", "Bone Plating"}
	RuneOvergrowth           = Rune{8451, "perk-images/Styles/Resolve/Overgrowth/Overgrowth.png", "Overgrowth"}
	RuneRevitalize           = Rune{8453, "perk-images/Styles/Resolve/Revitalize/Revitalize.png", "Revitalize"}
	RuneUnflinching          = Rune{8242, "perk-images/Styles/Sorcery/Unflinching/Unflinching.png", "Unflinching"}
	RuneSummonAery           = Rune{8214, "perk-images/Styles/Sorcery/SummonAery/SummonAery.png", "Summon Aery"}
	RuneArcaneComet          = Rune{8229, "perk-images/Styles/Sorcery/ArcaneComet/ArcaneComet.png", "Arcane Comet"}
	RunePhaseRush            = Rune{8230, "perk-images/Styles/Sorcery/PhaseRush/PhaseRush.png", "Phase Rush"}
	RuneNullifyingOrb        = Rune{8224, "perk-images/Styles/Sorcery/NullifyingOrb/Axiom_Arcanist.png", "Axiom Arcanist"}
	RuneManaflowBand         = Rune{8226, "perk-images/Styles/Sorcery/ManaflowBand/ManaflowBand.png", "Manaflow Band"}
	RuneNimbusCloak          = Rune{8275, "perk-images/Styles/Sorcery/NimbusCloak/6361.png", "Nimbus Cloak"}
	RuneTranscendence        = Rune{8210, "perk-images/Styles/Sorcery/Transcendence/Transcendence.png", "Transcendence"}
	RuneCelerity             = Rune{8234, "perk-images/Styles/Sorcery/Celerity/CelerityTemp.png", "Celerity"}
	RuneAbsoluteFocus        = Rune{8233, "perk-images/Styles/Sorcery/AbsoluteFocus/AbsoluteFocus.png", "Absolute Focus"}
	RuneScorch               = Rune{8237, "perk-images/Styles/Sorcery/Scorch/Scorch.png", "Scorch"}
	RuneWaterwalking         = Rune{8232, "perk-images/Styles/Sorcery/Waterwalking/Waterwalking.png", "Waterwalking"}
	RuneGatheringStorm       = Rune{8236, "perk-images/Styles/Sorcery/GatheringStorm/GatheringStorm.png", "Gathering Storm"}
)

var RuneMap = map[int]Rune{
	8100: RuneTreeDomination,
	8300: RuneTreeInspiration,
	8000: RuneTreePrecision,
	8400: RuneTreeResolve,
	8200: RuneTreeSorcery,

	8112: RuneElectrocute,
	8128: RuneDarkHarvest,
	9923: RuneHailOfBlades,
	8126: RuneCheapShot,
	8139: RuneTasteOfBlood,
	8143: RuneSuddenImpact,
	8137: RuneSixthSense,
	8140: RuneGrislyMementos,
	8141: RuneDeepWard,
	8135: RuneTreasureHunter,
	8105: RuneRelentlessHunter,
	8106: RuneUltimateHunter,
	8351: RuneGlacialAugment,
	8360: RuneUnsealedSpellbook,
	8369: RuneFirstStrike,
	8306: RuneHextechFlashtraption,
	8304: RuneMagicalFootwear,
	8321: RuneCashBack,
	8313: RunePerfectTiming,
	8352: RuneTimeWarpTonic,
	8345: RuneBiscuitDelivery,
	8347: RuneCosmicInsight,
	8410: RuneApproachVelocity,
	8316: RuneJackOfAllTrades,
	8005: RunePressTheAttack,
	8008: RuneLethalTempo,
	8021: RuneFleetFootwork,
	8010: RuneConqueror,
	9101: RuneAbsorbLife,
	9111: RuneTriumph,
	8009: RunePresenceOfMind,
	9104: RuneLegendAlacrity,
	9105: RuneLegendHaste,
	9103: RuneLegendBloodline,
	8014: RuneCoupDeGrace,
	8017: RuneCutDown,
	8299: RuneLastStand,
	8437: RuneGraspOfTheUndying,
	8439: RuneAftershock,
	8465: RuneGuardian,
	8446: RuneDemolish,
	8463: RuneFontOfLife,
	8401: RuneShieldBash,
	8429: RuneConditioning,
	8444: RuneSecondWind,
	8473: RuneBonePlating,
	8451: RuneOvergrowth,
	8453: RuneRevitalize,
	8242: RuneUnflinching,
	8214: RuneSummonAery,
	8229: RuneArcaneComet,
	8230: RunePhaseRush,
	8224: RuneNullifyingOrb,
	8226: RuneManaflowBand,
	8275: RuneNimbusCloak,
	8210: RuneTranscendence,
	8234: RuneCelerity,
	8233: RuneAbsoluteFocus,
	8237: RuneScorch,
	8232: RuneWaterwalking,
	8236: RuneGatheringStorm,
}
