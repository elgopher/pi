// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pimouse

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/pievent"
)

// EventButton is published when the player presses or releases a mouse button.
// It may be published more than once during a single game tick.
type EventButton struct {
	Type   EventButtonType
	Button Button
}

type EventButtonType string

const (
	EventButtonUp   EventButtonType = "up"
	EventButtonDown EventButtonType = "down"
)

// EventMove is published when the mouse is moved.
//
// It may be published more than once during a single game tick.
type EventMove struct {
	Position pi.Position
	Previous pi.Position
}

func ButtonTarget() pievent.Target[EventButton] {
	return buttonTarget
}

func ButtonDebugTarget() pievent.Target[EventButton] {
	return buttonDebugTarget
}

func MoveTarget() pievent.Target[EventMove] {
	return moveTarget
}

func MoveDebugTarget() pievent.Target[EventMove] {
	return moveDebugTarget
}
