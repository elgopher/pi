// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import "github.com/elgopher/pi/internal/input"

type MouseButton int

const (
	MouseLeft   MouseButton = 0
	MouseMiddle MouseButton = 1
	MouseRight  MouseButton = 2
)

var (
	// MouseBtnDuration has how many frames in a row a mouse button was pressed:
	// Index of array is equal to mouse button constant.
	MouseBtnDuration [3]uint

	// MousePos is the position of mouse in screen coordinates.
	MousePos Position
)

// MouseBtn returns true if the mouse button is being pressed at this moment.
func MouseBtn(b MouseButton) bool {
	if b < MouseLeft || b > MouseRight {
		return false
	}

	return MouseBtnDuration[b] > 0
}

// MouseBtnp returns true when the mouse button has just been pressed.
// It also returns true after the next 15 frames, and then every 4 frames.
// This simulates keyboard-like repeating.
func MouseBtnp(b MouseButton) bool {
	if b < MouseLeft || b > MouseRight {
		return false
	}

	return input.IsPressedRepeatably(MouseBtnDuration[b])
}
