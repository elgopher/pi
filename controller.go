// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"github.com/elgopher/pi/internal/input"
)

// Button is a virtual button on any game controller. The game controller can be a gamepad or a keyboard.
//
// Button is used by Btn, Btnp, BtnPlayer, BtnpPlayer, BtnBits and BtnpBits.
type Button int

// Keyboard mappings:
//
//	player 0: [DPAD] - cursors  [O] - Z C N   [X] - X V M
//	player 1: [DPAD] - SFED     [O] - LSHIFT  [X] - TAB W Q A
//
// First connected gamepad controller is player 0, second player 1 and so on.
// On XBox controller [O] is A and Y, [X] is B and X.
const (
	Left  Button = 0
	Right Button = 1
	Up    Button = 2
	Down  Button = 3
	O     Button = 4 // O is a first fire button
	X     Button = 5 // X is a second fire button
)

var Controllers [8]Controller // 0th element is for Player 0, 1st for Player 1 etc.

type Controller struct {
	// BtnDuration is how many frames button was pressed:
	// Index of array is equal to controller button constant.
	BtnDuration [6]uint
}

// Btn returns true if a controller button is being pressed at this moment by player 0.
func Btn(button Button) bool {
	return BtnPlayer(button, 0)
}

// BtnPlayer returns true if a controller button is being pressed at this moment
// by specific player. The player can be 0..7.
func BtnPlayer(button Button, player int) bool {
	return isPressed(buttonDuration(player, button))
}

// Btnp returns true when the controller button has just been pressed.
// It also returns true after the next 15 frames, and then every 4 frames.
// This simulates keyboard-like repeating.
func Btnp(button Button) bool {
	return BtnpPlayer(button, 0)
}

// BtnpPlayer returns true when the controller button has just been pressed.
// It also returns true after the next 15 frames, and then every 4 frames.
// This simulates keyboard-like repeating. The player can be 0..7.
func BtnpPlayer(button Button, player int) bool {
	return input.IsPressedRepeatably(buttonDuration(player, button))
}

// BtnBits returns the state of all controller buttons for players 0 and 1 as bitset.
//
// The first byte contains the button states for player 0 (bits 0 through 5, bits 6 and 7 are unused).
// The second byte contains the button states for player 1 (bits 8 through 13).
//
// Bit 0 is Left, 1 is Right, bit 5 is the X button.
//
// A bit of 1 means the button is pressed.
func BtnBits() int {
	return buttonBits(0, isPressed) + buttonBits(1, isPressed)<<8
}

// BtnpBits returns the state of all controller buttons for players 0 and 1 as bitset.
//
// The first byte contains the button states for player 0 (bits 0 through 5, bits 6 and 7 are unused).
// The second byte contains the button states for player 1 (bits 8 through 13).
//
// Bit 0 is Left, 1 is Right, bit 5 is the X button.
//
// A bit of 1 means the button has just been pressed.
func BtnpBits() int {
	return buttonBits(0, input.IsPressedRepeatably) + buttonBits(1, input.IsPressedRepeatably)<<8
}

func buttonDuration(player int, button Button) uint {
	if button < Left || button > X {
		return 0
	}

	if player < 0 || player > 7 {
		return 0
	}

	return Controllers[player].BtnDuration[button]
}

func isPressed(duration uint) bool {
	return duration > 0
}

func buttonBits(player int, isSet func(uint) bool) int {
	c := Controllers[player]
	var b int
	for i := 0; i <= int(X); i++ {
		if isSet(c.BtnDuration[i]) {
			b += 1 << i
		}
	}
	return b
}
