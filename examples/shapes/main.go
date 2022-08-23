// Example showing how to draw shapes and use a mouse.
package main

import (
	"embed"
	"fmt"
	"math"

	"github.com/elgopher/pi"
)

//go:embed sprite-sheet.png
var resources embed.FS

var drawShape func(x0, y0, x1, y1 int)

const (
	shapeColor = 15
	textColor  = 2
)

var drawFunctions = []func(x0, y0, x1, y1 int){
	func(x0, y0, x1, y1 int) {
		pi.Rect(x0, y0, x1, y1, shapeColor)
		command := fmt.Sprintf("Rect(%d,%d,%d,%d,%d)", x0, y0, x1, y1, shapeColor)
		pi.Print(command, textColor)
	},
	func(x0, y0, x1, y1 int) {
		pi.RectFill(x0, y0, x1, y1, shapeColor)
		command := fmt.Sprintf("RectFill(%d,%d,%d,%d,%d)", x0, y0, x1, y1, shapeColor)
		pi.Print(command, textColor)
	},
	func(x0, y0, x1, y1 int) {
		pi.Line(x0, y0, x1, y1, shapeColor)
		command := fmt.Sprintf("Line(%d,%d,%d,%d,%d)", x0, y0, x1, y1, shapeColor)
		pi.Print(command, textColor)
	},
	func(x0, y0, x1, y1 int) {
		r := radius(x0, y0, x1, y1)
		pi.Circ(x0, y0, r, shapeColor)

		command := fmt.Sprintf("Circ(%d,%d,%d,%d)", x0, y0, r, shapeColor)
		pi.Print(command, textColor)
	},
	func(x0, y0, x1, y1 int) {
		r := radius(x0, y0, x1, y1)
		pi.CircFill(x0, y0, r, shapeColor)
		command := fmt.Sprintf("CircFill(%d,%d,%d,%d)", x0, y0, r, shapeColor)
		pi.Print(command, textColor)
	},
}

var current = 0

var x0, y0 int

func main() {
	pi.Resources = resources
	pi.Draw = func() {
		pi.Cls()

		// change the shape if right mouse button was just pressed
		if pi.MouseBtnp(pi.MouseRight) {
			current++
			if current == len(drawFunctions) {
				current = 0
			}
		}

		// set initial coordinates on start dragging
		if pi.MouseBtn(pi.MouseLeft) && drawShape == nil {
			x0, y0 = pi.MousePos()
			drawShape = drawFunctions[current]

		}

		pi.Cursor(8, 128-6-6)

		// set coordinates during dragging
		if drawShape != nil {
			x1, y1 := pi.MousePos()
			drawShape(x0, y0, x1, y1)
		}

		if !pi.MouseBtn(pi.MouseLeft) {
			drawShape = nil
		}

		drawMousePointer()
	}

	pi.RunOrPanic()
}

func drawMousePointer() {
	x, y := pi.MousePos()
	pi.Spr(current, x, y)
}

func radius(x0, y0, x1, y1 int) int {
	dx := math.Abs(float64(x0 - x1))
	dy := math.Abs(float64(y0 - y1))
	return int(math.Max(dx, dy))
}
