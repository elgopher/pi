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

func loadResources(resources fs.ReadFileFS) error {
	return loadSpriteSheet(resources)
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

//nolint:govet
var defaultPalette = [256]image.RGB{
	{0, 0, 0},
	{0x1D, 0x2B, 0x53},
	{0x7E, 0x25, 0x53},
	{0x00, 0x87, 0x51},
	{0xAB, 0x52, 0x36},
	{0x5F, 0x57, 0x4F},
	{0xC2, 0xC3, 0xC7},
	{0xff, 0xf1, 0xe8},
	{0xFF, 0x00, 0x4D},
	{0xFF, 0xA3, 0x00},
	{0xFF, 0xEC, 0x27},
	{0x00, 0xE4, 0x36},
	{0x29, 0xAD, 0xFF},
	{0x83, 0x76, 0x9C},
	{0xFF, 0x77, 0xA8},
	{0xFF, 0xCC, 0xAA},
}
