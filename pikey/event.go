// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pikey

import "github.com/elgopher/pi/pievent"

type Event struct {
	Type EventType
	Key  Key
}

type EventType string

const (
	EventUp   EventType = "up"
	EventDown EventType = "down"
)

func Target() pievent.Target[Event] {
	return target
}

// events are published all the time - even when game is paused.
func DebugTarget() pievent.Target[Event] {
	return debugTarget
}
