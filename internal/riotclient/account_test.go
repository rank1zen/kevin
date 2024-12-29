package riotclient

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccountByRiotId(t *testing.T) {
	client := &Client{}

	ctx := context.Background()

	m, err := client.AccountGetByRiotId(ctx, "orrange", "na1")
	if assert.NoError(t, err) {
		assert.Equal(t, "orrange", m.GameName)
		assert.Equal(t, "NA1", m.TagLine)
		assert.Equal(t, "0bEBr8VSevIGuIyJRLw12BKo3Li4mxvHpy_7l94W6p5SRrpv00U3cWAx7hC4hqf_efY8J4omElP9-Q", m.Puuid)
	}
}
