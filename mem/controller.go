// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package mem

var Controllers [8]Controller // 0th element is for Player 0, 1st for Player 1 etc.

type Controller struct {
	// BtnDuration is how many frames button was pressed:
	// Index of array is equal to controller button constant.
	BtnDuration [6]uint
}

const (
	ControllerLeft  = 0
	ControllerRight = 1
	ControllerUp    = 2
	ControllerDown  = 3
	ControllerO     = 4 // O is a first fire button
	ControllerX     = 5 // X is a second fire button
)
