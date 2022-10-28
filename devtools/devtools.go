// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package devtools

import (
	"fmt"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/devtools/internal/inspector"
)

var (
	// BgColor is used to draw background behind the text
	BgColor byte = 1
	// FgColor is used to print text, draw icons etc.
	FgColor byte = 7
)

// MustRun runs the game with devtools.
//
// Any time you can pause them game by pressing F12. This will
// show screen inspector. F12 again resumes the game.
func MustRun(runBackend func() error) {
	update := pi.Update
	draw := pi.Draw

	inspector.BgColor, inspector.FgColor = BgColor, FgColor
	fmt.Println("Press F12 to pause the game and show devtools.")

	pi.Update = func() {
		updateDevTools()

		if !gamePaused && update != nil {
			update()
			handleStoppedGame()
		}
	}

	pi.Draw = func() {
		if !gamePaused && draw != nil {
			draw()
			handleStoppedGame()
		} else {
			drawDevTools()
		}
	}

	if err := runBackend(); err != nil {
		panic(fmt.Sprintf("Something terrible happened! Pi cannot be run: %v\n", err))
	}
}

func handleStoppedGame() {
	if pi.GameLoopStopped {
		pauseGame()
		pi.GameLoopStopped = false
	}
}
