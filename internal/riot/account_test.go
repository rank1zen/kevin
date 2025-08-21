package riot_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rank1zen/kevin/internal/riot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountServiceGetAccountByRiotID(t *testing.T) {
	ctx := context.Background()

	client, server := MakeTestClient(t, http.StatusOK, "sample/riot/account/v1/accounts/by-riot-id/orrange-NA1.json",
		func(r *http.Request) {
			for _, tc := range []struct {
				Name             string
				Expected, Actual any
			}{
				{
					Name:     "expects correct endpoint",
					Expected: "/riot/account/v1/accounts/by-riot-id/orrange/NA1",
					Actual:   r.URL.Path,
				},
			} {
				t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
			}
		},
	)

	defer server.Close()

	res, err := client.Account.GetAccountByRiotID(ctx, riot.RegionNA1, "orrange", "NA1")
	require.NoError(t, err)

	for _, tc := range []struct {
		Name             string
		Expected, Actual any
	}{
		{
			Name:     "expects correct puuid for orrange",
			Expected: OrrangePUUID,
			Actual:   res.PUUID,
		},
		{
			Name:     "expects correct name for orrange",
			Expected: OrrangePUUID,
			Actual:   res.GameName,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) { assert.Equal(t, tc.Expected, tc.Actual) })
	}
}
