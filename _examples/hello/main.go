// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package main

import (
	"github.com/elgopher/pi"          // import pi core package
	"github.com/elgopher/pi/picofont" // import very small pico-8 font
	"github.com/elgopher/pi/piebiten" // import backend
)

func main() {
	pi.SetScreenSize(47, 9) // set custom screen size
	pi.Draw = func() {      // draw will be executed each frame
		picofont.Print("HELLO WORLD", 2, 2)
	}
	piebiten.Run() // run backend
}
