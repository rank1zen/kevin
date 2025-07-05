package testdata

import (
	"os"

	"github.com/rank1zen/kevin/internal/riot"
)

func a() *riot.Match {
	a := os.DirFS("../testdata")
	a.Open("a")
}
