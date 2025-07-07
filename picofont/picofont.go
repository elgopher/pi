// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package picofont provides the Pico-8 font created by Zep.
//
// The font is available under the CC0 license:
// https://creativecommons.org/publicdomain/zero/1.0/
package picofont

import (
	_ "embed"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/pifont"
)

//go:embed "font.png"
var fontPng []byte

// Sheet provides 256 characters.
// The first 128 characters (0–127) have 4px width,
// while the last 128 characters (128–255) have 8px.
// See font.png for details.
var Sheet = pifont.Sheet{
	Height:  8,
	FgColor: 1,
	BgColor: 0,
}

// Print writes text on the screen using the Pico-8 font.
//
// Returns the x, y position where you can continue writing text.
func Print(text string, x, y int) (currentX, currentY int) {
	return Sheet.Print(text, x, y)
}

func init() {
	prev0, prev1 := pi.Palette[0], pi.Palette[1]
	defer func() {
		pi.Palette[0], pi.Palette[1] = prev0, prev1
	}()

	pi.Palette[0] = 0x000000
	pi.Palette[1] = 0xFFF1E8
	canvas := pi.DecodeCanvas(fontPng)

	Sheet.Chars = map[rune]pi.Sprite{}
	idx := 0
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			width := 8
			if idx < 128 {
				width = 4
			}
			Sheet.Chars[rune(idx)] = pi.Sprite{
				Area: pi.Area[int]{
					X: x * 8, Y: y * 8, W: width, H: 8,
				},
				Source: canvas,
			}
			idx += 1
		}
	}
}
