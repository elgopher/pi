// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pidebug

import "github.com/elgopher/pi/pievent"

type Event string

const (
	EventPaused  Event = "paused"
	EventResumed Event = "resumed"
)

func Target() pievent.Target[Event] {
	return target
}

var target = pievent.NewTarget[Event]()
