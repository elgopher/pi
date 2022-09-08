// Example showing how to change screen resolution and run π functions before game loop.
package main

import (
	"github.com/elgopher/pi"
)

func main() {
	// set screen resolution:
	pi.ScreenWidth = 44
	pi.ScreenHeight = 44

	// boot the game with custom screen resolution:
	pi.MustBoot()

	// once boot is executed all drawing functions are available:
	pi.Cursor(0, 18)
	pi.Print("TINY SCREEN", 7) // print text on the screen before game loop

	// Run the game loop.
	pi.MustRun()
	// Update and Draw functions were not set therefore screen will be fixed.
}
