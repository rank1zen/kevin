package shared

import (
	"github.com/rank1zen/kevin/internal/component"
	"github.com/rank1zen/kevin/internal/ddragon"
)

func NewChampionSprite(championID int, size component.TextSize) component.Sprite {
	champion := ddragon.ChampionMap[championID]
	x, y := champion.X/48, champion.Y/48

	var pixelSize, truePixelSize int
	var bgHeight, bgWidth int
	// HACK: extra cuts the black border around the edges.
	var extra int

	if size == component.TextSizeLG {
		pixelSize = 28
		extra = 2 // TODO: figure out what this should be.
	} else {
		pixelSize = 36
		extra = 3
	}

	truePixelSize = pixelSize + 2*extra
	bgHeight = truePixelSize * 3
	bgWidth = truePixelSize * 10
	if champion.Sprite == "champion5.png" {
		// HACK: the last sprite map is not the usual height.
		bgHeight = truePixelSize * 2
	}

	sprite := component.Sprite{
		SpriteMap: "/static/sprite/" + champion.Sprite,
		BGHeight:  bgHeight,
		BGWidth:   bgWidth,
		Height:    pixelSize,
		Width:     pixelSize,
		X:         x*truePixelSize + extra,
		Y:         y*truePixelSize + extra,
		Round:     component.RoundSizeLG,
	}

	return sprite
}

// Note: size is only for future compatibility, currently not used.
func NewSummonerSpellSprite(summonerID int, size component.TextSize) component.Sprite {
	summ := ddragon.SummonerMap[summonerID]
	x, y := summ.X/48, summ.Y/48

	pixelSize := 16

	sprite := component.Sprite{
		SpriteMap: "/static/sprite/" + "tiny_" + summ.Sprite,
		BGHeight:  pixelSize * 4,
		BGWidth:   pixelSize * 10,
		Height:    pixelSize,
		Width:     pixelSize,
		X:         pixelSize * x,
		Y:         pixelSize * y,
	}

	return sprite
}

// Note: size is only for future compatibility, currently not used.
func NewItemSprite(itemID int, size component.TextSize) component.Sprite {
	item := ddragon.ItemsMap[itemID]

	x, y := item.X/48, item.Y/48

	sprite := component.Sprite{
		SpriteMap: "/static/sprite/" + "small_" + item.Sprite,
		BGHeight:  280,
		BGWidth:   280,
		Height:    28,
		Width:     28,
		X:         x * 28,
		Y:         y * 28,
		Round:     component.RoundSizeMD,
	}

	return sprite
}

// Note: size is only for future compatibility, currently not used.
func NewWardSprite(wardID int, size component.TextSize) component.Sprite {
	item := ddragon.ItemsMap[wardID]

	x, y := item.X/48, item.Y/48

	sprite := component.Sprite{
		SpriteMap: "/static/sprite/" + "small_" + item.Sprite,
		BGHeight:  280,
		BGWidth:   280,
		Height:    28,
		Width:     28,
		X:         x * 28,
		Y:         y * 28,
		Round:     component.RoundSizeFull,
	}

	return sprite
}
