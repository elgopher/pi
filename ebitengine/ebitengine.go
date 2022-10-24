// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitengine

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elgopher/pi/mem"
)

var gameStoppedErr = errors.New("game stopped")

const tps = 30

// Backend opens the window and runs the game loop. It must be
// called from the main thread.
func Backend() error {
	lastTime = time.Now()

	ebiten.SetTPS(tps)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(mem.ScreenWidth*scale(), mem.ScreenHeight*scale())
	ebiten.SetWindowSizeLimits(mem.ScreenWidth, mem.ScreenHeight, -1, -1)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetWindowTitle("Pi Game")

	if err := ebiten.RunGame(&game{}); err != nil {
		if err == gameStoppedErr {
			return nil
		}

		return fmt.Errorf("running game using Ebitengine failed: %w", err)
	}

	return nil
}

func scale() int {
	return int(math.Round(ebiten.DeviceScaleFactor() * 3))
}
