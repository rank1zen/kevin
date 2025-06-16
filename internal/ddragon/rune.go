package ddragon

type RuneID int

type Rune struct {
	ID       RuneID
	IconPath string
	Name     string
}

var (
	Domination  = Rune{8100, "perk-images/Styles/7200_Domination.png", "Domination"}
	Inspiration = Rune{8300, "perk-images/Styles/7203_Whimsy.png", "Inspiration"}
	Precision   = Rune{8000, "perk-images/Styles/7201_Precision.png", "Precision"}
	Resolve     = Rune{8400, "perk-images/Styles/7204_Resolve.png", "Resolve"}
	Sorcery     = Rune{8200, "perk-images/Styles/7202_Sorcery.png", "Sorcery"}
)

var (
	Electrocute          = Rune{8112, "perk-images/Styles/Domination/Electrocute/Electrocute.png", "Electrocute"}
	DarkHarvest          = Rune{8128, "perk-images/Styles/Domination/DarkHarvest/DarkHarvest.png", "Dark Harvest"}
	HailOfBlades         = Rune{9923, "perk-images/Styles/Domination/HailOfBlades/HailOfBlades.png", "Hail of Blades"}
	CheapShot            = Rune{8126, "perk-images/Styles/Domination/CheapShot/CheapShot.png", "Cheap Shot"}
	TasteOfBlood         = Rune{8139, "perk-images/Styles/Domination/TasteOfBlood/GreenTerror_TasteOfBlood.png", "Taste of Blood"}
	SuddenImpact         = Rune{8143, "perk-images/Styles/Domination/SuddenImpact/SuddenImpact.png", "Sudden Impact"}
	SixthSense           = Rune{8137, "perk-images/Styles/Domination/SixthSense/SixthSense.png", "Sixth Sense"}
	GrislyMementos       = Rune{8140, "perk-images/Styles/Domination/GrislyMementos/GrislyMementos.png", "Grisly Mementos"}
	DeepWard             = Rune{8141, "perk-images/Styles/Domination/DeepWard/DeepWard.png", "Deep Ward"}
	TreasureHunter       = Rune{8135, "perk-images/Styles/Domination/TreasureHunter/TreasureHunter.png", "Treasure Hunter"}
	RelentlessHunter     = Rune{8105, "perk-images/Styles/Domination/RelentlessHunter/RelentlessHunter.png", "Relentless Hunter"}
	UltimateHunter       = Rune{8106, "perk-images/Styles/Domination/UltimateHunter/UltimateHunter.png", "Ultimate Hunter"}
	GlacialAugment       = Rune{8351, "perk-images/Styles/Inspiration/GlacialAugment/GlacialAugment.png", "Glacial Augment"}
	UnsealedSpellbook    = Rune{8360, "perk-images/Styles/Inspiration/UnsealedSpellbook/UnsealedSpellbook.png", "Unsealed Spellbook"}
	FirstStrike          = Rune{8369, "perk-images/Styles/Inspiration/FirstStrike/FirstStrike.png", "First Strike"}
	HextechFlashtraption = Rune{8306, "perk-images/Styles/Inspiration/HextechFlashtraption/HextechFlashtraption.png", "Hextech Flashtraption"}
	MagicalFootwear      = Rune{8304, "perk-images/Styles/Inspiration/MagicalFootwear/MagicalFootwear.png", "Magical Footwear"}
	CashBack             = Rune{8321, "perk-images/Styles/Inspiration/CashBack/CashBack2.png", "Cash Back"}
	PerfectTiming        = Rune{8313, "perk-images/Styles/Inspiration/PerfectTiming/AlchemistCabinet.png", "Triple Tonic"}
	TimeWarpTonic        = Rune{8352, "perk-images/Styles/Inspiration/TimeWarpTonic/TimeWarpTonic.png", "Time Warp Tonic"}
	BiscuitDelivery      = Rune{8345, "perk-images/Styles/Inspiration/BiscuitDelivery/BiscuitDelivery.png", "Biscuit Delivery"}
	CosmicInsight        = Rune{8347, "perk-images/Styles/Inspiration/CosmicInsight/CosmicInsight.png", "Cosmic Insight"}
	ApproachVelocity     = Rune{8410, "perk-images/Styles/Resolve/ApproachVelocity/ApproachVelocity.png", "Approach Velocity"}
	JackOfAllTrades      = Rune{8316, "perk-images/Styles/Inspiration/JackOfAllTrades/JackofAllTrades2.png", "Jack Of All Trades"}
	PressTheAttack       = Rune{8005, "perk-images/Styles/Precision/PressTheAttack/PressTheAttack.png", "Press the Attack"}
	LethalTempo          = Rune{8008, "perk-images/Styles/Precision/LethalTempo/LethalTempoTemp.png", "Lethal Tempo"}
	FleetFootwork        = Rune{8021, "perk-images/Styles/Precision/FleetFootwork/FleetFootwork.png", "Fleet Footwork"}
	Conqueror            = Rune{8010, "perk-images/Styles/Precision/Conqueror/Conqueror.png", "Conqueror"}
	AbsorbLife           = Rune{9101, "perk-images/Styles/Precision/AbsorbLife/AbsorbLife.png", "Absorb Life"}
	Triumph              = Rune{9111, "perk-images/Styles/Precision/Triumph.png", "Triumph"}
	PresenceOfMind       = Rune{8009, "perk-images/Styles/Precision/PresenceOfMind/PresenceOfMind.png", "Presence of Mind"}
	LegendAlacrity       = Rune{9104, "perk-images/Styles/Precision/LegendAlacrity/LegendAlacrity.png", "Legend: Alacrity"}
	LegendHaste          = Rune{9105, "perk-images/Styles/Precision/LegendHaste/LegendHaste.png", "Legend: Haste"}
	LegendBloodline      = Rune{9103, "perk-images/Styles/Precision/LegendBloodline/LegendBloodline.png", "Legend: Bloodline"}
	CoupDeGrace          = Rune{8014, "perk-images/Styles/Precision/CoupDeGrace/CoupDeGrace.png", "Coup de Grace"}
	CutDown              = Rune{8017, "perk-images/Styles/Precision/CutDown/CutDown.png", "Cut Down"}
	LastStand            = Rune{8299, "perk-images/Styles/Sorcery/LastStand/LastStand.png", "Last Stand"}
	GraspOfTheUndying    = Rune{8437, "perk-images/Styles/Resolve/GraspOfTheUndying/GraspOfTheUndying.png", "Grasp of the Undying"}
	Aftershock           = Rune{8439, "perk-images/Styles/Resolve/VeteranAftershock/VeteranAftershock.png", "Aftershock"}
	Guardian             = Rune{8465, "perk-images/Styles/Resolve/Guardian/Guardian.png", "Guardian"}
	Demolish             = Rune{8446, "perk-images/Styles/Resolve/Demolish/Demolish.png", "Demolish"}
	FontOfLife           = Rune{8463, "perk-images/Styles/Resolve/FontOfLife/FontOfLife.png", "Font of Life"}
	ShieldBash           = Rune{8401, "perk-images/Styles/Resolve/MirrorShell/MirrorShell.png", "Shield Bash"}
	Conditioning         = Rune{8429, "perk-images/Styles/Resolve/Conditioning/Conditioning.png", "Conditioning"}
	SecondWind           = Rune{8444, "perk-images/Styles/Resolve/SecondWind/SecondWind.png", "Second Wind"}
	BonePlating          = Rune{8473, "perk-images/Styles/Resolve/BonePlating/BonePlating.png", "Bone Plating"}
	Overgrowth           = Rune{8451, "perk-images/Styles/Resolve/Overgrowth/Overgrowth.png", "Overgrowth"}
	Revitalize           = Rune{8453, "perk-images/Styles/Resolve/Revitalize/Revitalize.png", "Revitalize"}
	Unflinching          = Rune{8242, "perk-images/Styles/Sorcery/Unflinching/Unflinching.png", "Unflinching"}
	SummonAery           = Rune{8214, "perk-images/Styles/Sorcery/SummonAery/SummonAery.png", "Summon Aery"}
	ArcaneComet          = Rune{8229, "perk-images/Styles/Sorcery/ArcaneComet/ArcaneComet.png", "Arcane Comet"}
	PhaseRush            = Rune{8230, "perk-images/Styles/Sorcery/PhaseRush/PhaseRush.png", "Phase Rush"}
	NullifyingOrb        = Rune{8224, "perk-images/Styles/Sorcery/NullifyingOrb/Axiom_Arcanist.png", "Axiom Arcanist"}
	ManaflowBand         = Rune{8226, "perk-images/Styles/Sorcery/ManaflowBand/ManaflowBand.png", "Manaflow Band"}
	NimbusCloak          = Rune{8275, "perk-images/Styles/Sorcery/NimbusCloak/6361.png", "Nimbus Cloak"}
	Transcendence        = Rune{8210, "perk-images/Styles/Sorcery/Transcendence/Transcendence.png", "Transcendence"}
	Celerity             = Rune{8234, "perk-images/Styles/Sorcery/Celerity/CelerityTemp.png", "Celerity"}
	AbsoluteFocus        = Rune{8233, "perk-images/Styles/Sorcery/AbsoluteFocus/AbsoluteFocus.png", "Absolute Focus"}
	Scorch               = Rune{8237, "perk-images/Styles/Sorcery/Scorch/Scorch.png", "Scorch"}
	Waterwalking         = Rune{8232, "perk-images/Styles/Sorcery/Waterwalking/Waterwalking.png", "Waterwalking"}
	GatheringStorm       = Rune{8236, "perk-images/Styles/Sorcery/GatheringStorm/GatheringStorm.png", "Gathering Storm"}
)

