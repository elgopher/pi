// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"github.com/elgopher/pi/pievent"
	"github.com/elgopher/pi/pikey"
)

func handleInputInConsoleMode() {
	right := pikey.Duration(pikey.Right)
	if right > 0 {
		if right == 1 || right > 10 {
			showNextSnapshot()
		}
	} else {
		left := pikey.Duration(pikey.Left)
		if left == 1 || left > 10 {
			showPrevSnapshot()
		}
	}
}

func registerShortcuts() {
	// CTRL+SHIFT+I
	onCtrlShiftI := func() {
		if !consoleMode {
			enterConsoleMode()
		} else {
			exitConsoleMode()
		}
	}
	pikey.RegisterShortcut(onCtrlShiftI, pikey.Ctrl, pikey.Shift, pikey.I)

	// F12
	f12Down := pikey.Event{Type: pikey.EventDown, Key: pikey.F12}
	pikey.DebugTarget().Subscribe(f12Down, func(pikey.Event, pievent.Handler) {
		if consoleMode {
			captureSnapshot()
		}
	})

	// Space
	spaceDown := pikey.Event{Type: pikey.EventDown, Key: pikey.Space}
	pikey.DebugTarget().Subscribe(spaceDown, func(pikey.Event, pievent.Handler) {
		pauseOrResume()
	})

	// Esc
	escDown := pikey.Event{Type: pikey.EventDown, Key: pikey.Esc}
	pikey.DebugTarget().Subscribe(escDown, func(pikey.Event, pievent.Handler) {
		if consoleMode {
			exitConsoleMode()
		}
	})
}
