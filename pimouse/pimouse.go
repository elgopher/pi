// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pimouse

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/internal/input"
	"github.com/elgopher/pi/pievent"
)

var (
	Position      pi.Position
	MovementDelta pi.Position // mouse movement delta since the last frame
)

func Duration(b Button) int {
	return buttonState.Duration(b)
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

var buttonState input.State[Button]

func init() {
	onButton := func(event EventButton, _ pievent.Handler) {
		switch event.Type {
		case EventButtonDown:
			buttonState.SetDownFrame(event.Button, pi.Frame)
		case EventButtonUp:
			buttonState.SetUpFrame(event.Button, pi.Frame)
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
