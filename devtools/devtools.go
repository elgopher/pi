// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package devtools

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/vm"
)

var (
	// BgColor is used to draw background behind the text
	BgColor byte = 1
	// FgColor is used to print text, draw icons etc.
	FgColor byte = 7
)

// MustRun runs the game using backend, similarly to pi.MustRun.
//
// Any time you can pause them game by pressing F12. This will
// show screen inspector. F12 again resumes the game.
func MustRun(runBackend func() error) {
	update := pi.Update
	draw := pi.Draw

	pi.Update = func() {
		updateDevTools()

		if !paused && update != nil {
			update()
			handleStoppedGame()
		}
	}

	pi.Draw = func() {
		if !paused && draw != nil {
			draw()
			handleStoppedGame()
		} else {
			drawDevTools()
		}
	}

	pi.MustRun(runBackend)
}

func handleStoppedGame() {
	if vm.GameLoopStopped {
		pauseGame()
		vm.GameLoopStopped = false
	}
}
