package frontend

import (
	"fmt"

	"github.com/rank1zen/kevin/internal"
	"github.com/rank1zen/kevin/internal/riot"
)

func fmtInt(x int) string {
	return fmt.Sprintf("%d", x)
}

func fmtInThousands(x int) string {
	return fmt.Sprintf("%.1fk", float32(x) / 1000)
}

func fmtPercentage(x float32) string {
	return fmt.Sprintf("%.0f%%", x*100)
}

func fmtRank(rank *internal.RankDetail) string {
	if rank == nil {
		return "Unranked"
	}

	var r string
	if rank.Tier == riot.TierChallenger || rank.Tier == riot.TierGrandmaster || rank.Tier == riot.TierMaster {
		r = string(rank.Tier)
	} else {
		r = string(rank.Tier) + " " + string(rank.Rank)
	}

	return fmt.Sprintf("%s %dLP • %dW - %dL", r, rank.LP, rank.Wins, rank.Losses)
}
