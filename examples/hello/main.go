// Example animating HELLO WORLD text on screen.
package main

import (
	"embed"

	"github.com/elgopher/pi"            // core package used to draw, print and control the input
	"github.com/elgopher/pi/devtools"   // tools used only during development
	"github.com/elgopher/pi/ebitengine" // engine capable of rendering the game on multiple operating systems
)

// Tell Go compiler to embed sprite-sheet file inside binary:
//
//go:embed sprite-sheet.png
var resources embed.FS

func main() {
	// Tell Pi to load resources embedded inside binary:
	pi.Load(resources)
	// Pi runs the game in a loop. 30 times per second, Pi
	// asks the game to draw a frame. This callback function
	// can be set by following code:
	pi.Draw = func() {
		// Clear entire screen each frame with color 0 (black):
		pi.Cls()
		// Draw "HELLO WORLD". Each letter is a different sprite.
		// Sprite 0 is H, 1 is E, 2 is L etc.
		for i := 0; i < 12; i++ {
			// Calculate the position of letter on screen:
			x := 20 + i*8
			y := pi.Cos(pi.Time()+float64(i)/64) * 60
			// Draw sprite:
			pi.Spr(i, x, 60+int(y))
		}
	}

	// Run game with devtools (Hit F12 to show screen inspector)
	devtools.MustRun(ebitengine.Run)
}
