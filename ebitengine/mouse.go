// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitengine

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elgopher/pi/vm"
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
			vm.MouseBtnDuration[i] += 1
		} else {
			vm.MouseBtnDuration[i] = 0
		}
	}

	x, y := ebiten.CursorPosition()
	vm.MousePos.X = x
	vm.MousePos.Y = y
}
