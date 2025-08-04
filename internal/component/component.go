package component

import (
	"context"

	"github.com/a-h/templ"
)

// Component models a [templ.Component]
type Component interface {
	// ToTempl builds the [templ.Component]
	ToTempl(ctx context.Context) templ.Component
}

type ComponentFunc func(ctx context.Context) templ.Component

func (f ComponentFunc) ToTempl(ctx context.Context) templ.Component {
	return f(ctx)
}

// TextSize follow Tailwind text size classes.
type TextSize int

const (
	TextSizeXS TextSize = iota
	TextSizeSM
	TextSizeBase
	TextSizeLG
	TextSizeXL
	TextSize2XL
	TextSize3XL
	TextSize4XL
	TextSize5XL
)
