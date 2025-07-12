// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piaudio

import "github.com/elgopher/pi/pievent"

// Event represents a type of piaudio-related event.
//
// Currently only EventTick is defined, but this can be extended in the future.
type Event string

const (
	// EventTick is published by the audio backend once per audio tick.
	EventTick Event = "tick"
)

// Target is the global event target for piaudio events.
//
// It allows other packages to subscribe to piaudio.Event notifications.
//
// Typically, EventTick will be published once per audio tick by the backend,
// enabling precise synchronization with audio scheduling.
var Target = pievent.NewTarget[Event]()
