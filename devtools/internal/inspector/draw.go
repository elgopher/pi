// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package inspector

import (
	"fmt"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/devtools/internal/icons"
	"github.com/elgopher/pi/devtools/internal/snapshot"
	"github.com/elgopher/pi/vm"
)

var BgColor, FgColor byte

func Draw() {
	snapshot.Draw()

	// I check the input in Draw function because only during Draw operation
	// I have access to screen restored from snapshot
	if pi.MouseBtnp(pi.MouseLeft) {
		x, y := pi.MousePos()
		fmt.Printf("Screen pixel (%d, %d) with color %d selected\n", x, y, pi.Pget(x, y))
	}

	moveBarIfNeeded()
	drawBar()
	drawPointer()
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
