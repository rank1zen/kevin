package frontend

import (
	"testing"
	"time"
)

func TestGroup(t *testing.T) {
	res := groupTimeByDay(
		[]time.Time{
			time.Date(2025, 9, 4, 23, 0, 0, 0, time.UTC),
			time.Date(2025, 9, 5, 0, 0, 0, 0, time.UTC),
		},
	)
	for _, k := range res {
		print(k)
	}
}
