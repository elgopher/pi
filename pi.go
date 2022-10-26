// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package pi provides API to develop retro games.
//
// Please note that the entire pi package is not concurrency-safe.
// This means that it is unsafe to run functions and access package
// variables from go-routines started by your code.
package pi

import (
	_ "embed"
	"io/fs"

	"github.com/elgopher/pi/font"
)

// User parameters.
var (
	// Update is a user provided function executed each frame (30 times per second).
	//
	// The purpose of this function is to handle user input, perform calculations, update
	// game state etc. Typically, this function does not draw on screen.
	Update func()

	// Draw is a user provided function executed at most each frame (up to 30 times per second).
	// π may skip calling this function if previous frame took too long.
	//
	// The purpose of this function is to draw on screen.
	Draw func()

	// Palette has all colors available in the game. Up to 256.
	// Palette is taken from loaded sprite sheet (which must be
	// a PNG file with indexed color mode). If sprite-sheet.png was not
	// found, then default 16 color palette is used.
	//
	// Can be freely read and updated. Changes will be visible immediately.
	Palette = defaultPalette
)

var (
	// DrawPalette contains mapping of colors used to replace color with
	// another one for all subsequent drawings.
	//
	// The index of array is original color, the value is color replacement.
	DrawPalette [256]byte

	// DisplayPalette contains mapping of colors used to replace color with
	// another one for the entire screen, at the end of a frame
	//
	// The index of array is original color, the value is color replacement.
	DisplayPalette [256]byte

	// ColorTransparency contains information whether given color is transparent.
	//
	// The index of array is a color number.
	ColorTransparency = defaultTransparency

	// TimeSeconds is the number of seconds since game was started
	TimeSeconds float64

	GameLoopStopped bool
)

// Load loads files like sprite-sheet.png, custom-font.png
func Load(resources fs.ReadFileFS) {
	if resources == nil {
		return
	}

	if err := loadGameResources(resources); err != nil {
		panic(err)
	}
}

func Reset() {
	Update = nil
	Draw = nil
	CameraReset()
	ClipReset()
	PalReset()
	PaltReset()
	systemFont.Data, _ = font.Load(systemFontPNG)
	customFont = defaultCustomFont
	customFont.Data = make([]byte, fontDataSize)
	screen = newScreen(defaultScreenWidth, defaultScreenHeight)
	sprSheet = newSpriteSheet(defaultSpriteSheetWidth, defaultSpriteSheetHeight)
	Palette = defaultPalette
}

func loadGameResources(resources fs.ReadFileFS) error {
	if err := loadSpriteSheet(resources); err != nil {
		return err
	}

	if err := loadCustomFont(resources); err != nil {
		return err
	}

	return nil
}

// Stop will stop the game loop after Update or Draw is finished.
// If you are using devtools, the game will be paused. Otherwise, the game
// will be closed.
func Stop() {
	GameLoopStopped = true
}

// Time returns the amount of time since game was run, as a (fractional) number of seconds.
//
// Calling Time() multiple times in the same frame will always return the same result.
func Time() float64 {
	return TimeSeconds
}
