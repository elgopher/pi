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

var sprSheet = newSpriteSheet(defaultSpriteSheetWidth, defaultSpriteSheetHeight)

// Sset sets the pixel color on the sprite sheet. It does not update the global Color.
func Sset(x, y int, color byte) {
	if x < 0 {
		return
	}
	if y < 0 {
		return
	}
	if x >= sprSheet.W {
		return
	}
	if y >= sprSheet.H {
		return
	}

	sprSheet.Pix[y*sprSheet.W+x] = color
}

// Sget gets the pixel color on the sprite sheet.
func Sget(x, y int) byte {
	if x < 0 {
		return 0
	}
	if y < 0 {
		return 0
	}
	if x >= sprSheet.W {
		return 0
	}
	if y >= sprSheet.H {
		return 0
	}

	return sprSheet.Pix[y*sprSheet.W+x]
}

func loadSpriteSheet(resources fs.ReadFileFS) error {
	fileContents, err := resources.ReadFile("sprite-sheet.png")
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error loading sprite-sheet.png file: %w", err)
	}

	return useSpriteSheet(fileContents)
}

func useSpriteSheet(b []byte) error {
	img, err := image.DecodePNG(bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("error decoding file: %w", err)
	}

	Palette = img.Palette
	sprSheet = newSpriteSheet(img.Width, img.Height)
	sprSheet.Pix = img.Pixels
	return nil
}

type SpriteSheet struct {
	// Width and height in pixels
	W, H int

	// Pix contains pixel colors for the entire sprite sheet.
	// Each pixel is one byte. It is initialized during pi.Boot.
	//
	// Pixels in the sprite-sheet are organized from left to right,
	// top to bottom. Slice element number 0 has pixel located
	// in the top-left corner. Slice element number 1 has a pixel color
	// on the right and so on.
	//
	// Can be freely read and updated.
	// Useful when you want to use your own functions for pixel manipulation.
	Pix []byte

	numberOfSprites int
	spritesInLine   int
}

func SprSheet() SpriteSheet {
	return sprSheet
}

func newSpriteSheet(w int, h int) SpriteSheet {
	if w%8 != 0 || w == 0 {
		panic(fmt.Sprintf("sprite sheet width %d is not a multiplcation of 8", w))
	}
	if h%8 != 0 || h == 0 {
		panic(fmt.Sprintf("sprite sheet height %d is not a multiplcation of 8", h))
	}

	size := w * h

	return SpriteSheet{
		W:               w,
		H:               h,
		Pix:             make([]byte, size),
		numberOfSprites: size / (SpriteWidth * SpriteHeight),
		spritesInLine:   w / SpriteWidth,
	}
}

// UseEmptySpriteSheet initializes empty sprite-sheet with given size. Could be used
// when you don't have sprite-sheet.png in resources.
func UseEmptySpriteSheet(w, h int) {
	sprSheet = newSpriteSheet(w, h)
}
