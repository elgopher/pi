// Example showing how to change screen resolution and run π functions before game loop.
package main

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/ebitengine"
)

func main() {
	// set screen resolution:
	pi.SetScreenSize(44, 44)

	// all drawing functions are available before running the game:
	pi.Print("TINY SCREEN", 0, 18, 7) // print text on the screen before game loop

	// Run the game loop.
	ebitengine.MustRun()
	// Update and Draw functions were not set therefore screen will not be changed by them
}
