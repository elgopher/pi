// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

// Game-loop function callbacks.
var (
	// Update is a user-provided function called every frame.
	//
	// This function handles user input, performs calculations, and updates
	// the game state. Typically, it does not draw anything on the screen.
	Update = func() {}

	// Draw is a user-provided function called up to once per frame.
	// Pi may skip calling this function if the previous frame took too long.
	//
	// The purpose of this function is to draw to the screen.
	Draw = func() {}
)
