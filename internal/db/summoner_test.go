package db

import (
	"context"
	"testing"
	"time"
)

func TestCreateSummonerRecord(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	db := setupDB(t)

	now := time.Now()

}
