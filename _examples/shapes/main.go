// Example showing how to draw shapes and use a mouse.
package main

import (
	_ "embed"
	"fmt"
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/picofont"
	"github.com/elgopher/pi/piebiten"
	"github.com/elgopher/pi/pimouse"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

const (
	shapeColor = 15
	textColor  = 2
)

var (
	//go:embed sprite-sheet.png
	spriteSheetPNG []byte

	drawShape func(start, stop pi.Position)

	drawFunctions = []func(start, stop pi.Position){
		func(start, stop pi.Position) {
			pi.Rect(start.X, start.Y, stop.X, stop.Y)
			command := fmt.Sprintf("pi.Rect(%d,%d,%d,%d)", start.X, start.Y, stop.X, stop.Y)
			printCmd(command)
		},
		func(start, stop pi.Position) {
			pi.RectFill(start.X, start.Y, stop.X, stop.Y)
			command := fmt.Sprintf("pi.RectFill(%d,%d,%d,%d)", start.X, start.Y, stop.X, stop.Y)
			printCmd(command)
		},
		func(start, stop pi.Position) {
			pi.Line(start.X, start.Y, stop.X, stop.Y)
			command := fmt.Sprintf("pi.Line(%d,%d,%d,%d)", start.X, start.Y, stop.X, stop.Y)
			printCmd(command)
		},
		func(start, stop pi.Position) {
			r := radius(start.X, start.Y, stop.X, stop.Y)
			pi.Circ(start.X, start.Y, r)

			command := fmt.Sprintf("pi.Circ(%d,%d,%d)", start.X, start.Y, r)
			printCmd(command)
		},
		func(start, stop pi.Position) {
			r := radius(start.X, start.Y, stop.X, stop.Y)
			pi.CircFill(start.X, start.Y, r)
			command := fmt.Sprintf("pi.CircFill(%d,%d,%d)", start.X, start.Y, r)
			printCmd(command)
		},
	}

	currentShapeIdx = 0

	shapeStart pi.Position

	cursorSprites []pi.Sprite
)

func main() {
	pi.SetScreenSize(128, 128)
	pi.TPS = 60 // pi.Update and pi.Draw wil be executed 60 times per second

	pi.Palette = pi.DecodePalette(spriteSheetPNG)
	spriteSheet := pi.DecodeCanvas(spriteSheetPNG)

	// create cursors sprite array for each shape
	for i := range drawFunctions {
		cursorSprite := pi.SpriteFrom(spriteSheet, i*8, 0, 8, 8)
		cursorSprites = append(cursorSprites, cursorSprite)
	}

	pi.Draw = func() {
		pi.Cls()

		// change the shape if right mouse button was just pressed
		if pimouse.Duration(pimouse.Right) == 1 {
			currentShapeIdx++
			if currentShapeIdx == len(drawFunctions) {
				currentShapeIdx = 0
			}
		}

		// set initial coordinates on start dragging
		if pimouse.Duration(pimouse.Left) > 0 && drawShape == nil {
			shapeStart = pimouse.Position
			drawShape = drawFunctions[currentShapeIdx]

		}

		// set coordinates during dragging
		if drawShape != nil {
			stop := pimouse.Position
			pi.SetColor(shapeColor)
			drawShape(shapeStart, stop)
		}

		if pimouse.Duration(pimouse.Left) == 0 {
			drawShape = nil
		}

		drawMousePointer()
	}

	ebiten.SetCursorMode(ebiten.CursorModeHidden) // hide cursor in Ebitengine
	piebiten.Run()
}

func drawMousePointer() {
	pi.Spr(cursorSprites[currentShapeIdx], pimouse.Position.X, pimouse.Position.Y)
}

func radius(x0, y0, x1, y1 int) int {
	dx := math.Abs(float64(x0 - x1))
	dy := math.Abs(float64(y0 - y1))
	return int(math.Max(dx, dy))
}

func printCmd(command string) {
	pi.SetColor(textColor)
	picofont.Print(command, 8, 128-6-6)
}
