// Example showing how to directly modify screen memory. Useful for doing low-level stuff such as writing your own
// functions to manipulate pixels.
package main

import (
	"math/rand"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/ebitengine"
)

func main() {
	pi.Draw = func() {
		pixels := pi.Scr().Pix()
		for i := 0; i < len(pixels); i++ {
			randomColor := byte(rand.Intn(16))
			pixels[i] = randomColor // put a random color to each pixel
		}
	}

	ebitengine.MustRun()
}
