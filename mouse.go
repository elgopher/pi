// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type MouseButton int

const (
	MouseLeft   MouseButton = 0
	MouseMiddle MouseButton = 1
	MouseRight  MouseButton = 2
)

var mapping = []ebiten.MouseButton{
	ebiten.MouseButtonLeft,
	ebiten.MouseButtonMiddle,
	ebiten.MouseButtonRight,
}

// MouseBtn returns true if the mouse button is being pressed at this moment.
func MouseBtn(b MouseButton) bool {
	if b < MouseLeft || b > MouseRight {
		return false
	}

	return ebiten.IsMouseButtonPressed(mapping[b])
}

// MouseBtnp returns true when the mouse button has just been pressed.
// It also returns true after the next 15 frames, and then every 4 frames.
// This simulates keyboard-like repeating.
func MouseBtnp(b MouseButton) bool {
	if b < MouseLeft || b > MouseRight {
		return false
	}

	return isPressedRepeatably(mouseButtonDuration[b])
}

// MousePos returns the position of mouse in screen coordinates.
func MousePos() (x, y int) {
	x, y = ebiten.CursorPosition()
	return
}

var mouseButtonDuration [3]int // left, middle, right

func updateMouse() {
	for i := 0; i < len(mapping); i++ {
		button := mapping[i]
		if ebiten.IsMouseButtonPressed(button) {
			mouseButtonDuration[i] += 1
		} else {
			mouseButtonDuration[i] = 0
		}
	}
}
