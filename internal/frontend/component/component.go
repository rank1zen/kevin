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
	TextSizeXS   TextSize = iota // 16px
	TextSizeSM                   // 20 px
	TextSizeBase                 // 24px
	TextSizeLG                   // 28px
	TextSizeXL                   // 32px
	TextSize2XL                  // 36px
	TextSize3XL                  // 40px
	TextSize4XL                  // 44px
	TextSize5XL                  // 48px
)

// RoundSize follow Tailwind rounding size classes.
type RoundSize int

const (
	RoundSizeNone RoundSize = iota
	RoundSizeXS
	RoundSizeSM
	RoundSizeMD
	RoundSizeLG
	RoundSizeXL
	RoundSize2XL
	RoundSize3XL
	RoundSize4XL
	RoundSizeFull
)

func (s RoundSize) class() string {
	classes := map[RoundSize]string{
		RoundSizeNone: "rounded-none",
		RoundSizeXS:   "rounded-xs",
		RoundSizeSM:   "rounded-sm",
		RoundSizeMD:   "rounded-md",
		RoundSizeLG:   "rounded-lg",
		RoundSizeXL:   "rounded-xl",
		RoundSize2XL:  "rounded-2xl",
		RoundSize3XL:  "rounded-3xl",
		RoundSize4XL:  "rounded-4xl",
		RoundSizeFull: "rounded-full",
	}

	return classes[s]
}
