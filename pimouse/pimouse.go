// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pimouse

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/pievent"
)

var (
	Position      pi.Position
	MovementDelta pi.Position // mouse movement delta since the last frame
)

func Duration(b Button) int {
	if buttonDownFrame[b] == 0 {
		return 0
	}
	return pi.Frame - buttonDownFrame[b]
}

type Button string

const (
	Left  Button = "Left"
	Right Button = "Right"
)

var buttonTarget = pievent.NewTarget[EventButton]()
var buttonDebugTarget = pievent.NewTarget[EventButton]()
var moveTarget = pievent.NewTarget[EventMove]()
var moveDebugTarget = pievent.NewTarget[EventMove]()

var buttonDownFrame = map[Button]int{}

func init() {
	onButton := func(event EventButton, _ pievent.Handler) {
		switch event.Type {
		case EventButtonDown:
			buttonDownFrame[event.Button] = pi.Frame
		case EventButtonUp:
			buttonDownFrame[event.Button] = 0
		}
	}
	buttonTarget.SubscribeAll(onButton)
	buttonDebugTarget.SubscribeAll(onButton)

	onMove := func(event EventMove, _ pievent.Handler) {
		Position = event.Position
		MovementDelta = event.Position.Subtract(event.Previous)
	}
	moveTarget.SubscribeAll(onMove)
	moveDebugTarget.SubscribeAll(onMove)
}
