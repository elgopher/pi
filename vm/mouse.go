// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package vm

// MouseBtnDuration has how many frames in a row a mouse button was pressed:
// Index of array is equal to mouse button constant.
var MouseBtnDuration [3]uint

const (
	MouseLeft   = 0
	MouseMiddle = 1
	MouseRight  = 2
)

// MousePos is the position of mouse in screen coordinates.
var MousePos Pos
