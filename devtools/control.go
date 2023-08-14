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
	stateCopy        stateWhenGameWasPaused
)

type stateWhenGameWasPaused struct {
	time   float64
	camera pi.Position
	clip   pi.Region
}

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

	stateCopy = stateWhenGameWasPaused{
		time:   pi.Time,
		camera: pi.Camera,
		clip:   pi.Scr().Clip(),
	}

	pi.Camera.Reset()
	pi.ClipReset()

	snapshot.Take()
}

func resumeGame() {
	gamePaused = false

	pi.Time = stateCopy.time
	pi.Camera = stateCopy.camera
	pi.Clip(stateCopy.clip.X, stateCopy.clip.Y, stateCopy.clip.W, stateCopy.clip.H)

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