var RuneMap = map[RuneID]Rune{
	8100: Domination,
	8300: Inspiration,
	8000: Precision,
	8400: Resolve,
	8200: Sorcery,

	8112: Electrocute,
	8128: DarkHarvest,
	9923: HailOfBlades,
	8126: CheapShot,
	8139: TasteOfBlood,
	8143: SuddenImpact,
	8137: SixthSense,
	8140: GrislyMementos,
	8141: DeepWard,
	8135: TreasureHunter,
	8105: RelentlessHunter,
	8106: UltimateHunter,
	8351: GlacialAugment,
	8360: UnsealedSpellbook,
	8369: FirstStrike,
	8306: HextechFlashtraption,
	8304: MagicalFootwear,
	8321: CashBack,
	8313: PerfectTiming,
	8352: TimeWarpTonic,
	8345: BiscuitDelivery,
	8347: CosmicInsight,
	8410: ApproachVelocity,
	8316: JackOfAllTrades,
	8005: PressTheAttack,
	8008: LethalTempo,
	8021: FleetFootwork,
	8010: Conqueror,
	9101: AbsorbLife,
	9111: Triumph,
	8009: PresenceOfMind,
	9104: LegendAlacrity,
	9105: LegendHaste,
	9103: LegendBloodline,
	8014: CoupDeGrace,
	8017: CutDown,
	8299: LastStand,
	8437: GraspOfTheUndying,
	8439: Aftershock,
	8465: Guardian,
	8446: Demolish,
	8463: FontOfLife,
	8401: ShieldBash,
	8429: Conditioning,
	8444: SecondWind,
	8473: BonePlating,
	8451: Overgrowth,
	8453: Revitalize,
	8242: Unflinching,
	8214: SummonAery,
	8229: ArcaneComet,
	8230: PhaseRush,
	8224: NullifyingOrb,
	8226: ManaflowBand,
	8275: NimbusCloak,
	8210: Transcendence,
	8234: Celerity,
	8233: AbsoluteFocus,
	8237: Scorch,
	8232: Waterwalking,
	8236: GatheringStorm,
}
