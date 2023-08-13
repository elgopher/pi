// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package inspector

import (
	"fmt"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/devtools/internal/snapshot"
	"github.com/elgopher/pi/key"
	"github.com/elgopher/pi/snap"
)

var BgColor, FgColor byte

var pixelColorAtMouseCoords byte

func Draw() {
	snapshot.Draw()
	pixelColorAtMouseCoords = pi.Get(pi.MousePos.X, pi.MousePos.Y)
	handleScreenshot()

	moveBarIfNeeded()

	if cursorOutOfWindow() {
		return
	}

	if toolbar.visible {
		toolbar.draw()
	} else {
		tool.Draw()
	}
}

func cursorOutOfWindow() bool {
	x, y := pi.MousePos.X, pi.MousePos.Y
	screen := pi.Scr()
	return x < 0 || x >= screen.Width() || y < 0 || y >= screen.Height()
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
