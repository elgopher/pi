// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import "image/color"

const maxColors = 64 // copied because Go does not support cycles

// PaletteMaker automatically assign indexes to all visited colors.
// If color was added before, then it is not added again.
type PaletteMaker[T ~uint32] struct { // generic type used instead of pi.RGB to avoid cyclical import
	colorAlreadyUsed map[color.Color]struct{}
	palette          [maxColors]T
	index            int
}

// Add returns true, when the number of colors has been exceeded
func (p *PaletteMaker[T]) Add(c color.Color) bool {
	if p.colorAlreadyUsed == nil {
		p.colorAlreadyUsed = map[color.Color]struct{}{}
	}

	_, found := p.colorAlreadyUsed[c]
	if found {
		return false
	} else if p.index == maxColors {
		return true
	}

	p.colorAlreadyUsed[c] = struct{}{}

	r, g, b, _ := c.RGBA()
	r, g, b = r&0xFF, g&0xFF, b&0xFF
	rgb := r<<16 + g<<8 + b

	p.palette[p.index] = T(rgb)

	p.index++

	return false
}

func (p *PaletteMaker[T]) Palette() [maxColors]T {
	return p.palette
}
