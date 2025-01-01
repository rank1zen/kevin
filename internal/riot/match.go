package riot

import (
	"context"

	"github.com/rank1zen/yujin/internal"
)

func (r *Riot) GetMatch(ctx context.Context, id internal.MatchID) (internal.RiotMatch, error) {
	m, err := r.client.GetMatch(ctx, id.String())
	if err != nil {
		return internal.RiotMatch{}, err
	}

	return internal.RiotMatch{
		ID:              internal.MatchID(m.Metadata.MatchId),
		DataVersion:     m.Metadata.DataVersion,
		Patch:           internal.GameVersion(m.Info.GameVersion),
		CreateTimestamp: riotUnixToDate(m.Info.GameCreation),
		StartTimestamp:  riotUnixToDate(m.Info.GameStartTimestamp),
		EndTimestamp:    riotUnixToDate(m.Info.GameEndTimestamp),
		Duration:        riotDurationToInterval(int(m.Info.GameDuration)),
		EndOfGameResult: m.Info.EndOfGameResult,
		// TODO: Participants:    nil,
	}, nil
}
