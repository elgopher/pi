// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package pifont provides functionality for rendering text
// using bitmap fonts.
package pifont

import (
	"github.com/elgopher/pi"
)

// Sheet is a character sheet used for rendering text.
type Sheet struct {
	Chars   map[rune]pi.Sprite
	Height  int
	FgColor pi.Color // font color on sprites
	BgColor pi.Color // background color on sprites
}

var intermediateCanvas pi.Canvas // text is first rendered here to change its color

var prevFgColorTable [pi.MaxColors]pi.Color
var prevBgColorTable [pi.MaxColors]pi.Color

// Print draws text using the current draw color.
//
// Returns the x, y position where you can continue writing text.
func (s Sheet) Print(str string, x, y int) (currentX, currentY int) {
	originalDrawTarget := pi.DrawTarget()
	if intermediateCanvas.W() != originalDrawTarget.W() || intermediateCanvas.H() != originalDrawTarget.H() {
		intermediateCanvas = pi.NewCanvas(originalDrawTarget.W(), originalDrawTarget.H())
	}

	currentColor := pi.GetColor()

	copy(prevFgColorTable[:], pi.ColorTables[0][s.FgColor][:])
	copy(prevBgColorTable[:], pi.ColorTables[0][s.BgColor][:])
	intermediateCanvas.Clear(s.BgColor)
	pi.Pal(s.FgColor, currentColor)
	pi.Palt(s.BgColor, true)
	pi.SetDrawTarget(intermediateCanvas)

	// first draw text in selected color on intermediateCanvas
	currentX, currentY = s.PrintOriginal(str, x, y)

	copy(pi.ColorTables[0][s.FgColor][:], prevFgColorTable[:])

	// now copy text in target color on original draw target
	coloredText := pi.Sprite{
		Area:   pi.Area[int]{X: x - pi.Camera.X, Y: y - pi.Camera.Y, W: currentX - x, H: currentY - y + s.Height},
		Source: intermediateCanvas,
	}
	pi.SetDrawTarget(originalDrawTarget)
	pi.Spr(coloredText, x, y)

	copy(pi.ColorTables[0][s.BgColor][:], prevBgColorTable[:])

	return
}

// PrintOriginal prints the text using its original colors.
func (s Sheet) PrintOriginal(str string, x, y int) (maxX, currentY int) {
	maxX = x
	currentX := x
	currentY = y
	for _, r := range str {
		if r == '\n' {
			currentX = x
			currentY += s.Height
			continue
		}
		sprite := s.Chars[r]
		pi.Spr(sprite, currentX, currentY)
		currentX += sprite.W
		maxX = max(maxX, currentX)
	}

	return
}

// PrintStroked prints the text with a stroke effect.
//
// The text is drawn using the specified foreground and stroke colors.
func (s Sheet) PrintStroked(text string, x, y int, fgColor, strokeColor pi.Color) (currentX, currentY int) {
	prevColor := pi.SetColor(strokeColor)
	for l := y - 1; l <= y+1; l++ {
		s.Print(text, x-1, l)
		s.Print(text, x, l)
		s.Print(text, x+1, l)
	}

	pi.SetColor(fgColor)
	currentX, currentY = s.Print(text, x, y)

	pi.SetColor(prevColor)

	return
}

// Size returns the dimensions of the text without rendering it to the draw target.
func (s Sheet) Size(text string) (width, height int) {
	originalDrawTarget := pi.SetDrawTarget(intermediateCanvas)
	defer pi.SetDrawTarget(originalDrawTarget)

	return s.PrintOriginal(text, 0, 0)
}
