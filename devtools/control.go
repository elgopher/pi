// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package devtools

import (
	"fmt"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/devtools/internal/snapshot"
)

var (
	gamePaused       bool
	pauseOnNextFrame bool
	timeWhenPaused   float64
)

var helpShown bool

func pauseGame() {
	fmt.Println("Game paused")
	if !helpShown {
		helpShown = true
		fmt.Println("\nPress right mouse button in the game window to show the toolbar.")
		fmt.Println("Press P in the game window to take screenshot.")
		fmt.Println("Press N in the game window to go to next frame.")
		fmt.Println("Press F12 in the game window to resume the game and exit devtools inspector.")
	}
	gamePaused = true
	timeWhenPaused = pi.Time
	snapshot.Take()
}

func resumeGame() {
	gamePaused = false
	pi.Time = timeWhenPaused
	snapshot.Draw()
	fmt.Println("Game resumed")
}

func resumeUntilNextFrame() {
	if gamePaused {
		resumeGame()
		pauseOnNextFrame = true
	} else {
		fmt.Println("Game not paused")
	}
}
