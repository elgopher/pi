// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pikey

import (
	"slices"

	"github.com/elgopher/pi/pievent"
	"github.com/elgopher/pi/piloop"
)

// RegisterShortcut registers a new keyboard shortcut.
//
// The function takes a callback onEvent, which will be called
// every time the user presses all the keys specified in `keys` simultaneously.
// The order of keys in the `keys` argument does not matter.
//
// It returns a *Shortcut object, which can later be used to unregister the shortcut.
//
// Example:
//
//	RegisterShortcut(func() {
//		log.Println("Pressed Ctrl + S")
//	}, KeyCtrl, KeyS)
//
// Notes:
//   - If the same key combination is registered multiple times,
//     all corresponding callbacks will be invoked.
func RegisterShortcut(onEvent func(), keys ...Key) *Shortcut {
	s := &Shortcut{
		keys:        keys,
		pressedKeys: make(map[Key]int),
		onEvent:     onEvent,
	}
	s.keyHandler = Target().SubscribeAll(s.onKey)
	s.updateHandler = piloop.Target().Subscribe(piloop.EventLateUpdate, s.onEventLateUpdate)
	return s
}

type Shortcut struct {
	keys    []Key
	onEvent func()

	pressedKeys map[Key]int

	keyHandler    pievent.Handler
	updateHandler pievent.Handler
}

func (s *Shortcut) Unregister() {
	Target().Unsubscribe(s.keyHandler)
	piloop.Target().Unsubscribe(s.updateHandler)
}

func (s *Shortcut) onKey(e Event, _ pievent.Handler) {
	switch e.Type {
	case EventDown:
		if slices.Contains(s.keys, e.Key) {
			s.pressedKeys[e.Key] = 1
		}
	case EventUp:
		if slices.Contains(s.keys, e.Key) {
			s.pressedKeys[e.Key] = 0
		}
	}
}

func (s *Shortcut) onEventLateUpdate(piloop.Event, pievent.Handler) {
	if s.justPressed() {
		s.onEvent()
	}
	for key, duration := range s.pressedKeys {
		if duration > 0 {
			s.pressedKeys[key] = duration + 1
		}
	}
}

func (s *Shortcut) justPressed() bool {
	if len(s.pressedKeys) < len(s.keys) {
		return false
	}

	for _, duration := range s.pressedKeys {
		if duration == 0 {
			return false
		}
	}

	for _, duration := range s.pressedKeys {
		if duration == 1 {
			return true
		}
	}

	return false
}
