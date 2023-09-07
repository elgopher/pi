// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package pi provides API to develop retro games.
//
// Like other game development engines, Pi runs your game in a loop:
//
//	for {
//	  pi.Update()
//	  pi.Draw()
//	  sleep() // sleep until next frame (30 frames per second)
//	 }
//
// Both pi.Update and pi.Draw functions are provided by you. By default,
// they do nothing. You can set them by using:
//
//	pi.Draw = func() {
//	  pi.Print("HELLO WORLD", 40, 60, 7)
//	}
//
// To run the game please use the ebitengine back-end by calling
// ebitengine.Run or ebitengine.MustRun.
//
// During development, you might want to use dev-tools which provide tools
// for screen inspection and REPL terminal, where you can write Go code
// live when your game is running. To start the game with dev-tools
// please use:
//
//	devtools.MustRun(ebitengine.Run)
//
// Please note that the entire pi package is not concurrency-safe.
// This means that it is unsafe to run functions and access package
// variables from go-routines started by your code.
package pi

import (
	_ "embed"
	"errors"
	"fmt"
	"io/fs"

	"github.com/elgopher/pi/audio"
	"github.com/elgopher/pi/font"
)

// Game-loop function callbacks
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
)

// Camera has camera offset used for all subsequent draw operations.
var Camera Position

// Time returns the amount of time since game was run, as a (fractional) number of seconds
//
// Time is updated each frame.
var Time float64

var GameLoopStopped bool

// Load loads files: sprite-sheet.png, custom-font.png and audio.sfx from resources parameter.
//
// Load looks for images with hard-coded names, eg for sprite-sheet it loads "sprite-sheet.png".
// Your file name must be exactly the same. And it cannot be inside subdirectory.
//
// sprite-sheet.png file must have an indexed color mode. This means that pixels in the sprite-sheet
// file refers to an index of the color defined in a small palette, also attached to the file in a
// form of mapping: index number->RGB color. Image must have an indexed color mode,
// because Pi loads the palette from the sprite-sheet.png file itself. Please use a pixel-art editor
// which supports indexed color mode, such as Aseprite (paid) or LibreSprite (free). Sprite-sheet
// width and height must be multiplication of 8. Each sprite is 8x8 pixels. The maximum number of pixels
// is 65536 (64KB).
//
// custom-font.png must also have an indexed color mode. Color with index 0 is treated as background.
// Any other color as foreground. The size of the image is fixed. It must be 128x128. Each char is 8x8.
// Char 0 is in the top-left corner. Char 1 to the right.
//
// To acquire the resources object, the easiest way is to include the resources in the game binary
// by using go:embed directive:
//
//	package main
//
//	// go:embed sprite-sheet.png
//	var resources embed.FS
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
	Camera.Reset()
	ClipReset()
	Pald.Reset()
	Pal.Reset()
	Palt.Reset()
	systemFont.Data, _ = font.Load(systemFontPNG)
	customFont = defaultCustomFont
	customFont.Data = make([]byte, fontDataSize)
	screen = NewPixMap(defaultScreenWidth, defaultScreenHeight)
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

	if err := loadAudio(resources); err != nil {
		return err
	}

	return nil
}

func loadAudio(resources fs.ReadFileFS) error {
	fileContents, err := resources.ReadFile("audio.sfx")
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	if err = audio.Load(fileContents); err != nil {
		return fmt.Errorf("error loading audio.sfx: %w", err)
	}

	return nil
}

// Stop will stop the game loop after Update or Draw is finished.
// If you are using devtools, the game will be paused. Otherwise, the game
// will be closed.
func Stop() {
	GameLoopStopped = true
}
