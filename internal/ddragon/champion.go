package ddragon

type Champion struct {
	ID         int
	Name       string
	FullIcon   string
	Sprite     string
	X, Y, W, H int
}

const (
	ChampionAatroxID       = 266
	ChampionAhriID         = 103
	ChampionAkaliID        = 84
	ChampionAkshanID       = 166
	ChampionAlistarID      = 12
	ChampionAmbessaID      = 799
	ChampionAmumuID        = 32
	ChampionAniviaID       = 34
	ChampionAnnieID        = 1
	ChampionApheliosID     = 523
	ChampionAsheID         = 22
	ChampionAurelionSolID  = 136
	ChampionAuroraID       = 893
	ChampionAzirID         = 268
	ChampionBardID         = 432
	ChampionBelvethID      = 200
	ChampionBlitzcrankID   = 53
	ChampionBrandID        = 63
	ChampionBraumID        = 201
	ChampionBriarID        = 233
	ChampionCaitlynID      = 51
	ChampionCamilleID      = 164
	ChampionCassiopeiaID   = 69
	ChampionChogathID      = 31
	ChampionCorkiID        = 42
	ChampionDariusID       = 122
	ChampionDianaID        = 131
	ChampionDravenID       = 119
	ChampionDrMundoID      = 36
	ChampionEkkoID         = 245
	ChampionEliseID        = 60
	ChampionEvelynnID      = 28
	ChampionEzrealID       = 81
	ChampionFiddlesticksID = 9
	ChampionFioraID        = 114
	ChampionFizzID         = 105
	ChampionGalioID        = 3
	ChampionGangplankID    = 41
	ChampionGarenID        = 86
	ChampionGnarID         = 150
	ChampionGragasID       = 79
	ChampionGravesID       = 104
	ChampionGwenID         = 887
	ChampionHecarimID      = 120
	ChampionHeimerdingerID = 74
	ChampionHweiID         = 910
	ChampionIllaoiID       = 420
	ChampionIreliaID       = 39
	ChampionIvernID        = 427
	ChampionJannaID        = 40
	ChampionJarvanIVID     = 59
	ChampionJaxID          = 24
	ChampionJayceID        = 126
	ChampionJhinID         = 202
	ChampionJinxID         = 222
	ChampionKaisaID        = 145
	ChampionKalistaID      = 429
	ChampionKarmaID        = 43
	ChampionKarthusID      = 30
	ChampionKassadinID     = 38
	ChampionKatarinaID     = 55
	ChampionKayleID        = 10
	ChampionKaynID         = 141
	ChampionKennenID       = 85
	ChampionKhazixID       = 121
	ChampionKindredID      = 203
	ChampionKledID         = 240
	ChampionKogMawID       = 96
	ChampionKSanteID       = 897
	ChampionLeblancID      = 7
	ChampionLeeSinID       = 64
	ChampionLeonaID        = 89
	ChampionLilliaID       = 876
	ChampionLissandraID    = 127
	ChampionLucianID       = 236
	ChampionLuluID         = 117
	ChampionLuxID          = 99
	ChampionMalphiteID     = 54
	ChampionMalzaharID     = 90
	ChampionMaokaiID       = 57
	ChampionMasterYiID     = 11
	ChampionMelID          = 800
	ChampionMilioID        = 902
	ChampionMissFortuneID  = 21
	ChampionMonkeyKingID   = 62
	ChampionMordekaiserID  = 82
	ChampionMorganaID      = 25
	ChampionNaafiriID      = 950
	ChampionNamiID         = 267
	ChampionNasusID        = 75
	ChampionNautilusID     = 111
	ChampionNeekoID        = 518
	ChampionNidaleeID      = 76
	ChampionNilahID        = 895
	ChampionNocturneID     = 56
	ChampionNunuID         = 20
	ChampionOlafID         = 2
	ChampionOriannaID      = 61
	ChampionOrnnID         = 516
	ChampionPantheonID     = 80
	ChampionPoppyID        = 78
	ChampionPykeID         = 555
	ChampionQiyanaID       = 246
	ChampionQuinnID        = 133
	ChampionRakanID        = 497
	ChampionRammusID       = 33
	ChampionRekSaiID       = 421
	ChampionRellID         = 526
	ChampionRenataID       = 888
	ChampionRenektonID     = 58
	ChampionRengarID       = 107
	ChampionRivenID        = 92
	ChampionRumbleID       = 68
	ChampionRyzeID         = 13
	ChampionSamiraID       = 360
	ChampionSejuaniID      = 113
	ChampionSennaID        = 235
	ChampionSeraphineID    = 147
	ChampionSettID         = 875
	ChampionShacoID        = 35
	ChampionShenID         = 98
	ChampionShyvanaID      = 102
	ChampionSingedID       = 27
	ChampionSionID         = 14
	ChampionSivirID        = 15
	ChampionSkarnerID      = 72
	ChampionSmolderID      = 901
	ChampionSonaID         = 37
	ChampionSorakaID       = 16
	ChampionSwainID        = 50
	ChampionSylasID        = 517
	ChampionSyndraID       = 134
	ChampionTahmKenchID    = 223
	ChampionTaliyahID      = 163
	ChampionTalonID        = 91
	ChampionTaricID        = 44
	ChampionTeemoID        = 17
	ChampionThreshID       = 412
	ChampionTristanaID     = 18
	ChampionTrundleID      = 48
	ChampionTryndamereID   = 23
	ChampionTwistedFateID  = 4
	ChampionTwitchID       = 29
	ChampionUdyrID         = 77
	ChampionUrgotID        = 6
	ChampionVarusID        = 110
	ChampionVayneID        = 67
	ChampionVeigarID       = 45
	ChampionVelkozID       = 161
	ChampionVexID          = 711
	ChampionViID           = 254
	ChampionViegoID        = 234
	ChampionViktorID       = 112
	ChampionVladimirID     = 8
	ChampionVolibearID     = 106
	ChampionWarwickID      = 19
	ChampionXayahID        = 498
	ChampionXerathID       = 101
	ChampionXinZhaoID      = 5
	ChampionYasuoID        = 157
	ChampionYoneID         = 777
	ChampionYorickID       = 83
	ChampionYuumiID        = 350
	ChampionZacID          = 154
	ChampionZedID          = 238
	ChampionZeriID         = 221
	ChampionZiggsID        = 115
	ChampionZileanID       = 26
	ChampionZoeID          = 142
	ChampionZyraID         = 143
)

