// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package vm

var (
	// ScreenData contains pixel colors for the screen visible by the player.
	// Each pixel is one byte. It is initialized during pi.Boot.
	//
	// Pixels on the screen are organized from left to right,
	// top to bottom. Slice element number 0 has pixel located
	// in the top-left corner. Slice element number 1 has pixel color
	// on the right and so on.
	//
	// Can be freely read and updated. Useful when you want to use your own
	// functions for pixel manipulation.
	// Pi will panic if you try to change the length of the slice.
	ScreenData []byte

	// ScreenWidth is the width of the screen (in pixels).
	ScreenWidth int

	// ScreenHeight is the height of the screen (in pixels).
	ScreenHeight int
)
