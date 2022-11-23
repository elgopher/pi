// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package ebitengine uses Ebitengine technology to run the game.
package ebitengine

import (
	"errors"
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elgopher/pi"
)

var gameStoppedErr = errors.New("game stopped")

const tps = 30

// Run opens the window and runs the game loop. It must be
// called from the main thread.
func Run() error {
	screen := pi.Scr() // TODO Screen size should be read each frame (and window size adjusted)

	ebiten.SetTPS(tps)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(screen.Width()*scale(), screen.Height()*scale())
	ebiten.SetWindowSizeLimits(screen.Width(), screen.Height(), -1, -1)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetWindowTitle("Pi Game")

	//var theAudio otoAudio
	//if err := theAudio.Start(); err != nil {
	//	return err
	//}
	//defer theAudio.Stop()

	theGame := &game{}

	if err := ebiten.RunGame(theGame); err != nil {
		if err == gameStoppedErr {
			return nil
		}

		return fmt.Errorf("running game using Ebitengine failed: %w", err)
	}

	return nil
}

// MustRun does the same as Run, but panics instead of returning an error.
//
// Useful for writing unit tests and quick and dirty prototypes. Do not use on production ;)
func MustRun() {
	if err := Run(); err != nil {
		panic(fmt.Sprintf("Something terrible happened! Pi cannot be run: %v\n", err))
	}
}

func scale() int {
	return int(math.Round(ebiten.DeviceScaleFactor() * 3))
}
