// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package piscope provides developer tools.
package piscope

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/piscope/internal"
)

// Toolbar colors
var (
	BgColor pi.Color = 8 // toolbar background color
	FgColor pi.Color = 1 // toolbar foreground color
)

// Start launches the developer tools.
//
// Pressing Ctrl+Shift+I will activate the tools in the game.
//
// When the tools are active, the following keyboard shortcuts are available:
//
//   - spacebar - pause/resume the game
//   - left arrow - show a snapshot of the previous frame
//   - right arrow - show a snapshot of the next frame, or resume the game and stop after one frame
//   - F12 - take a screenshot and save it to a file
//   - Esc or Ctrl+Shift+I - exit the tools
//
// Currently, piscope requires the game to have a resolution of at least
// 128 pixels horizontally and 16 pixels vertically. Additionally, the game's
// palette must use at least 2 colors.
func Start() {
	internal.Start(&BgColor, &FgColor)
}
