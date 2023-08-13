// Example showing how to draw shapes and use a mouse.
package main

import (
	"embed"
	"fmt"
	"math"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/ebitengine"
)

const (
	shapeColor = 15
	textColor  = 2
)

var (
	//go:embed sprite-sheet.png
	resources embed.FS

	drawShape func(start, stop pi.Position)

	drawFunctions = []func(start, stop pi.Position){
		func(start, stop pi.Position) {
			pi.Rect(start.X, start.Y, stop.X, stop.Y, shapeColor)
			command := fmt.Sprintf("Rect(%d,%d,%d,%d,%d)", start.X, start.Y, stop.X, stop.Y, shapeColor)
			printCmd(command)
		},
		func(start, stop pi.Position) {
			pi.RectFill(start.X, start.Y, stop.X, stop.Y, shapeColor)
			command := fmt.Sprintf("RectFill(%d,%d,%d,%d,%d)", start.X, start.Y, stop.X, stop.Y, shapeColor)
			printCmd(command)
		},
		func(start, stop pi.Position) {
			pi.Line(start.X, start.Y, stop.X, stop.Y, shapeColor)
			command := fmt.Sprintf("Line(%d,%d,%d,%d,%d)", start.X, start.Y, stop.X, stop.Y, shapeColor)
			printCmd(command)
		},
		func(start, stop pi.Position) {
			r := radius(start.X, start.Y, stop.X, stop.Y)
			pi.Circ(start.X, start.Y, r, shapeColor)

			command := fmt.Sprintf("Circ(%d,%d,%d,%d)", start.X, start.Y, r, shapeColor)
			printCmd(command)
		},
		func(start, stop pi.Position) {
			r := radius(start.X, start.Y, stop.X, stop.Y)
			pi.CircFill(start.X, start.Y, r, shapeColor)
			command := fmt.Sprintf("CircFill(%d,%d,%d,%d)", start.X, start.Y, r, shapeColor)
			printCmd(command)
		},
	}

	current = 0

	start pi.Position
)

func main() {
	pi.Load(resources)
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
			start = pi.MousePos
			drawShape = drawFunctions[current]

		}

		// set coordinates during dragging
		if drawShape != nil {
			stop := pi.MousePos
			drawShape(start, stop)
		}

		if !pi.MouseBtn(pi.MouseLeft) {
			drawShape = nil
		}

		drawMousePointer()
	}

	ebitengine.MustRun()
}

func drawMousePointer() {
	pi.Spr(current, pi.MousePos.X, pi.MousePos.Y)
}

func radius(x0, y0, x1, y1 int) int {
	dx := math.Abs(float64(x0 - x1))
	dy := math.Abs(float64(y0 - y1))
	return int(math.Max(dx, dy))
}

func printCmd(command string) {
	pi.Print(command, 8, 128-6-6, textColor)
}
