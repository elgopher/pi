// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

// Game-loop function callbacks.
var (
	// Init is a user-defined function called once when the game starts.
	//
	// Use this function to initialize the game state, load assets,
	// or prepare data needed before the main loop begins.
	Init = func() {}

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
