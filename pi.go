// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package pi provides API to develop retro games.
//
// Please note that the entire pi package is not concurrency-safe.
// This means that it is unsafe to run functions and access package
// variables from go-routines started by your code.
package pi

import (
	"embed"
	_ "embed"
	"fmt"
	"io/fs"

	"github.com/elgopher/pi/font"
	"github.com/elgopher/pi/vm"
)

// User parameters. Will be used during Boot (and Run).
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

	// SpriteSheetWidth will be used if sprite-sheet.png was not found.
	SpriteSheetWidth = defaultSpriteSheetWidth
	// SpriteSheetHeight will be used if sprite-sheet.png was not found.
	SpriteSheetHeight = defaultSpriteSheetHeight
	// Palette has all colors available in the game. Up to 256.
	// Palette is taken from loaded sprite sheet (which must be
	// a PNG file with indexed color mode). If sprite-sheet.png was not
	// found, then default 16 color palette is used.
	//
	// Can be freely read and updated. Changes will be visible immediately.
	Palette = defaultPalette

	// ScreenWidth specifies the width of the screen (in pixels).
	ScreenWidth = defaultScreenWidth
	// ScreenHeight specifies the height of the screen (in pixels).
	ScreenHeight = defaultScreenHeight

	Resources fs.ReadFileFS // Resources contains files like sprite-sheet.png, custom-font.png
)

var booted bool

// Run boots the game and runs the game loop using specified backend. For ebitengine
// backend it must be called from the main thread, for example:
//
//	err := Run(ebitengine.Backend)
//
// Run does not boot the game once the game has been booted. Thanks to this,
// the user can call Boot directly and draw to the screen
// before the game loop starts.
//
// It returns error when something terrible happened during initialization.
func Run(runBackend func() error) error {
	if !booted {
		if err := Boot(); err != nil {
			return fmt.Errorf("booting game failed: %w", err)
		}
	}

	return runBackend()
}

// MustRun does the same as Run, but panics instead of returning an error.
//
// Useful for writing unit tests and quick and dirty prototypes. Do not use on production ;)
func MustRun(runBackend func() error) {
	if err := Run(runBackend); err != nil {
		panic(fmt.Sprintf("Something terrible happened! Pi cannot be run: %v\n", err))
	}
}

// Reset resets all user parameters to default values. Useful in unit tests.
func Reset() {
	Update = nil
	Draw = nil
	Resources = nil
	ScreenWidth = defaultScreenWidth
	ScreenHeight = defaultScreenHeight
	SpriteSheetWidth = defaultSpriteSheetWidth
	SpriteSheetHeight = defaultSpriteSheetHeight
	Palette = defaultPalette
}

// Boot initializes the engine based on user parameters such as ScreenWidth and ScreenHeight.
// It loads the resources like sprite-sheet.png.
//
// If sprite-sheet.png was not found in pi.Resources, then empty sprite-sheet is used with
// the size of pi.SpriteSheetWidth * pi.SpriteSheetHeight.
//
// Boot also resets all draw state information like camera position and clipping region.
//
// Boot can be run multiple times. This is useful for writing unit tests.
func Boot() error {
	if err := validateUserParameters(); err != nil {
		return err
	}

	if err := font.Load(systemFontPNG, vm.SystemFont.Data[:]); err != nil {
		return err
	}

	if Resources == nil {
		Resources = embed.FS{}
	}

	if err := loadGameResources(Resources); err != nil {
		return err
	}

	numberOfSprites = (vm.SpriteSheetWidth * vm.SpriteSheetHeight) / (SpriteWidth * SpriteHeight)

	spritesInLine = vm.SpriteSheetWidth / SpriteWidth

	vm.ScreenWidth = ScreenWidth
	vm.ScreenHeight = ScreenHeight
	screenSize := vm.ScreenWidth * vm.ScreenHeight
	vm.ScreenData = make([]byte, screenSize)
	zeroScreenData = make([]byte, screenSize)
	lineOfScreenWidth = make([]byte, vm.ScreenWidth)

	ClipReset()
	CameraReset()
	PaltReset()
	PalReset()

	vm.Update = Update
	vm.Draw = Draw

	booted = true

	return nil
}

func validateUserParameters() error {
	if SpriteSheetWidth%8 != 0 || SpriteSheetWidth == 0 {
		return fmt.Errorf("sprite sheet width %d is not a multiplcation of 8", SpriteSheetWidth)
	}
	if SpriteSheetHeight%8 != 0 || SpriteSheetHeight == 0 {
		return fmt.Errorf("sprite sheet height %d is not a multiplcation of 8", SpriteSheetHeight)
	}

	if ScreenWidth <= 0 {
		return fmt.Errorf("screen width %d is not greather than 0", ScreenWidth)
	}
	if ScreenHeight <= 0 {
		return fmt.Errorf("screen height %d is not greather than 0", ScreenHeight)
	}

	return nil
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

// MustBoot does the same as Boot, but panics instead of returning an error.
//
// Useful for writing unit tests and quick and dirty prototypes. Do not use on production ;)
func MustBoot() {
	if err := Boot(); err != nil {
		panic("init failed " + err.Error())
	}
}

// Stop will stop the game loop after Update or Draw is finished.
// If you are using devtools, the game will be paused. Otherwise, the game
// will be closed.
func Stop() {
	vm.GameLoopStopped = true
}

// Time returns the amount of time since game was run, as a (fractional) number of seconds.
//
// Calling Time() multiple times in the same frame will always return the same result.
func Time() float64 {
	return vm.TimeSeconds
}
