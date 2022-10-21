// Example animating HELLO WORLD text on screen.
package main

import (
	"embed"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/ebitengine"
)

//go:embed sprite-sheet.png
var resources embed.FS

func main() {
	pi.Resources = resources
	pi.Draw = func() {
		pi.Cls()
		// Draw "HELLO WORLD". Each letter is a different sprite.
		// Sprite 0 is H, 1 is E, 2 is L etc.
		for i := 0; i < 12; i++ {
			// calculate the position of letter on screen:
			x := 20 + i*8
			y := pi.Cos(pi.Time()+float64(i)/64) * 60
			// draw sprite:
			pi.Spr(i, x, 60+int(y))
		}
	}
	pi.MustRun(ebitengine.Backend)
}
