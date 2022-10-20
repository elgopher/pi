// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package vm

var (
	// SpriteSheetWidth is a sprite-sheet width in pixels
	SpriteSheetWidth int
	// SpriteSheetHeight is a sprite-sheet height in pixels
	SpriteSheetHeight int

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
)
