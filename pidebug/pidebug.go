// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package pidebug provides an API for pausing the game for debugging purposes.
// It is essential for building developer tools.
package pidebug

var paused = false

// SetPaused pauses or resumes the game.
//
// The game does not pause immediately, but only at the end of the current game
// loop iteration.
//
// When the game is paused, most targets stop processing events (for example,
// piloop.Target or pikey.Target).
// Developer tools should use their debug-mode equivalents instead, such as
// piloop.DebugTarget or pikey.DebugTarget.
func SetPaused(p bool) bool {
	prev := paused
	paused = p
	if !prev && p {
		target.Publish(EventPaused)
	}
	if prev && !p {
		target.Publish(EventResumed)
	}
	return prev
}

func Paused() bool {
	return paused
}
