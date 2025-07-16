// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"github.com/elgopher/pi"
	"github.com/hajimehoshi/ebiten/v2"
)

var RememberWindow bool

func RunOrErr() error {
	game := RunEbitenGame()

	setWindowSize := func() {
		monitor := ebiten.Monitor()

		width, height, minW, minH := WindowAutoSize(monitor)
		ebiten.SetWindowSize(width, height)
		ebiten.SetWindowSizeLimits(minW, minH, -1, -1)
	}
	setWindowSize()
	game.windowState.restore()

	// here we intentionally set only a subset of Ebiten parameters,
	// so the user can configure the rest as needed
	ebiten.SetTPS(pi.TPS())
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowClosingHandled(true)
	ebiten.SetWindowTitle("Pi")

	err := ebiten.RunGameWithOptions(game, &ebiten.RunGameOptions{SingleThread: true})
	if err != nil {
		panic(err)
	}

	return nil
}
