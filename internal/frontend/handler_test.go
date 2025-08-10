package frontend_test

import (
	"testing"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/stretchr/testify/assert"
)

func TestGetNameTag(t *testing.T) {
	tests := []struct {
		TestName string
		Query    string
		Name     string
		Tag      string
	}{
		{
			TestName: "expects empty tag",
			Query:    "orrange#",
			Name:     "orrange",
			Tag:      "",
		},
		{
			TestName: "expects both name and tag",
			Query:    "orrange#123",
			Name:     "orrange",
			Tag:      "123",
		},
		{
			TestName: "expects empty name",
			Query:    "#123",
			Name:     "",
			Tag:      "123",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.TestName,
			func(t *testing.T) {
				name, tag := frontend.GetNameTag(tt.Query)
				assert.Equal(t, tt.Name, name)
				assert.Equal(t, tt.Tag, tag)
			},
		)
	}
}
