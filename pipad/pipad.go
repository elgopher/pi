// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package pipad provides functionality for a virtual gamepad
// inspired by controllers from the 1990s.
//
// It has no analog sticks or analog buttons. All buttons are digital (on/off),
// with no pressure sensitivity. The virtual gamepad has the following layout:
//
//	  ##                   Y
//	######              X     B
//	  ##                   A
//
// On the left side is the d-pad/stick for movement in all directions.
// On the right side are the "fire" buttons named A, B, X, and Y.
//
// This package lets you check which buttons were pressed on
// any modern controller (Xbox, Steam Deck, PlayStation, Switch).
// Because the Xbox controller is the most common on PC,
// this package uses Xbox button naming. The "fire" buttons
// are on the right side — A, B, X, Y. Directional input
// on the Xbox controller is handled by the left stick and d-pad.
// However, regardless of the user's controller, the same code
// will work across all these devices, even if their button
// labels differ.
//
// When displaying instructions to the user about which button to press,
// you have several options:
//   - Show the Xbox button name (or icon). This is the least inclusive
//     but works well for prototyping.
//   - Let the user choose their controller type in your game's settings.
//     Depending on their choice, you can show the appropriate
//     button name or icon. You'd need to build this mapping yourself,
//     typically covering the two most popular pads: Xbox (and similar)
//     and PlayStation.
//   - Instead of showing the button name, display an icon representing
//     its position on the controller. Most controllers have a cluster
//     of four buttons on the right side (PlayStation, Switch).
//     For example, the "A" button is at the bottom of this cluster
//     on an Xbox controller, but on PlayStation it's the "X" button,
//     and on Switch it's "B". All of these map to the same constant pipad.A,
//     so your code can handle all controllers without changes.
//   - Use action aliases that the user can configure in your game's settings,
//     such as "jump", "shoot", or "block". This is the most flexible approach,
//     though not always the most intuitive for the player.
//
// The pipad package supports up to 16 buttons but defines only
// 8 public constants: A, B, X, Y, Left, Right, Top, and Bottom.
// This is intentional, encouraging you to design your game
// to use no more than 8 buttons. Why? Because retro games
// didn't use that many buttons, and this ensures compatibility
// with nearly any modern controller. Of course, you can allow
// users to map additional buttons if they want — for example,
// by capturing the button name from an EventButton and saving
// it in a settings file. Just remember that some controllers
// have even more than 16 buttons, but Pi does not support them.
package pipad

import (
	"log"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/pievent"
)

type Button string

// directional buttons
const (
	Left   Button = "Left"
	Right  Button = "Right"
	Top    Button = "Top"
	Bottom Button = "Bottom"
)

// fire buttons
const (
	A Button = "A"
	B Button = "B"
	X Button = "X"
	Y Button = "Y"
)

// Duration returns button press duration for any controller
func Duration(b Button) int {
	if buttonAnyDownFrame[b] == 0 {
		return 0
	}
	return pi.Frame - buttonAnyDownFrame[b]
}

// PlayerCount returns the number of connected controllers
func PlayerCount() int {
	return playerCount
}

func PlayerDuration(b Button, player int) int {
	if buttonDownFrame[player] == nil {
		return 0
	}

	if buttonDownFrame[player][b] == 0 {
		return 0
	}
	return pi.Frame - buttonDownFrame[player][b]
}

var buttonDownFrame = map[int]map[Button]int{}
var buttonAnyDownFrame = map[Button]int{}

func init() {
	ButtonTarget().SubscribeAll(onButton)
	ConnectionTarget().SubscribeAll(onConnection)
}

func onButton(event EventButton, _ pievent.Handler) {
	if buttonDownFrame[event.Player] == nil {
		buttonDownFrame[event.Player] = map[Button]int{}
	}

	switch event.Type {
	case EventDown:
		buttonDownFrame[event.Player][event.Button] = pi.Frame
		buttonAnyDownFrame[event.Button] = pi.Frame
	case EventUp:
		buttonDownFrame[event.Player][event.Button] = 0
		buttonAnyDownFrame[event.Button] = 0
	}
}

var playerCount = 0

func onConnection(event EventConnection, _ pievent.Handler) {
	if event.Type == EventDisconnected {
		log.Println("Controller disconnected", event.Player)
		buttonDownFrame[event.Player] = map[Button]int{}
		playerCount -= 1
	} else {
		log.Println("Controller connected", event.Player)
		buttonDownFrame[event.Player] = map[Button]int{} // na wszelki
		playerCount += 1
	}
}
