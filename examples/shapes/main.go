// Example showing how to draw shapes.
package main

import (
	"github.com/elgopher/pi"
)

func main() {
	pi.Draw = func() {
		pi.Cls()
		pi.Color(7)

		pi.Camera(-5, -5) // move every shape 5 pixels to the right, 5 pixels to the bottom

		// draw a filled square with side length=50
		pi.RectFill(0, 0, 49, 49) // x0=0, y0=0, x1=49, y1=49 (coords are inclusive)

		pi.Color(8)
		// draw a filled rectangle 10x20
		pi.RectFill(19+10, 15+20, 19, 15)

		pi.Color(3)
		// draw rect without filling. Will be drawn on top of existing pixels
		pi.Rect(10, 10, 80, 80)

		// draw line from x0,y0 to x1,y1 inclusive
		pi.Line(10, 10, 80, 80)

		pi.Line(80, 10, 10, 80)
	}

	pi.RunOrPanic()
}
