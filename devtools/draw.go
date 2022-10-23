// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package devtools

import (
	"fmt"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/devtools/internal/icons"
	"github.com/elgopher/pi/devtools/internal/snapshot"
	"github.com/elgopher/pi/vm"
)

func drawDevTools() {
	snapshot.Draw()
	moveBarIfNeeded()
	drawBar()
	drawPointer()
}

var isBarOnTop bool

func moveBarIfNeeded() {
	_, mouseY := pi.MousePos()
	switch {
	case isBarOnTop && mouseY <= 12:
		isBarOnTop = false
	case !isBarOnTop && mouseY >= vm.ScreenHeight-12:
		isBarOnTop = true
	}
}

func drawBar() {
	mouseX, mouseY := pi.MousePos()
	var barY int
	if !isBarOnTop {
		barY = vm.ScreenHeight - 7
	}

	pi.RectFill(0, barY, vm.ScreenWidth, barY+6, BgColor)

	mostX := printCoords(mouseX, mouseY, 1, barY+1)
	color := pi.Pget(mouseX, mouseY)
	printPixelColor(color, mostX+4, barY+1)
}

func printCoords(mouseX int, mouseY int, x, y int) int {
	coords := fmt.Sprintf("%d %d", mouseX, mouseY)
	return pi.Print(coords, x, y, FgColor)
}

func printPixelColor(color byte, x int, y int) int {
	c := fmt.Sprintf("%d", color)
	return pi.Print(c, x, y, color)
}

func drawPointer() {
	x, y := pi.MousePos()
	icons.Draw(icons.Pointer, x, y, FgColor)
}
