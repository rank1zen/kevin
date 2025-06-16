package ddragon

type ChampionID int

type Champion struct {
	ID         ChampionID
	Name       string
	FullIcon   string
	Sprite     string
	X, Y, W, H int
}

var (
	Aatrox       = Champion{266, "Aatrox", "Aatrox.png", "champion0.png", 0, 0, 48, 48}
	Ahri         = Champion{103, "Ahri", "Ahri.png", "champion0.png", 48, 0, 48, 48}
	Akali        = Champion{84, "Akali", "Akali.png", "champion0.png", 96, 0, 48, 48}
	Akshan       = Champion{166, "Akshan", "Akshan.png", "champion0.png", 144, 0, 48, 48}
	Alistar      = Champion{12, "Alistar", "Alistar.png", "champion0.png", 192, 0, 48, 48}
	Ambessa      = Champion{799, "Ambessa", "Ambessa.png", "champion0.png", 240, 0, 48, 48}
	Amumu        = Champion{32, "Amumu", "Amumu.png", "champion0.png", 288, 0, 48, 48}
	Anivia       = Champion{34, "Anivia", "Anivia.png", "champion0.png", 336, 0, 48, 48}
	Annie        = Champion{1, "Annie", "Annie.png", "champion0.png", 384, 0, 48, 48}
	Aphelios     = Champion{523, "Aphelios", "Aphelios.png", "champion0.png", 432, 0, 48, 48}
	Ashe         = Champion{22, "Ashe", "Ashe.png", "champion0.png", 0, 48, 48, 48}
	AurelionSol  = Champion{136, "Aurelion Sol", "AurelionSol.png", "champion0.png", 48, 48, 48, 48}
	Aurora       = Champion{893, "Aurora", "Aurora.png", "champion0.png", 96, 48, 48, 48}
	Azir         = Champion{268, "Azir", "Azir.png", "champion0.png", 144, 48, 48, 48}
	Bard         = Champion{432, "Bard", "Bard.png", "champion0.png", 192, 48, 48, 48}
	Belveth      = Champion{200, "Bel'Veth", "Belveth.png", "champion0.png", 240, 48, 48, 48}
	Blitzcrank   = Champion{53, "Blitzcrank", "Blitzcrank.png", "champion0.png", 288, 48, 48, 48}
	Brand        = Champion{63, "Brand", "Brand.png", "champion0.png", 336, 48, 48, 48}
	Braum        = Champion{201, "Braum", "Braum.png", "champion0.png", 384, 48, 48, 48}
	Briar        = Champion{233, "Briar", "Briar.png", "champion0.png", 432, 48, 48, 48}
	Caitlyn      = Champion{51, "Caitlyn", "Caitlyn.png", "champion0.png", 0, 96, 48, 48}
	Camille      = Champion{164, "Camille", "Camille.png", "champion0.png", 48, 96, 48, 48}
	Cassiopeia   = Champion{69, "Cassiopeia", "Cassiopeia.png", "champion0.png", 96, 96, 48, 48}
	Chogath      = Champion{31, "Cho'Gath", "Chogath.png", "champion0.png", 144, 96, 48, 48}
	Corki        = Champion{42, "Corki", "Corki.png", "champion0.png", 192, 96, 48, 48}
	Darius       = Champion{122, "Darius", "Darius.png", "champion0.png", 240, 96, 48, 48}
	Diana        = Champion{131, "Diana", "Diana.png", "champion0.png", 288, 96, 48, 48}
	Draven       = Champion{119, "Draven", "Draven.png", "champion0.png", 336, 96, 48, 48}
	DrMundo      = Champion{36, "Dr. Mundo", "DrMundo.png", "champion0.png", 384, 96, 48, 48}
	Ekko         = Champion{245, "Ekko", "Ekko.png", "champion0.png", 432, 96, 48, 48}
	Elise        = Champion{60, "Elise", "Elise.png", "champion1.png", 0, 0, 48, 48}
	Evelynn      = Champion{28, "Evelynn", "Evelynn.png", "champion1.png", 48, 0, 48, 48}
	Ezreal       = Champion{81, "Ezreal", "Ezreal.png", "champion1.png", 96, 0, 48, 48}
	Fiddlesticks = Champion{9, "Fiddlesticks", "Fiddlesticks.png", "champion1.png", 144, 0, 48, 48}
	Fiora        = Champion{114, "Fiora", "Fiora.png", "champion1.png", 192, 0, 48, 48}
	Fizz         = Champion{105, "Fizz", "Fizz.png", "champion1.png", 240, 0, 48, 48}
	Galio        = Champion{3, "Galio", "Galio.png", "champion1.png", 288, 0, 48, 48}
	Gangplank    = Champion{41, "Gangplank", "Gangplank.png", "champion1.png", 336, 0, 48, 48}
	Garen        = Champion{86, "Garen", "Garen.png", "champion1.png", 384, 0, 48, 48}
	Gnar         = Champion{150, "Gnar", "Gnar.png", "champion1.png", 432, 0, 48, 48}
	Gragas       = Champion{79, "Gragas", "Gragas.png", "champion1.png", 0, 48, 48, 48}
	Graves       = Champion{104, "Graves", "Graves.png", "champion1.png", 48, 48, 48, 48}
	Gwen         = Champion{887, "Gwen", "Gwen.png", "champion1.png", 96, 48, 48, 48}
	Hecarim      = Champion{120, "Hecarim", "Hecarim.png", "champion1.png", 144, 48, 48, 48}
	Heimerdinger = Champion{74, "Heimerdinger", "Heimerdinger.png", "champion1.png", 192, 48, 48, 48}
	Hwei         = Champion{910, "Hwei", "Hwei.png", "champion1.png", 240, 48, 48, 48}
	Illaoi       = Champion{420, "Illaoi", "Illaoi.png", "champion1.png", 288, 48, 48, 48}
	Irelia       = Champion{39, "Irelia", "Irelia.png", "champion1.png", 336, 48, 48, 48}
	Ivern        = Champion{427, "Ivern", "Ivern.png", "champion1.png", 384, 48, 48, 48}
	Janna        = Champion{40, "Janna", "Janna.png", "champion1.png", 432, 48, 48, 48}
	JarvanIV     = Champion{59, "Jarvan IV", "JarvanIV.png", "champion1.png", 0, 96, 48, 48}
	Jax          = Champion{24, "Jax", "Jax.png", "champion1.png", 48, 96, 48, 48}
	Jayce        = Champion{126, "Jayce", "Jayce.png", "champion1.png", 96, 96, 48, 48}
	Jhin         = Champion{202, "Jhin", "Jhin.png", "champion1.png", 144, 96, 48, 48}
	Jinx         = Champion{222, "Jinx", "Jinx.png", "champion1.png", 192, 96, 48, 48}
	Kaisa        = Champion{145, "Kai'Sa", "Kaisa.png", "champion1.png", 240, 96, 48, 48}
	Kalista      = Champion{429, "Kalista", "Kalista.png", "champion1.png", 288, 96, 48, 48}
	Karma        = Champion{43, "Karma", "Karma.png", "champion1.png", 336, 96, 48, 48}
	Karthus      = Champion{30, "Karthus", "Karthus.png", "champion1.png", 384, 96, 48, 48}
	Kassadin     = Champion{38, "Kassadin", "Kassadin.png", "champion1.png", 432, 96, 48, 48}
	Katarina     = Champion{55, "Katarina", "Katarina.png", "champion2.png", 0, 0, 48, 48}
	Kayle        = Champion{10, "Kayle", "Kayle.png", "champion2.png", 48, 0, 48, 48}
	Kayn         = Champion{141, "Kayn", "Kayn.png", "champion2.png", 96, 0, 48, 48}
	Kennen       = Champion{85, "Kennen", "Kennen.png", "champion2.png", 144, 0, 48, 48}
	Khazix       = Champion{121, "Kha'Zix", "Khazix.png", "champion2.png", 192, 0, 48, 48}
	Kindred      = Champion{203, "Kindred", "Kindred.png", "champion2.png", 240, 0, 48, 48}
	Kled         = Champion{240, "Kled", "Kled.png", "champion2.png", 288, 0, 48, 48}
	KogMaw       = Champion{96, "Kog'Maw", "KogMaw.png", "champion2.png", 336, 0, 48, 48}
	KSante       = Champion{897, "K'Sante", "KSante.png", "champion2.png", 384, 0, 48, 48}
	Leblanc      = Champion{7, "LeBlanc", "Leblanc.png", "champion2.png", 432, 0, 48, 48}
	LeeSin       = Champion{64, "Lee Sin", "LeeSin.png", "champion2.png", 0, 48, 48, 48}
	Leona        = Champion{89, "Leona", "Leona.png", "champion2.png", 48, 48, 48, 48}
	Lillia       = Champion{876, "Lillia", "Lillia.png", "champion2.png", 96, 48, 48, 48}
	Lissandra    = Champion{127, "Lissandra", "Lissandra.png", "champion2.png", 144, 48, 48, 48}
	Lucian       = Champion{236, "Lucian", "Lucian.png", "champion2.png", 192, 48, 48, 48}
	Lulu         = Champion{117, "Lulu", "Lulu.png", "champion2.png", 240, 48, 48, 48}
	Lux          = Champion{99, "Lux", "Lux.png", "champion2.png", 288, 48, 48, 48}
	Malphite     = Champion{54, "Malphite", "Malphite.png", "champion2.png", 336, 48, 48, 48}
	Malzahar     = Champion{90, "Malzahar", "Malzahar.png", "champion2.png", 384, 48, 48, 48}
	Maokai       = Champion{57, "Maokai", "Maokai.png", "champion2.png", 432, 48, 48, 48}
	MasterYi     = Champion{11, "Master Yi", "MasterYi.png", "champion2.png", 0, 96, 48, 48}
	Mel          = Champion{800, "Mel", "Mel.png", "champion2.png", 48, 96, 48, 48}
	Milio        = Champion{902, "Milio", "Milio.png", "champion2.png", 96, 96, 48, 48}
	MissFortune  = Champion{21, "Miss Fortune", "MissFortune.png", "champion2.png", 144, 96, 48, 48}
	MonkeyKing   = Champion{62, "Wukong", "MonkeyKing.png", "champion2.png", 192, 96, 48, 48}
	Mordekaiser  = Champion{82, "Mordekaiser", "Mordekaiser.png", "champion2.png", 240, 96, 48, 48}
	Morgana      = Champion{25, "Morgana", "Morgana.png", "champion2.png", 288, 96, 48, 48}
	Naafiri      = Champion{950, "Naafiri", "Naafiri.png", "champion2.png", 336, 96, 48, 48}
	Nami         = Champion{267, "Nami", "Nami.png", "champion2.png", 384, 96, 48, 48}
	Nasus        = Champion{75, "Nasus", "Nasus.png", "champion2.png", 432, 96, 48, 48}
	Nautilus     = Champion{111, "Nautilus", "Nautilus.png", "champion3.png", 0, 0, 48, 48}
	Neeko        = Champion{518, "Neeko", "Neeko.png", "champion3.png", 48, 0, 48, 48}
	Nidalee      = Champion{76, "Nidalee", "Nidalee.png", "champion3.png", 96, 0, 48, 48}
	Nilah        = Champion{895, "Nilah", "Nilah.png", "champion3.png", 144, 0, 48, 48}
	Nocturne     = Champion{56, "Nocturne", "Nocturne.png", "champion3.png", 192, 0, 48, 48}
	Nunu         = Champion{20, "Nunu & Willump", "Nunu.png", "champion3.png", 240, 0, 48, 48}
	Olaf         = Champion{2, "Olaf", "Olaf.png", "champion3.png", 288, 0, 48, 48}
	Orianna      = Champion{61, "Orianna", "Orianna.png", "champion3.png", 336, 0, 48, 48}
	Ornn         = Champion{516, "Ornn", "Ornn.png", "champion3.png", 384, 0, 48, 48}
	Pantheon     = Champion{80, "Pantheon", "Pantheon.png", "champion3.png", 432, 0, 48, 48}
	Poppy        = Champion{78, "Poppy", "Poppy.png", "champion3.png", 0, 48, 48, 48}
	Pyke         = Champion{555, "Pyke", "Pyke.png", "champion3.png", 48, 48, 48, 48}
	Qiyana       = Champion{246, "Qiyana", "Qiyana.png", "champion3.png", 96, 48, 48, 48}
	Quinn        = Champion{133, "Quinn", "Quinn.png", "champion3.png", 144, 48, 48, 48}
	Rakan        = Champion{497, "Rakan", "Rakan.png", "champion3.png", 192, 48, 48, 48}
	Rammus       = Champion{33, "Rammus", "Rammus.png", "champion3.png", 240, 48, 48, 48}
	RekSai       = Champion{421, "Rek'Sai", "RekSai.png", "champion3.png", 288, 48, 48, 48}
	Rell         = Champion{526, "Rell", "Rell.png", "champion3.png", 336, 48, 48, 48}
	Renata       = Champion{888, "Renata Glasc", "Renata.png", "champion3.png", 384, 48, 48, 48}
	Renekton     = Champion{58, "Renekton", "Renekton.png", "champion3.png", 432, 48, 48, 48}
	Rengar       = Champion{107, "Rengar", "Rengar.png", "champion3.png", 0, 96, 48, 48}
	Riven        = Champion{92, "Riven", "Riven.png", "champion3.png", 48, 96, 48, 48}
	Rumble       = Champion{68, "Rumble", "Rumble.png", "champion3.png", 96, 96, 48, 48}
	Ryze         = Champion{13, "Ryze", "Ryze.png", "champion3.png", 144, 96, 48, 48}
	Samira       = Champion{360, "Samira", "Samira.png", "champion3.png", 192, 96, 48, 48}
	Sejuani      = Champion{113, "Sejuani", "Sejuani.png", "champion3.png", 240, 96, 48, 48}
	Senna        = Champion{235, "Senna", "Senna.png", "champion3.png", 288, 96, 48, 48}
	Seraphine    = Champion{147, "Seraphine", "Seraphine.png", "champion3.png", 336, 96, 48, 48}
	Sett         = Champion{875, "Sett", "Sett.png", "champion3.png", 384, 96, 48, 48}
	Shaco        = Champion{35, "Shaco", "Shaco.png", "champion3.png", 432, 96, 48, 48}
	Shen         = Champion{98, "Shen", "Shen.png", "champion4.png", 0, 0, 48, 48}
	Shyvana      = Champion{102, "Shyvana", "Shyvana.png", "champion4.png", 48, 0, 48, 48}
	Singed       = Champion{27, "Singed", "Singed.png", "champion4.png", 96, 0, 48, 48}
	Sion         = Champion{14, "Sion", "Sion.png", "champion4.png", 144, 0, 48, 48}
	Sivir        = Champion{15, "Sivir", "Sivir.png", "champion4.png", 192, 0, 48, 48}
	Skarner      = Champion{72, "Skarner", "Skarner.png", "champion4.png", 240, 0, 48, 48}
	Smolder      = Champion{901, "Smolder", "Smolder.png", "champion4.png", 288, 0, 48, 48}
	Sona         = Champion{37, "Sona", "Sona.png", "champion4.png", 336, 0, 48, 48}
	Soraka       = Champion{16, "Soraka", "Soraka.png", "champion4.png", 384, 0, 48, 48}
	Swain        = Champion{50, "Swain", "Swain.png", "champion4.png", 432, 0, 48, 48}
	Sylas        = Champion{517, "Sylas", "Sylas.png", "champion4.png", 0, 48, 48, 48}
	Syndra       = Champion{134, "Syndra", "Syndra.png", "champion4.png", 48, 48, 48, 48}
	TahmKench    = Champion{223, "Tahm Kench", "TahmKench.png", "champion4.png", 96, 48, 48, 48}
	Taliyah      = Champion{163, "Taliyah", "Taliyah.png", "champion4.png", 144, 48, 48, 48}
	Talon        = Champion{91, "Talon", "Talon.png", "champion4.png", 192, 48, 48, 48}
	Taric        = Champion{44, "Taric", "Taric.png", "champion4.png", 240, 48, 48, 48}
	Teemo        = Champion{17, "Teemo", "Teemo.png", "champion4.png", 288, 48, 48, 48}
	Thresh       = Champion{412, "Thresh", "Thresh.png", "champion4.png", 336, 48, 48, 48}
	Tristana     = Champion{18, "Tristana", "Tristana.png", "champion4.png", 384, 48, 48, 48}
	Trundle      = Champion{48, "Trundle", "Trundle.png", "champion4.png", 432, 48, 48, 48}
	Tryndamere   = Champion{23, "Tryndamere", "Tryndamere.png", "champion4.png", 0, 96, 48, 48}
	TwistedFate  = Champion{4, "Twisted Fate", "TwistedFate.png", "champion4.png", 48, 96, 48, 48}
	Twitch       = Champion{29, "Twitch", "Twitch.png", "champion4.png", 96, 96, 48, 48}
	Udyr         = Champion{77, "Udyr", "Udyr.png", "champion4.png", 144, 96, 48, 48}
	Urgot        = Champion{6, "Urgot", "Urgot.png", "champion4.png", 192, 96, 48, 48}
	Varus        = Champion{110, "Varus", "Varus.png", "champion4.png", 240, 96, 48, 48}
	Vayne        = Champion{67, "Vayne", "Vayne.png", "champion4.png", 288, 96, 48, 48}
	Veigar       = Champion{45, "Veigar", "Veigar.png", "champion4.png", 336, 96, 48, 48}
	Velkoz       = Champion{161, "Vel'Koz", "Velkoz.png", "champion4.png", 384, 96, 48, 48}
	Vex          = Champion{711, "Vex", "Vex.png", "champion4.png", 432, 96, 48, 48}
	Vi           = Champion{254, "Vi", "Vi.png", "champion5.png", 0, 0, 48, 48}
	Viego        = Champion{234, "Viego", "Viego.png", "champion5.png", 48, 0, 48, 48}
	Viktor       = Champion{112, "Viktor", "Viktor.png", "champion5.png", 96, 0, 48, 48}
	Vladimir     = Champion{8, "Vladimir", "Vladimir.png", "champion5.png", 144, 0, 48, 48}
	Volibear     = Champion{106, "Volibear", "Volibear.png", "champion5.png", 192, 0, 48, 48}
	Warwick      = Champion{19, "Warwick", "Warwick.png", "champion5.png", 240, 0, 48, 48}
	Xayah        = Champion{498, "Xayah", "Xayah.png", "champion5.png", 288, 0, 48, 48}
	Xerath       = Champion{101, "Xerath", "Xerath.png", "champion5.png", 336, 0, 48, 48}
	XinZhao      = Champion{5, "Xin Zhao", "XinZhao.png", "champion5.png", 384, 0, 48, 48}
	Yasuo        = Champion{157, "Yasuo", "Yasuo.png", "champion5.png", 432, 0, 48, 48}
	Yone         = Champion{777, "Yone", "Yone.png", "champion5.png", 0, 48, 48, 48}
	Yorick       = Champion{83, "Yorick", "Yorick.png", "champion5.png", 48, 48, 48, 48}
	Yuumi        = Champion{350, "Yuumi", "Yuumi.png", "champion5.png", 96, 48, 48, 48}
	Zac          = Champion{154, "Zac", "Zac.png", "champion5.png", 144, 48, 48, 48}
	Zed          = Champion{238, "Zed", "Zed.png", "champion5.png", 192, 48, 48, 48}
	Zeri         = Champion{221, "Zeri", "Zeri.png", "champion5.png", 240, 48, 48, 48}
	Ziggs        = Champion{115, "Ziggs", "Ziggs.png", "champion5.png", 288, 48, 48, 48}
	Zilean       = Champion{26, "Zilean", "Zilean.png", "champion5.png", 336, 48, 48, 48}
	Zoe          = Champion{142, "Zoe", "Zoe.png", "champion5.png", 384, 48, 48, 48}
	Zyra         = Champion{143, "Zyra", "Zyra.png", "champion5.png", 432, 48, 48, 48}
)

