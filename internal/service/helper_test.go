package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSoloQMatchFilter(t *testing.T) {
	options := soloQMatchFilter(time.Unix(1749596377, 0), time.Unix(1749596377, 0))
	expected := new(int64)
	*expected = 1749596377
	assert.Equal(t, expected, options.StartTime)
}
