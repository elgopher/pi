// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"fmt"
	"image/color"
	"math"
)

// ClosestColorPicker uses generic types to avoid importing pi.RGB and pi.Color (no cycles)
type ClosestColorPicker[RGB ~uint32, Color ~uint8] struct {
	Palette [64]RGB
	Cache   map[color.Color]Color
}

func (c *ClosestColorPicker[RGB, Color]) IndexInPalette(rgba color.Color) (Color, error) {
	closestColor, ok := c.Cache[rgba] // accessing the cache 3 million times takes 59% of image decoding time
	if ok {
		return closestColor, nil
	}

	if len(c.Cache) == maxColors {
		return 0, fmt.Errorf("too many colors in the image to decode. The max number is %d", maxColors)
	}

	smallestDistance := math.MaxFloat64

	for i, paletteCol := range c.Palette {
		r, g, b, _ := rgba.RGBA()
		r, g, b = r&0xff, g&0xff, b&0xff
		r2, g2, b2 := uint32(paletteCol>>16&0xff),
			uint32(paletteCol>>8&0xff),
			uint32(paletteCol&0xff)
		if r == r2 && g == g2 && b == b2 {
			// found perfect match. Short circuit
			closestColor = Color(i)
			break
		}
		distance := perceptualColorDistance(r, g, b, r2, g2, b2)

		if distance < smallestDistance {
			smallestDistance = distance
			closestColor = Color(i)
		}
	}

	c.Cache[rgba] = closestColor // without caching the code is extremely slow

	return closestColor, nil
}

func perceptualColorDistance(r1, g1, b1, r2, g2, b2 uint32) float64 {
	rd := float64(r1 - r2)
	gd := float64(g1 - g2)
	bd := float64(b1 - b2)
	return math.Sqrt(0.299*rd*rd + 0.587*gd*gd + 0.114*bd*bd)
}
