// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package piloop defines events published during the game loop.
//
// It enables adding logic from any component,
// including those created by third parties.
package piloop

import "github.com/elgopher/pi/pievent"

func Target() pievent.Target[Event] {
	return target
}

func DebugTarget() pievent.Target[Event] {
	return debugTarget
}

var (
	target      = pievent.NewTarget[Event]()
	debugTarget = pievent.NewTarget[Event]()
)
