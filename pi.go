// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package pi provides API to develop retro games.
package pi

import (
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"time"
)

const (
	defaultSpriteSheetWidth  = 128
	defaultSpriteSheetHeight = 128
	defaultScreenWidth       = 128
	defaultScreenHeight      = 128
)

// User parameters. Will be used during Boot (and Run).
var (
	Update = func() {} // Update is a user provided function executed each frame.
	Draw   = func() {} // Draw is a user provided function executed each frame (if he can).

	Resources fs.ReadFileFS // Resources contains files like sprite-sheet.png

	// SpriteSheetWidth will be used if sprite-sheet.png was not found.
	SpriteSheetWidth = defaultSpriteSheetWidth
	// SpriteSheetHeight will be used if sprite-sheet.png was not found.
	SpriteSheetHeight = defaultSpriteSheetHeight

	// ScreenWidth specifies the width of the screen (in pixels).
	ScreenWidth = defaultScreenWidth
	// ScreenHeight specifies the height of the screen (in pixels).
	ScreenHeight = defaultScreenHeight
)

// Run boots the game, opens the window and run the game loop. It must be
// called from the main thread.
//
// It returns error when something terrible happened during initialization.
func Run() error {
	if err := Boot(); err != nil {
		return fmt.Errorf("booting game failed: %w", err)
	}

	timeStarted = time.Now()

	return run()
}

// RunOrPanic does the same as Run, but panics instead of returning an error.
//
// Useful for writing unit tests and quick and dirty prototypes. Do not use on production ;)
func RunOrPanic() {
	if err := Run(); err != nil {
		panic(fmt.Sprintf("Something terrible happened! Pi cannot be run: %v\n", err))
	}
}

// Reset resets all user parameters to default values. Useful in unit tests.
func Reset() {
	Update = nil
	Draw = nil
	Resources = nil
	SpriteSheetWidth = defaultSpriteSheetWidth
	SpriteSheetHeight = defaultSpriteSheetHeight
	ScreenWidth = defaultScreenWidth
	ScreenHeight = defaultScreenHeight
	Color = 6
}

// Boot initializes the engine based on user parameters such as ScreenWidth and ScreenHeight.
// It loads the resources like sprite-sheet.png.
func Boot() error {
	if SpriteSheetWidth%8 != 0 || SpriteSheetWidth == 0 {
		return fmt.Errorf("sprite sheet width %d is not a multiplcation of 8", SpriteSheetWidth)
	}
	if SpriteSheetHeight%8 != 0 || SpriteSheetHeight == 0 {
		return fmt.Errorf("sprite sheet height %d is not a multiplcation of 8", SpriteSheetHeight)
	}
	ssWidth = SpriteSheetWidth
	ssHeight = SpriteSheetHeight
	numberOfSprites = (ssWidth * ssHeight) / (SpriteWidth * SpriteHeight)

	if Resources == nil {
		Resources = embed.FS{}
	}

	if err := loadResources(Resources); err != nil {
		return err
	}

	spritesInLine = ssWidth / SpriteWidth
	spritesRows = ssHeight / SpriteHeight

	scrWidth = ScreenWidth
	scrHeight = ScreenHeight
	screenSize := scrWidth * scrHeight
	ScreenData = make([]byte, screenSize)
	zeroScreenData = make([]byte, screenSize)
	lineOfScreenWidth = make([]byte, scrWidth)

	Clip(0, 0, scrWidth, scrHeight)
	Camera(0, 0)

	return nil
}

// BootOrPanic does the same as Boot, but panics instead of returning an error.
//
// Useful for writing unit tests and quick and dirty prototypes. Do not use on production ;)
func BootOrPanic() {
	if err := Boot(); err != nil {
		panic("init failed " + err.Error())
	}
}

var gameLoopStopped bool

// Stop will stop the game loop after Update is finished, but before Draw.
// For now the entire app will be closed, but later it will show a dev console instead,
// where developer will be able to resume the game.
func Stop() {
	gameLoopStopped = true
}
