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
