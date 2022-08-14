// Example plotting sin and cos on screen
package main

import (
	"math"

	"github.com/elgopher/pi"
)

var start float64

func main() {
	pi.Update = func() {
		start += 1
	}

	pi.Draw = func() {
		pi.Cls()
		pi.Color(1)

		draw(32, 8, pi.Sin)
		draw(96, 11, pi.Cos)

	}
	pi.RunOrPanic()
}

func draw(line int, color byte, f func(x float64) float64) {
	pi.Color(1)
	drawHorizontalAxis(line)

	pi.Color(color)
	for x := 0.0; x < 128; x++ {
		angle := (x + start) / 128
		dy := math.Round(f(angle) * 16)
		pi.Pset(int(x), line+int(dy))
	}
}

func drawHorizontalAxis(line int) {
	for x := 0; x < 128; x++ {
		pi.Pset(x, line)
	}
}