var ChampionMap = map[int]Champion{
	266: ChampionAatrox,
	103: ChampionAhri,
	84:  ChampionAkali,
	166: ChampionAkshan,
	12:  ChampionAlistar,
	799: ChampionAmbessa,
	32:  ChampionAmumu,
	34:  ChampionAnivia,
	1:   ChampionAnnie,
	523: ChampionAphelios,
	22:  ChampionAshe,
	136: ChampionAurelionSol,
	893: ChampionAurora,
	268: ChampionAzir,
	432: ChampionBard,
	200: ChampionBelveth,
	53:  ChampionBlitzcrank,
	63:  ChampionBrand,
	201: ChampionBraum,
	233: ChampionBriar,
	51:  ChampionCaitlyn,
	164: ChampionCamille,
	69:  ChampionCassiopeia,
	31:  ChampionChogath,
	42:  ChampionCorki,
	122: ChampionDarius,
	131: ChampionDiana,
	119: ChampionDraven,
	36:  ChampionDrMundo,
	245: ChampionEkko,
	60:  ChampionElise,
	28:  ChampionEvelynn,
	81:  ChampionEzreal,
	9:   ChampionFiddlesticks,
	114: ChampionFiora,
	105: ChampionFizz,
	3:   ChampionGalio,
	41:  ChampionGangplank,
	86:  ChampionGaren,
	150: ChampionGnar,
	79:  ChampionGragas,
	104: ChampionGraves,
	887: ChampionGwen,
	120: ChampionHecarim,
	74:  ChampionHeimerdinger,
	910: ChampionHwei,
	420: ChampionIllaoi,
	39:  ChampionIrelia,
	427: ChampionIvern,
	40:  ChampionJanna,
	59:  ChampionJarvanIV,
	24:  ChampionJax,
	126: ChampionJayce,
	202: ChampionJhin,
	222: ChampionJinx,
	145: ChampionKaisa,
	429: ChampionKalista,
	43:  ChampionKarma,
	30:  ChampionKarthus,
	38:  ChampionKassadin,
	55:  ChampionKatarina,
	10:  ChampionKayle,
	141: ChampionKayn,
	85:  ChampionKennen,
	121: ChampionKhazix,
	203: ChampionKindred,
	240: ChampionKled,
	96:  ChampionKogMaw,
	897: ChampionKSante,
	7:   ChampionLeblanc,
	64:  ChampionLeeSin,
	89:  ChampionLeona,
	876: ChampionLillia,
	127: ChampionLissandra,
	236: ChampionLucian,
	117: ChampionLulu,
	99:  ChampionLux,
	54:  ChampionMalphite,
	90:  ChampionMalzahar,
	57:  ChampionMaokai,
	11:  ChampionMasterYi,
	800: ChampionMel,
	902: ChampionMilio,
	21:  ChampionMissFortune,
	62:  ChampionMonkeyKing,
	82:  ChampionMordekaiser,
	25:  ChampionMorgana,
	950: ChampionNaafiri,
	267: ChampionNami,
	75:  ChampionNasus,
	111: ChampionNautilus,
	518: ChampionNeeko,
	76:  ChampionNidalee,
	895: ChampionNilah,
	56:  ChampionNocturne,
	20:  ChampionNunu,
	2:   ChampionOlaf,
	61:  ChampionOrianna,
	516: ChampionOrnn,
	80:  ChampionPantheon,
	78:  ChampionPoppy,
	555: ChampionPyke,
	246: ChampionQiyana,
	133: ChampionQuinn,
	497: ChampionRakan,
	33:  ChampionRammus,
	421: ChampionRekSai,
	526: ChampionRell,
	888: ChampionRenata,
	58:  ChampionRenekton,
	107: ChampionRengar,
	92:  ChampionRiven,
	68:  ChampionRumble,
	13:  ChampionRyze,
	360: ChampionSamira,
	113: ChampionSejuani,
	235: ChampionSenna,
	147: ChampionSeraphine,
	875: ChampionSett,
	35:  ChampionShaco,
	98:  ChampionShen,
	102: ChampionShyvana,
	27:  ChampionSinged,
	14:  ChampionSion,
	15:  ChampionSivir,
	72:  ChampionSkarner,
	901: ChampionSmolder,
	37:  ChampionSona,
	16:  ChampionSoraka,
	50:  ChampionSwain,
	517: ChampionSylas,
	134: ChampionSyndra,
	223: ChampionTahmKench,
	163: ChampionTaliyah,
	91:  ChampionTalon,
	44:  ChampionTaric,
	17:  ChampionTeemo,
	412: ChampionThresh,
	18:  ChampionTristana,
	48:  ChampionTrundle,
	23:  ChampionTryndamere,
	4:   ChampionTwistedFate,
	29:  ChampionTwitch,
	77:  ChampionUdyr,
	6:   ChampionUrgot,
	110: ChampionVarus,
	67:  ChampionVayne,
	45:  ChampionVeigar,
	161: ChampionVelkoz,
	711: ChampionVex,
	254: ChampionVi,
	234: ChampionViego,
	112: ChampionViktor,
	8:   ChampionVladimir,
	106: ChampionVolibear,
	19:  ChampionWarwick,
	498: ChampionXayah,
	101: ChampionXerath,
	5:   ChampionXinZhao,
	157: ChampionYasuo,
	777: ChampionYone,
	83:  ChampionYorick,
	350: ChampionYuumi,
	154: ChampionZac,
	238: ChampionZed,
	221: ChampionZeri,
	115: ChampionZiggs,
	26:  ChampionZilean,
	142: ChampionZoe,
	143: ChampionZyra,
}

