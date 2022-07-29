// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"

	"github.com/elgopher/pi/image"
)

const (
	SpriteWidth, SpriteHeight = 8, 8
)

// Sprite-sheet data
var (
	// SpriteSheetData contains pixel colors for the entire sprite sheet.
	// Each pixel is one byte. It is initialized during pi.Boot.
	//
	// Pixels in the sprite-sheet are organized from left to right,
	// top to bottom. Slice element number 0 has pixel located
	// in the top-left corner. Slice element number 1 has a pixel color
	// on the right and so on.
	//
	// Can be freely read and updated.
	// Useful when you want to use your own functions for pixel manipulation.
	// Pi will panic if you try to change the length of the slice.
	SpriteSheetData []byte

	ssWidth, ssHeight int
	numberOfSprites   int
	spritesInLine     int

	//nolint:govet
	defaultPalette = [256]image.RGB{
		{0, 0, 0},          // 0 - black
		{0x1D, 0x2B, 0x53}, // 1 - dark blue
		{0x7E, 0x25, 0x53}, // 2 - dark purple
		{0x00, 0x87, 0x51}, // 3 - dark green
		{0xAB, 0x52, 0x36}, // 4 - brown
		{0x5F, 0x57, 0x4F}, // 5 - dark gray
		{0xC2, 0xC3, 0xC7}, // 6 - light gray
		{0xff, 0xf1, 0xe8}, // 7 - white
		{0xFF, 0x00, 0x4D}, // 8 - red
		{0xFF, 0xA3, 0x00}, // 9 - orange
		{0xFF, 0xEC, 0x27}, // 10 - yellow
		{0x00, 0xE4, 0x36}, // 11 - green
		{0x29, 0xAD, 0xFF}, // 12 - blue
		{0x83, 0x76, 0x9C}, // 13 - indigo
		{0xFF, 0x77, 0xA8}, // 14 - pink
		{0xFF, 0xCC, 0xAA}, // 15 - peach
	}
)

// Sset sets the pixel color on the sprite sheet. It does not update the global Color.
func Sset(x, y int, color byte) {
	if x < 0 {
		return
	}
	if y < 0 {
		return
	}
	if x >= ssWidth {
		return
	}
	if y >= ssHeight {
		return
	}

	SpriteSheetData[y*ssWidth+x] = color
}

// Sget gets the pixel color on the sprite sheet.
func Sget(x, y int) byte {
	if x < 0 {
		return 0
	}
	if y < 0 {
		return 0
	}
	if x >= ssWidth {
		return 0
	}
	if y >= ssHeight {
		return 0
	}

	return SpriteSheetData[y*ssWidth+x]
}

func loadSpriteSheet(resources fs.ReadFileFS) error {
	fileContents, err := resources.ReadFile("sprite-sheet.png")
	if errors.Is(err, fs.ErrNotExist) {
		useDefaultSpriteSheet()
		return nil
	}

	if err != nil {
		return fmt.Errorf("error loading sprite-sheet.png file: %w", err)
	}

	if err = useSpriteSheet(fileContents); err != nil {
		return err
	}

	return nil
}

func useDefaultSpriteSheet() {
	fmt.Printf("sprite-sheet.png file not found. Using empty sprite sheet %dx%d\n",
		SpriteSheetWidth, SpriteSheetHeight)

	size := SpriteSheetWidth * SpriteSheetHeight
	SpriteSheetData = make([]byte, size)
}

func useSpriteSheet(b []byte) error {
	img, err := image.DecodePNG(bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("error decoding file: %w", err)
	}

	SpriteSheetData = img.Pixels
	SpriteSheetWidth = img.Width
	SpriteSheetHeight = img.Height
	Palette = img.Palette
	return nil
}
