// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

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
