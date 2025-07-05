package frontend_test

import (
	"context"
	"os"
	"testing"

	"github.com/rank1zen/kevin/internal/frontend"
	"github.com/stretchr/testify/assert"
)

func TestText(t *testing.T) {
	ctx := context.Background()

	component := frontend.Text{
		S: "Hi bro",
	}

	err := component.Render(ctx, os.Stdout)
	assert.NoError(t, err)
}
