package ddragon_test

import (
	"log"
	"testing"

	"github.com/rank1zen/yujin/internal/ddragon"
)

func TestTest(t *testing.T) {
	u := ddragon.GetRuneIconUrl(8005)
	log.Print(u)
}
