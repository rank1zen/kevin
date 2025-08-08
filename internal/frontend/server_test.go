package frontend_test

import (
	"testing"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/stretchr/testify/assert"
)

func TestParseRiotID(t *testing.T) {
	for _, test := range []struct{
		TestName string
		RiotID string
		Err error
		Name string
		Tag string
	}{
		{
			TestName: "expects orrange#NA1",
			RiotID:   "orrange-NA1",
			Name:     "orrange",
			Tag:      "NA1",
		},
		{
			TestName: "expects invalid error for missing tag",
			RiotID:   "orrange-",
			Err:      frontend.ErrInvalidRiotID,
		},
		{
			TestName: "expects invalid error for missing tag",
			RiotID:   "orrange",
			Err:      frontend.ErrInvalidRiotID,
		},
		{
			TestName: "expects invalid  error for two seperators",
			RiotID:   "orrange-NA1-NA1",
			Err:      frontend.ErrInvalidRiotID,
		},
	} {
		t.Run(
			test.TestName,
			func(t *testing.T) {
				name, tag, err := frontend.ParseRiotID(test.RiotID)
				if test.Err != nil {
					assert.ErrorIs(t, err, test.Err)
				} else {
					assert.Equal(t, test.Name, name)
					assert.Equal(t, test.Tag, tag)
				}
			},
		)
	}
}
