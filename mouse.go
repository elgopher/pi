// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"github.com/elgopher/pi/internal/input"
	"github.com/elgopher/pi/vm"
)

type MouseButton int

const (
	MouseLeft   MouseButton = vm.MouseLeft
	MouseMiddle MouseButton = vm.MouseMiddle
	MouseRight  MouseButton = vm.MouseRight
)

// MouseBtn returns true if the mouse button is being pressed at this moment.
func MouseBtn(b MouseButton) bool {
	if b < MouseLeft || b > MouseRight {
		return false
	}

	return vm.MouseBtnDuration[b] > 0
}

// MouseBtnp returns true when the mouse button has just been pressed.
// It also returns true after the next 15 frames, and then every 4 frames.
// This simulates keyboard-like repeating.
func MouseBtnp(b MouseButton) bool {
	if b < MouseLeft || b > MouseRight {
		return false
	}

	return input.IsPressedRepeatably(vm.MouseBtnDuration[b])
}

// MousePos returns the position of mouse in screen coordinates.
func MousePos() (x, y int) {
	return vm.MousePos.X, vm.MousePos.Y
}
