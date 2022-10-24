// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package mem

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
