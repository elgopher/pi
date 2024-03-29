// Example plotting sin and cos on screen
package main

import (
	"math"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/ebitengine"
)

var start float64

func main() {
	pi.Update = func() {
		start += 1
	}

	pi.Draw = func() {
		pi.Cls()
		draw(32, 8, pi.Sin)
		draw(96, 11, pi.Cos)
	}

	ebitengine.MustRun()
}

func draw(line int, color byte, f func(x float64) float64) {
	drawHorizontalAxis(line)

	for x := 0.0; x < 128; x++ {
		angle := (x + start) / 128
		dy := math.Round(f(angle) * 16)
		pi.Set(int(x), line+int(dy), color)
	}
}

func drawHorizontalAxis(line int) {
	for x := 0; x < 128; x++ {
		pi.Set(x, line, 1)
	}
}
