// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package piebiten enables running your game using the [Ebitengine] backend.
//
// Ebitengine is a cross-platform game engine that supports Windows, macOS,
// Linux, FreeBSD, web browsers, Android, iOS, and even Nintendo Switch.
//
// To launch your game, use [Run] or [RunOrErr].
//
// This package also provides advanced functions for integrating Pi
// with your own Ebitengine-based game, such as [CopyCanvasToEbitenImage].
//
// [Ebitengine]: https://ebitengine.org
package piebiten

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/piebiten/internal"
)

// RememberWindow determines whether the game should open
// at its last window position, size, and monitor when set to true
var RememberWindow = false

// Run starts the Ebitengine backend. It panics if something goes wrong.
//
// If you want to handle errors gracefully, use [RunOrErr] instead.
//
// This function must be called from the first goroutine (the main thread).
func Run() {
	if err := RunOrErr(); err != nil {
		panic("piebiten.Run failed: " + err.Error())
	}
}

// RunOrErr starts the Ebitengine backend and returns an error if something goes wrong.
//
// This function must be called from the first goroutine (the main thread).
func RunOrErr() error {
	if internal.CurrentGoroutineID() != 1 {
		return errors.New("must be run from main goroutine 1")
	}
	internal.RememberWindow = RememberWindow
	return internal.RunOrErr() //nolint:wrapcheck
}

// CopyCanvasToEbitenImage copies the canvas to dst using the current
// palette in pi.Palette and the palette mapping in pi.PaletteMapping.
func CopyCanvasToEbitenImage(canvas pi.Canvas, dst *ebiten.Image) {
	internal.CopyCanvasToEbitenImage(canvas, dst)
}