var ChampionMap = map[ChampionID]Champion{
	266: Aatrox,
	103: Ahri,
	84:  Akali,
	166: Akshan,
	12:  Alistar,
	799: Ambessa,
	32:  Amumu,
	34:  Anivia,
	1:   Annie,
	523: Aphelios,
	22:  Ashe,
	136: AurelionSol,
	893: Aurora,
	268: Azir,
	432: Bard,
	200: Belveth,
	53:  Blitzcrank,
	63:  Brand,
	201: Braum,
	233: Briar,
	51:  Caitlyn,
	164: Camille,
	69:  Cassiopeia,
	31:  Chogath,
	42:  Corki,
	122: Darius,
	131: Diana,
	119: Draven,
	36:  DrMundo,
	245: Ekko,
	60:  Elise,
	28:  Evelynn,
	81:  Ezreal,
	9:   Fiddlesticks,
	114: Fiora,
	105: Fizz,
	3:   Galio,
	41:  Gangplank,
	86:  Garen,
	150: Gnar,
	79:  Gragas,
	104: Graves,
	887: Gwen,
	120: Hecarim,
	74:  Heimerdinger,
	910: Hwei,
	420: Illaoi,
	39:  Irelia,
	427: Ivern,
	40:  Janna,
	59:  JarvanIV,
	24:  Jax,
	126: Jayce,
	202: Jhin,
	222: Jinx,
	145: Kaisa,
	429: Kalista,
	43:  Karma,
	30:  Karthus,
	38:  Kassadin,
	55:  Katarina,
	10:  Kayle,
	141: Kayn,
	85:  Kennen,
	121: Khazix,
	203: Kindred,
	240: Kled,
	96:  KogMaw,
	897: KSante,
	7:   Leblanc,
	64:  LeeSin,
	89:  Leona,
	876: Lillia,
	127: Lissandra,
	236: Lucian,
	117: Lulu,
	99:  Lux,
	54:  Malphite,
	90:  Malzahar,
	57:  Maokai,
	11:  MasterYi,
	800: Mel,
	902: Milio,
	21:  MissFortune,
	62:  MonkeyKing,
	82:  Mordekaiser,
	25:  Morgana,
	950: Naafiri,
	267: Nami,
	75:  Nasus,
	111: Nautilus,
	518: Neeko,
	76:  Nidalee,
	895: Nilah,
	56:  Nocturne,
	20:  Nunu,
	2:   Olaf,
	61:  Orianna,
	516: Ornn,
	80:  Pantheon,
	78:  Poppy,
	555: Pyke,
	246: Qiyana,
	133: Quinn,
	497: Rakan,
	33:  Rammus,
	421: RekSai,
	526: Rell,
	888: Renata,
	58:  Renekton,
	107: Rengar,
	92:  Riven,
	68:  Rumble,
	13:  Ryze,
	360: Samira,
	113: Sejuani,
	235: Senna,
	147: Seraphine,
	875: Sett,
	35:  Shaco,
	98:  Shen,
	102: Shyvana,
	27:  Singed,
	14:  Sion,
	15:  Sivir,
	72:  Skarner,
	901: Smolder,
	37:  Sona,
	16:  Soraka,
	50:  Swain,
	517: Sylas,
	134: Syndra,
	223: TahmKench,
	163: Taliyah,
	91:  Talon,
	44:  Taric,
	17:  Teemo,
	412: Thresh,
	18:  Tristana,
	48:  Trundle,
	23:  Tryndamere,
	4:   TwistedFate,
	29:  Twitch,
	77:  Udyr,
	6:   Urgot,
	110: Varus,
	67:  Vayne,
	45:  Veigar,
	161: Velkoz,
	711: Vex,
	254: Vi,
	234: Viego,
	112: Viktor,
	8:   Vladimir,
	106: Volibear,
	19:  Warwick,
	498: Xayah,
	101: Xerath,
	5:   XinZhao,
	157: Yasuo,
	777: Yone,
	83:  Yorick,
	350: Yuumi,
	154: Zac,
	238: Zed,
	221: Zeri,
	115: Ziggs,
	26:  Zilean,
	142: Zoe,
	143: Zyra,
}
