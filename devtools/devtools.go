// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package devtools

import (
	"fmt"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/devtools/internal/inspector"
	"github.com/elgopher/pi/devtools/internal/terminal"
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
	if update == nil {
		update = func() {}
	}

	draw := pi.Draw
	if draw == nil {
		draw = func() {}
	}

	inspector.BgColor, inspector.FgColor = BgColor, FgColor
	fmt.Println("Press F12 in the game window to pause the game and activate devtools inspector.")
	fmt.Println("Terminal activated. Type help for help.")

	pi.Update = func() {
		updateDevTools()

		if !gamePaused {
			update()
			handleStoppedGame()
		}
	}

	pi.Draw = func() {
		if !gamePaused {
			draw()
			handleStoppedGame()
		} else {
			drawDevTools()
		}
	}

	if err := interpreterInstance.SetUpdate(&update); err != nil {
		panic(fmt.Sprintf("problem exporting Update function: %s", err))
	}

	if err := interpreterInstance.SetDraw(&draw); err != nil {
		panic(fmt.Sprintf("problem exporting Draw function: %s", err))
	}

	terminal.StartReadingCommands()
	defer terminal.StopReadingCommandsFromStdin()

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
