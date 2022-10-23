// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package inspector

import (
	"fmt"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/devtools/internal/icons"
	"github.com/elgopher/pi/devtools/internal/rgb"
	"github.com/elgopher/pi/devtools/internal/snapshot"
	"github.com/elgopher/pi/key"
	"github.com/elgopher/pi/snap"
	"github.com/elgopher/pi/vm"
)

var BgColor, FgColor byte

var pixelColorAtMouseCoords byte

func Draw() {
	snapshot.Draw()
	pixelColorAtMouseCoords = pi.Pget(pi.MousePos())
	handleScreenshot()

	moveBarIfNeeded()
	drawBar()

	drawDistanceLine()

	drawPointer()
}

func drawBar() {
	mouseX, mouseY := pi.MousePos()
	var barY int
	if !isBarOnTop {
		barY = vm.ScreenHeight - 7
	}

	pi.RectFill(0, barY, vm.ScreenWidth, barY+6, BgColor)

	textX := 1
	textY := barY + 1

	if distance.measuring {
		printDistance(textX, textY)
	} else {
		mostX := printCoords(mouseX, mouseY, textX, textY)
		printPixelColor(pixelColorAtMouseCoords, mostX+4, textY)
	}
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
	icons.Draw(icons.Pointer, x, y, choosePointerColor(x, y))
}

func choosePointerColor(x, y int) byte {
	c := pixelColorAtMouseCoords
	if rgb.BrightnessDelta(vm.Palette[FgColor], vm.Palette[c]) >= rgb.BrightnessDelta(vm.Palette[BgColor], vm.Palette[c]) {
		return FgColor
	}

	return BgColor
}

func drawDistanceLine() {
	if distance.measuring {
		x, y := pi.MousePos()
		pi.Line(distance.startX, distance.startY, x, y, BgColor)
	}
}

func printDistance(x, y int) int {
	if distance.measuring {
		dist, width, height := calcDistance()
		text := fmt.Sprintf("D: %.1f W: %d H: %d", dist, width, height)
		return pi.Print(text, x, y, FgColor)
	}

	return x
}

func handleScreenshot() {
	if key.Btnp(key.P) {
		path, err := snap.Take()
		if err != nil {
			fmt.Println("Problem taking screenshot:", err)
			return
		}
		fmt.Println("Screenshot taken and stored to: " + path)
	}
}
