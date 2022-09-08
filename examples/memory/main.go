// Example showing how to directly modify screen memory. Useful for doing low-level stuff such as writing your own
// functions to manipulate pixels.
package main

import (
	"math/rand"

	"github.com/elgopher/pi"
)

func main() {
	pi.Draw = func() {
		for i := 0; i < len(pi.ScreenData); i++ {
			randomColor := byte(rand.Intn(16))
			pi.ScreenData[i] = randomColor // put a random color to each pixel
		}
	}

	pi.MustRun()
}
