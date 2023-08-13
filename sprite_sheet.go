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
	SpriteHeight = 8
	SpriteWidth  = 8
)

var sprSheet = newSpriteSheet(defaultSpriteSheetWidth, defaultSpriteSheetHeight)

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
	copy(sprSheet.Pix(), img.Pix)
	return nil
}

type spriteSheet struct {
	PixMap

	numberOfSprites int
	spritesInLine   int
}

// SprSheet returns the sprite-sheet PixMap
func SprSheet() PixMap {
	return sprSheet.PixMap
}

func newSpriteSheet(w int, h int) spriteSheet {
	if w%8 != 0 || w == 0 {
		panic(fmt.Sprintf("sprite sheet width %d is not a multiplcation of 8", w))
	}
	if h%8 != 0 || h == 0 {
		panic(fmt.Sprintf("sprite sheet height %d is not a multiplcation of 8", h))
	}

	size := w * h

	return spriteSheet{
		PixMap:          NewPixMap(w, h),
		numberOfSprites: size / (SpriteWidth * SpriteHeight),
		spritesInLine:   w / SpriteWidth,
	}
}

// UseEmptySpriteSheet initializes empty sprite-sheet with given size. Could be used
// when you don't have sprite-sheet.png in resources.
func UseEmptySpriteSheet(w, h int) {
	sprSheet = newSpriteSheet(w, h)
}
