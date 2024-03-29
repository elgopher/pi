// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package inspector

import (
	"fmt"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/devtools/internal/icons"
	"github.com/elgopher/pi/devtools/internal/rgb"
)

type Measure struct {
}

func (m *Measure) Icon() byte {
	return icons.MeasureTool
}

func (m *Measure) Draw() {
	m.drawBar()
	m.drawDistanceLine()
	m.drawPointer()
}

func (m *Measure) drawBar() {
	screen := pi.Scr()
	var barY int
	if !isBarOnTop {
		barY = screen.Height() - 7
	}

	pi.RectFill(0, barY, screen.Width(), barY+6, BgColor)

	textX := 1
	textY := barY + 1

	if distance.measuring {
		m.printDistance(textX, textY)
	} else {
		mostX := m.printCoords(pi.MousePos.X, pi.MousePos.Y, textX, textY)
		m.printPixelColor(pixelColorAtMouseCoords, mostX+4, textY)
	}
}

func (m *Measure) printCoords(mouseX int, mouseY int, x, y int) int {
	coords := fmt.Sprintf("%d %d", mouseX, mouseY)
	return pi.Print(coords, x, y, FgColor)
}

func (m *Measure) printPixelColor(color byte, x int, y int) int {
	c := fmt.Sprintf("%d", color)
	return pi.Print(c, x, y, color)
}

func (m *Measure) choosePointerColor(x, y int) byte {
	c := pixelColorAtMouseCoords
	if rgb.BrightnessDelta(pi.Palette[FgColor], pi.Palette[c]) >= rgb.BrightnessDelta(pi.Palette[BgColor], pi.Palette[c]) {
		return FgColor
	}

	return BgColor
}

func (m *Measure) drawDistanceLine() {
	if distance.measuring {
		pi.Line(distance.startX, distance.startY, pi.MousePos.X, pi.MousePos.Y, BgColor)
	}
}

func (m *Measure) printDistance(x, y int) int {
	if distance.measuring {
		dist, width, height := calcDistance()
		text := fmt.Sprintf("D: %.1f W: %d H: %d", dist, width, height)
		return pi.Print(text, x, y, FgColor)
	}

	return x
}

func (m *Measure) drawPointer() {
	x, y := pi.MousePos.X, pi.MousePos.Y
	color := m.choosePointerColor(x, y)
	icons.Draw(x, y, color, icons.Pointer)
	icons.Draw(x+2, y+2, color, tool.Icon())
}

func (m *Measure) Update() {
	x, y := pi.MousePos.X, pi.MousePos.Y
	switch {
	case pi.MouseBtnp(pi.MouseLeft) && !distance.measuring:
		distance.measuring = true
		distance.startX, distance.startY = x, y
		fmt.Printf("\nMeasuring started at (%d, %d)\n", x, y)
	case !pi.MouseBtn(pi.MouseLeft) && distance.measuring:
		distance.measuring = false
		dist, width, height := calcDistance()
		fmt.Printf("Measuring stopped at (%d, %d). Distance is: %f, width: %d, height: %d.\n", x, y, dist, width, height)
	}
}
