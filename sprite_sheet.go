// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"

	"github.com/elgopher/pi/image"
	"github.com/elgopher/pi/vm"
)

const (
	SpriteWidth, SpriteHeight = 8, 8
)

// Sprite-sheet data
var (
	numberOfSprites int
	spritesInLine   int
)

// Sset sets the pixel color on the sprite sheet. It does not update the global Color.
func Sset(x, y int, color byte) {
	if x < 0 {
		return
	}
	if y < 0 {
		return
	}
	if x >= vm.SpriteSheetWidth {
		return
	}
	if y >= vm.SpriteSheetHeight {
		return
	}

	vm.SpriteSheetData[y*vm.SpriteSheetWidth+x] = color
}

// Sget gets the pixel color on the sprite sheet.
func Sget(x, y int) byte {
	if x < 0 {
		return 0
	}
	if y < 0 {
		return 0
	}
	if x >= vm.SpriteSheetWidth {
		return 0
	}
	if y >= vm.SpriteSheetHeight {
		return 0
	}

	return vm.SpriteSheetData[y*vm.SpriteSheetWidth+x]
}

func loadSpriteSheet(resources fs.ReadFileFS) error {
	fileContents, err := resources.ReadFile("sprite-sheet.png")
	if errors.Is(err, fs.ErrNotExist) {
		useEmptySpriteSheet()
		return nil
	}

	if err != nil {
		return fmt.Errorf("error loading sprite-sheet.png file: %w", err)
	}

	return useSpriteSheet(fileContents)
}

func useEmptySpriteSheet() {
	vm.SpriteSheetWidth = SpriteSheetWidth
	vm.SpriteSheetHeight = SpriteSheetHeight
	vm.Palette = Palette

	fmt.Printf("sprite-sheet.png file not found. Using empty sprite sheet %dx%d\n",
		vm.SpriteSheetWidth, vm.SpriteSheetHeight)

	size := vm.SpriteSheetWidth * vm.SpriteSheetHeight
	vm.SpriteSheetData = make([]byte, size)
}

func useSpriteSheet(b []byte) error {
	img, err := image.DecodePNG(bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("error decoding file: %w", err)
	}

	vm.SpriteSheetData = img.Pixels
	vm.Palette = img.Palette
	SpriteSheetWidth = img.Width
	SpriteSheetHeight = img.Height
	vm.SpriteSheetWidth = img.Width
	vm.SpriteSheetHeight = img.Height
	return nil
}
