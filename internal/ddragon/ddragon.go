package ddragon

import "fmt"

type DDragon struct {
	base string
}

func New(base string) *DDragon {
	return &DDragon{base: base}
}

func (d *DDragon) GetChampionImage(id int) string {
	if champion, ok := ChampionMap[id]; ok {
		return fmt.Sprintf("%s/img/champion/%s", d.base, champion.FullIcon)
	}

	return fmt.Sprintf("%s/img/champion/%s", d.base, ChampionAatrox.FullIcon)
}

func (d *DDragon) GetSummonerSpellImage(id int) string {
	if spell, ok := SummonerMap[id]; ok {
		return fmt.Sprintf("%s/img/spell/%s", d.base, spell.FullIcon)
	}

	return fmt.Sprintf("%s/img/spell/%s", d.base, SummonerMap[21].FullIcon)
}

// GetItemImage returns the image URL for the given item ID. If the item with id
// is not found, it returns an empty string.
func (d *DDragon) GetItemImage(id int) string {
	if item, ok := ItemsMap[id]; ok {
		return fmt.Sprintf("%s/img/item/%s", d.base, item.FullIcon)
	}

	return ""
}
