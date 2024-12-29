package riot

import (
	"time"

	"github.com/rank1zen/yujin/internal/riotclient"
)

type Riot struct {
	client *riotclient.Client
}

func riotUnixToDate(ts int64) time.Time {
	return time.Unix(ts/1000, 0)
}

func riotDurationToInterval(dur int) time.Duration {
	return time.Duration(dur) * time.Second
}