var (
	ChampionAatrox       = Champion{266, "Aatrox", "Aatrox.png", "champion0.png", 0, 0, 48, 48}
	ChampionAhri         = Champion{103, "Ahri", "Ahri.png", "champion0.png", 48, 0, 48, 48}
	ChampionAkali        = Champion{84, "Akali", "Akali.png", "champion0.png", 96, 0, 48, 48}
	ChampionAkshan       = Champion{166, "Akshan", "Akshan.png", "champion0.png", 144, 0, 48, 48}
	ChampionAlistar      = Champion{12, "Alistar", "Alistar.png", "champion0.png", 192, 0, 48, 48}
	ChampionAmbessa      = Champion{799, "Ambessa", "Ambessa.png", "champion0.png", 240, 0, 48, 48}
	ChampionAmumu        = Champion{32, "Amumu", "Amumu.png", "champion0.png", 288, 0, 48, 48}
	ChampionAnivia       = Champion{34, "Anivia", "Anivia.png", "champion0.png", 336, 0, 48, 48}
	ChampionAnnie        = Champion{1, "Annie", "Annie.png", "champion0.png", 384, 0, 48, 48}
	ChampionAphelios     = Champion{523, "Aphelios", "Aphelios.png", "champion0.png", 432, 0, 48, 48}
	ChampionAshe         = Champion{22, "Ashe", "Ashe.png", "champion0.png", 0, 48, 48, 48}
	ChampionAurelionSol  = Champion{136, "Aurelion Sol", "AurelionSol.png", "champion0.png", 48, 48, 48, 48}
	ChampionAurora       = Champion{893, "Aurora", "Aurora.png", "champion0.png", 96, 48, 48, 48}
	ChampionAzir         = Champion{268, "Azir", "Azir.png", "champion0.png", 144, 48, 48, 48}
	ChampionBard         = Champion{432, "Bard", "Bard.png", "champion0.png", 192, 48, 48, 48}
	ChampionBelveth      = Champion{200, "Bel'Veth", "Belveth.png", "champion0.png", 240, 48, 48, 48}
	ChampionBlitzcrank   = Champion{53, "Blitzcrank", "Blitzcrank.png", "champion0.png", 288, 48, 48, 48}
	ChampionBrand        = Champion{63, "Brand", "Brand.png", "champion0.png", 336, 48, 48, 48}
	ChampionBraum        = Champion{201, "Braum", "Braum.png", "champion0.png", 384, 48, 48, 48}
	ChampionBriar        = Champion{233, "Briar", "Briar.png", "champion0.png", 432, 48, 48, 48}
	ChampionCaitlyn      = Champion{51, "Caitlyn", "Caitlyn.png", "champion0.png", 0, 96, 48, 48}
	ChampionCamille      = Champion{164, "Camille", "Camille.png", "champion0.png", 48, 96, 48, 48}
	ChampionCassiopeia   = Champion{69, "Cassiopeia", "Cassiopeia.png", "champion0.png", 96, 96, 48, 48}
	ChampionChogath      = Champion{31, "Cho'Gath", "Chogath.png", "champion0.png", 144, 96, 48, 48}
	ChampionCorki        = Champion{42, "Corki", "Corki.png", "champion0.png", 192, 96, 48, 48}
	ChampionDarius       = Champion{122, "Darius", "Darius.png", "champion0.png", 240, 96, 48, 48}
	ChampionDiana        = Champion{131, "Diana", "Diana.png", "champion0.png", 288, 96, 48, 48}
	ChampionDraven       = Champion{119, "Draven", "Draven.png", "champion0.png", 336, 96, 48, 48}
	ChampionDrMundo      = Champion{36, "Dr. Mundo", "DrMundo.png", "champion0.png", 384, 96, 48, 48}
	ChampionEkko         = Champion{245, "Ekko", "Ekko.png", "champion0.png", 432, 96, 48, 48}
	ChampionElise        = Champion{60, "Elise", "Elise.png", "champion1.png", 0, 0, 48, 48}
	ChampionEvelynn      = Champion{28, "Evelynn", "Evelynn.png", "champion1.png", 48, 0, 48, 48}
	ChampionEzreal       = Champion{81, "Ezreal", "Ezreal.png", "champion1.png", 96, 0, 48, 48}
	ChampionFiddlesticks = Champion{9, "Fiddlesticks", "Fiddlesticks.png", "champion1.png", 144, 0, 48, 48}
	ChampionFiora        = Champion{114, "Fiora", "Fiora.png", "champion1.png", 192, 0, 48, 48}
	ChampionFizz         = Champion{105, "Fizz", "Fizz.png", "champion1.png", 240, 0, 48, 48}
	ChampionGalio        = Champion{3, "Galio", "Galio.png", "champion1.png", 288, 0, 48, 48}
	ChampionGangplank    = Champion{41, "Gangplank", "Gangplank.png", "champion1.png", 336, 0, 48, 48}
	ChampionGaren        = Champion{86, "Garen", "Garen.png", "champion1.png", 384, 0, 48, 48}
	ChampionGnar         = Champion{150, "Gnar", "Gnar.png", "champion1.png", 432, 0, 48, 48}
	ChampionGragas       = Champion{79, "Gragas", "Gragas.png", "champion1.png", 0, 48, 48, 48}
	ChampionGraves       = Champion{104, "Graves", "Graves.png", "champion1.png", 48, 48, 48, 48}
	ChampionGwen         = Champion{887, "Gwen", "Gwen.png", "champion1.png", 96, 48, 48, 48}
	ChampionHecarim      = Champion{120, "Hecarim", "Hecarim.png", "champion1.png", 144, 48, 48, 48}
	ChampionHeimerdinger = Champion{74, "Heimerdinger", "Heimerdinger.png", "champion1.png", 192, 48, 48, 48}
	ChampionHwei         = Champion{910, "Hwei", "Hwei.png", "champion1.png", 240, 48, 48, 48}
	ChampionIllaoi       = Champion{420, "Illaoi", "Illaoi.png", "champion1.png", 288, 48, 48, 48}
	ChampionIrelia       = Champion{39, "Irelia", "Irelia.png", "champion1.png", 336, 48, 48, 48}
	ChampionIvern        = Champion{427, "Ivern", "Ivern.png", "champion1.png", 384, 48, 48, 48}
	ChampionJanna        = Champion{40, "Janna", "Janna.png", "champion1.png", 432, 48, 48, 48}
	ChampionJarvanIV     = Champion{59, "Jarvan IV", "JarvanIV.png", "champion1.png", 0, 96, 48, 48}
	ChampionJax          = Champion{24, "Jax", "Jax.png", "champion1.png", 48, 96, 48, 48}
	ChampionJayce        = Champion{126, "Jayce", "Jayce.png", "champion1.png", 96, 96, 48, 48}
	ChampionJhin         = Champion{202, "Jhin", "Jhin.png", "champion1.png", 144, 96, 48, 48}
	ChampionJinx         = Champion{222, "Jinx", "Jinx.png", "champion1.png", 192, 96, 48, 48}
	ChampionKaisa        = Champion{145, "Kai'Sa", "Kaisa.png", "champion1.png", 240, 96, 48, 48}
	ChampionKalista      = Champion{429, "Kalista", "Kalista.png", "champion1.png", 288, 96, 48, 48}
	ChampionKarma        = Champion{43, "Karma", "Karma.png", "champion1.png", 336, 96, 48, 48}
	ChampionKarthus      = Champion{30, "Karthus", "Karthus.png", "champion1.png", 384, 96, 48, 48}
	ChampionKassadin     = Champion{38, "Kassadin", "Kassadin.png", "champion1.png", 432, 96, 48, 48}
	ChampionKatarina     = Champion{55, "Katarina", "Katarina.png", "champion2.png", 0, 0, 48, 48}
	ChampionKayle        = Champion{10, "Kayle", "Kayle.png", "champion2.png", 48, 0, 48, 48}
	ChampionKayn         = Champion{141, "Kayn", "Kayn.png", "champion2.png", 96, 0, 48, 48}
	ChampionKennen       = Champion{85, "Kennen", "Kennen.png", "champion2.png", 144, 0, 48, 48}
	ChampionKhazix       = Champion{121, "Kha'Zix", "Khazix.png", "champion2.png", 192, 0, 48, 48}
	ChampionKindred      = Champion{203, "Kindred", "Kindred.png", "champion2.png", 240, 0, 48, 48}
	ChampionKled         = Champion{240, "Kled", "Kled.png", "champion2.png", 288, 0, 48, 48}
	ChampionKogMaw       = Champion{96, "Kog'Maw", "KogMaw.png", "champion2.png", 336, 0, 48, 48}
	ChampionKSante       = Champion{897, "K'Sante", "KSante.png", "champion2.png", 384, 0, 48, 48}
	ChampionLeblanc      = Champion{7, "LeBlanc", "Leblanc.png", "champion2.png", 432, 0, 48, 48}
	ChampionLeeSin       = Champion{64, "Lee Sin", "LeeSin.png", "champion2.png", 0, 48, 48, 48}
	ChampionLeona        = Champion{89, "Leona", "Leona.png", "champion2.png", 48, 48, 48, 48}
	ChampionLillia       = Champion{876, "Lillia", "Lillia.png", "champion2.png", 96, 48, 48, 48}
	ChampionLissandra    = Champion{127, "Lissandra", "Lissandra.png", "champion2.png", 144, 48, 48, 48}
	ChampionLucian       = Champion{236, "Lucian", "Lucian.png", "champion2.png", 192, 48, 48, 48}
	ChampionLulu         = Champion{117, "Lulu", "Lulu.png", "champion2.png", 240, 48, 48, 48}
	ChampionLux          = Champion{99, "Lux", "Lux.png", "champion2.png", 288, 48, 48, 48}
	ChampionMalphite     = Champion{54, "Malphite", "Malphite.png", "champion2.png", 336, 48, 48, 48}
	ChampionMalzahar     = Champion{90, "Malzahar", "Malzahar.png", "champion2.png", 384, 48, 48, 48}
	ChampionMaokai       = Champion{57, "Maokai", "Maokai.png", "champion2.png", 432, 48, 48, 48}
	ChampionMasterYi     = Champion{11, "Master Yi", "MasterYi.png", "champion2.png", 0, 96, 48, 48}
	ChampionMel          = Champion{800, "Mel", "Mel.png", "champion2.png", 48, 96, 48, 48}
	ChampionMilio        = Champion{902, "Milio", "Milio.png", "champion2.png", 96, 96, 48, 48}
	ChampionMissFortune  = Champion{21, "Miss Fortune", "MissFortune.png", "champion2.png", 144, 96, 48, 48}
	ChampionMonkeyKing   = Champion{62, "Wukong", "MonkeyKing.png", "champion2.png", 192, 96, 48, 48}
	ChampionMordekaiser  = Champion{82, "Mordekaiser", "Mordekaiser.png", "champion2.png", 240, 96, 48, 48}
	ChampionMorgana      = Champion{25, "Morgana", "Morgana.png", "champion2.png", 288, 96, 48, 48}
	ChampionNaafiri      = Champion{950, "Naafiri", "Naafiri.png", "champion2.png", 336, 96, 48, 48}
	ChampionNami         = Champion{267, "Nami", "Nami.png", "champion2.png", 384, 96, 48, 48}
	ChampionNasus        = Champion{75, "Nasus", "Nasus.png", "champion2.png", 432, 96, 48, 48}
	ChampionNautilus     = Champion{111, "Nautilus", "Nautilus.png", "champion3.png", 0, 0, 48, 48}
	ChampionNeeko        = Champion{518, "Neeko", "Neeko.png", "champion3.png", 48, 0, 48, 48}
	ChampionNidalee      = Champion{76, "Nidalee", "Nidalee.png", "champion3.png", 96, 0, 48, 48}
	ChampionNilah        = Champion{895, "Nilah", "Nilah.png", "champion3.png", 144, 0, 48, 48}
	ChampionNocturne     = Champion{56, "Nocturne", "Nocturne.png", "champion3.png", 192, 0, 48, 48}
	ChampionNunu         = Champion{20, "Nunu & Willump", "Nunu.png", "champion3.png", 240, 0, 48, 48}
	ChampionOlaf         = Champion{2, "Olaf", "Olaf.png", "champion3.png", 288, 0, 48, 48}
	ChampionOrianna      = Champion{61, "Orianna", "Orianna.png", "champion3.png", 336, 0, 48, 48}
	ChampionOrnn         = Champion{516, "Ornn", "Ornn.png", "champion3.png", 384, 0, 48, 48}
	ChampionPantheon     = Champion{80, "Pantheon", "Pantheon.png", "champion3.png", 432, 0, 48, 48}
	ChampionPoppy        = Champion{78, "Poppy", "Poppy.png", "champion3.png", 0, 48, 48, 48}
	ChampionPyke         = Champion{555, "Pyke", "Pyke.png", "champion3.png", 48, 48, 48, 48}
	ChampionQiyana       = Champion{246, "Qiyana", "Qiyana.png", "champion3.png", 96, 48, 48, 48}
	ChampionQuinn        = Champion{133, "Quinn", "Quinn.png", "champion3.png", 144, 48, 48, 48}
	ChampionRakan        = Champion{497, "Rakan", "Rakan.png", "champion3.png", 192, 48, 48, 48}
	ChampionRammus       = Champion{33, "Rammus", "Rammus.png", "champion3.png", 240, 48, 48, 48}
	ChampionRekSai       = Champion{421, "Rek'Sai", "RekSai.png", "champion3.png", 288, 48, 48, 48}
	ChampionRell         = Champion{526, "Rell", "Rell.png", "champion3.png", 336, 48, 48, 48}
	ChampionRenata       = Champion{888, "Renata Glasc", "Renata.png", "champion3.png", 384, 48, 48, 48}
	ChampionRenekton     = Champion{58, "Renekton", "Renekton.png", "champion3.png", 432, 48, 48, 48}
	ChampionRengar       = Champion{107, "Rengar", "Rengar.png", "champion3.png", 0, 96, 48, 48}
	ChampionRiven        = Champion{92, "Riven", "Riven.png", "champion3.png", 48, 96, 48, 48}
	ChampionRumble       = Champion{68, "Rumble", "Rumble.png", "champion3.png", 96, 96, 48, 48}
	ChampionRyze         = Champion{13, "Ryze", "Ryze.png", "champion3.png", 144, 96, 48, 48}
	ChampionSamira       = Champion{360, "Samira", "Samira.png", "champion3.png", 192, 96, 48, 48}
	ChampionSejuani      = Champion{113, "Sejuani", "Sejuani.png", "champion3.png", 240, 96, 48, 48}
	ChampionSenna        = Champion{235, "Senna", "Senna.png", "champion3.png", 288, 96, 48, 48}
	ChampionSeraphine    = Champion{147, "Seraphine", "Seraphine.png", "champion3.png", 336, 96, 48, 48}
	ChampionSett         = Champion{875, "Sett", "Sett.png", "champion3.png", 384, 96, 48, 48}
	ChampionShaco        = Champion{35, "Shaco", "Shaco.png", "champion3.png", 432, 96, 48, 48}
	ChampionShen         = Champion{98, "Shen", "Shen.png", "champion4.png", 0, 0, 48, 48}
	ChampionShyvana      = Champion{102, "Shyvana", "Shyvana.png", "champion4.png", 48, 0, 48, 48}
	ChampionSinged       = Champion{27, "Singed", "Singed.png", "champion4.png", 96, 0, 48, 48}
	ChampionSion         = Champion{14, "Sion", "Sion.png", "champion4.png", 144, 0, 48, 48}
	ChampionSivir        = Champion{15, "Sivir", "Sivir.png", "champion4.png", 192, 0, 48, 48}
	ChampionSkarner      = Champion{72, "Skarner", "Skarner.png", "champion4.png", 240, 0, 48, 48}
	ChampionSmolder      = Champion{901, "Smolder", "Smolder.png", "champion4.png", 288, 0, 48, 48}
	ChampionSona         = Champion{37, "Sona", "Sona.png", "champion4.png", 336, 0, 48, 48}
	ChampionSoraka       = Champion{16, "Soraka", "Soraka.png", "champion4.png", 384, 0, 48, 48}
	ChampionSwain        = Champion{50, "Swain", "Swain.png", "champion4.png", 432, 0, 48, 48}
	ChampionSylas        = Champion{517, "Sylas", "Sylas.png", "champion4.png", 0, 48, 48, 48}
	ChampionSyndra       = Champion{134, "Syndra", "Syndra.png", "champion4.png", 48, 48, 48, 48}
	ChampionTahmKench    = Champion{223, "Tahm Kench", "TahmKench.png", "champion4.png", 96, 48, 48, 48}
	ChampionTaliyah      = Champion{163, "Taliyah", "Taliyah.png", "champion4.png", 144, 48, 48, 48}
	ChampionTalon        = Champion{91, "Talon", "Talon.png", "champion4.png", 192, 48, 48, 48}
	ChampionTaric        = Champion{44, "Taric", "Taric.png", "champion4.png", 240, 48, 48, 48}
	ChampionTeemo        = Champion{17, "Teemo", "Teemo.png", "champion4.png", 288, 48, 48, 48}
	ChampionThresh       = Champion{412, "Thresh", "Thresh.png", "champion4.png", 336, 48, 48, 48}
	ChampionTristana     = Champion{18, "Tristana", "Tristana.png", "champion4.png", 384, 48, 48, 48}
	ChampionTrundle      = Champion{48, "Trundle", "Trundle.png", "champion4.png", 432, 48, 48, 48}
	ChampionTryndamere   = Champion{23, "Tryndamere", "Tryndamere.png", "champion4.png", 0, 96, 48, 48}
	ChampionTwistedFate  = Champion{4, "Twisted Fate", "TwistedFate.png", "champion4.png", 48, 96, 48, 48}
	ChampionTwitch       = Champion{29, "Twitch", "Twitch.png", "champion4.png", 96, 96, 48, 48}
	ChampionUdyr         = Champion{77, "Udyr", "Udyr.png", "champion4.png", 144, 96, 48, 48}
	ChampionUrgot        = Champion{6, "Urgot", "Urgot.png", "champion4.png", 192, 96, 48, 48}
	ChampionVarus        = Champion{110, "Varus", "Varus.png", "champion4.png", 240, 96, 48, 48}
	ChampionVayne        = Champion{67, "Vayne", "Vayne.png", "champion4.png", 288, 96, 48, 48}
	ChampionVeigar       = Champion{45, "Veigar", "Veigar.png", "champion4.png", 336, 96, 48, 48}
	ChampionVelkoz       = Champion{161, "Vel'Koz", "Velkoz.png", "champion4.png", 384, 96, 48, 48}
	ChampionVex          = Champion{711, "Vex", "Vex.png", "champion4.png", 432, 96, 48, 48}
	ChampionVi           = Champion{254, "Vi", "Vi.png", "champion5.png", 0, 0, 48, 48}
	ChampionViego        = Champion{234, "Viego", "Viego.png", "champion5.png", 48, 0, 48, 48}
	ChampionViktor       = Champion{112, "Viktor", "Viktor.png", "champion5.png", 96, 0, 48, 48}
	ChampionVladimir     = Champion{8, "Vladimir", "Vladimir.png", "champion5.png", 144, 0, 48, 48}
	ChampionVolibear     = Champion{106, "Volibear", "Volibear.png", "champion5.png", 192, 0, 48, 48}
	ChampionWarwick      = Champion{19, "Warwick", "Warwick.png", "champion5.png", 240, 0, 48, 48}
	ChampionXayah        = Champion{498, "Xayah", "Xayah.png", "champion5.png", 288, 0, 48, 48}
	ChampionXerath       = Champion{101, "Xerath", "Xerath.png", "champion5.png", 336, 0, 48, 48}
	ChampionXinZhao      = Champion{5, "Xin Zhao", "XinZhao.png", "champion5.png", 384, 0, 48, 48}
	ChampionYasuo        = Champion{157, "Yasuo", "Yasuo.png", "champion5.png", 432, 0, 48, 48}
	ChampionYone         = Champion{777, "Yone", "Yone.png", "champion5.png", 0, 48, 48, 48}
	ChampionYorick       = Champion{83, "Yorick", "Yorick.png", "champion5.png", 48, 48, 48, 48}
	ChampionYuumi        = Champion{350, "Yuumi", "Yuumi.png", "champion5.png", 96, 48, 48, 48}
	ChampionZac          = Champion{154, "Zac", "Zac.png", "champion5.png", 144, 48, 48, 48}
	ChampionZed          = Champion{238, "Zed", "Zed.png", "champion5.png", 192, 48, 48, 48}
	ChampionZeri         = Champion{221, "Zeri", "Zeri.png", "champion5.png", 240, 48, 48, 48}
	ChampionZiggs        = Champion{115, "Ziggs", "Ziggs.png", "champion5.png", 288, 48, 48, 48}
	ChampionZilean       = Champion{26, "Zilean", "Zilean.png", "champion5.png", 336, 48, 48, 48}
	ChampionZoe          = Champion{142, "Zoe", "Zoe.png", "champion5.png", 384, 48, 48, 48}
	ChampionZyra         = Champion{143, "Zyra", "Zyra.png", "champion5.png", 432, 48, 48, 48}
)
