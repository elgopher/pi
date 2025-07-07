// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import "fmt"

var screen Canvas

func Screen() Canvas {
	return screen
}

// SetScreenSize sets the screen size in pixels.
//
// The current screen will be replaced with a new, cleared screen.
// The new screen also automatically becomes the current draw target.
//
// The maximum number of pixels is 131072 (128 KB).
func SetScreenSize(width, height int) {
	if width == screen.width && height == screen.height {
		return
	}
	if width <= 0 {
		panic(fmt.Sprintf("screen width %d is not greather than 0", width))
	}
	if height <= 0 {
		panic(fmt.Sprintf("screen height %d is not greather than 0", height))
	}

	const maxScreenSize = 128 * 1024
	if width*height > maxScreenSize {
		panic(fmt.Sprintf("number of pixels for screen resolution %dx%d is "+
			"higher than maximum %d. Please use smaller screen.",
			width, height, maxScreenSize))
	}

	screen = NewCanvas(width, height)
	SetDrawTarget(screen)
}

// Cls clears the screen using color 0.
func Cls() {
	Screen().Clear(0)
}

func init() {
	SetScreenSize(320, 180)
}
