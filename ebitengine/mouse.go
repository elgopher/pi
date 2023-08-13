// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitengine

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elgopher/pi"
)

var mouseMapping = []ebiten.MouseButton{
	ebiten.MouseButtonLeft,
	ebiten.MouseButtonMiddle,
	ebiten.MouseButtonRight,
}

func updateMouse() {
	for i := 0; i < len(mouseMapping); i++ {
		button := mouseMapping[i]
		if ebiten.IsMouseButtonPressed(button) {
			pi.MouseBtnDuration[i] += 1
		} else {
			pi.MouseBtnDuration[i] = 0
		}
	}

	x, y := ebiten.CursorPosition()
	pi.MousePos.X = x
	pi.MousePos.Y = y
}
