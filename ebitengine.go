// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"errors"
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var gameStoppedErr = errors.New("game stopped")

func run() error {
	ebiten.SetMaxTPS(30)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(scrWidth*scale(), scrHeight*scale())
	ebiten.SetWindowSizeLimits(scrWidth, scrHeight, -1, -1)
	ebiten.SetWindowTitle("Pi Game")

	if err := ebiten.RunGame(&ebitengineGame{}); err != nil {
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

type ebitengineGame struct {
	screenDataRGBA []byte // reused RGBA pixels
}

func (e *ebitengineGame) Update() error {
	updateTime()
	updateController()

	if Update != nil {
		Update()
	}

	if gameLoopStopped {
		return gameStoppedErr
	}

	return nil
}

func (e *ebitengineGame) Draw(screen *ebiten.Image) {
	if Draw != nil {
		Draw()
	}

	e.replaceScreenPixels(screen)
}

func (e *ebitengineGame) replaceScreenPixels(screen *ebiten.Image) {
	if e.screenDataRGBA == nil || len(e.screenDataRGBA)/4 != len(ScreenData) {
		e.screenDataRGBA = make([]byte, len(ScreenData)*4)
	}

	offset := 0
	for _, col := range ScreenData {
		rgb := Palette[displayPalette[col]]
		e.screenDataRGBA[offset] = rgb.R
		e.screenDataRGBA[offset+1] = rgb.G
		e.screenDataRGBA[offset+2] = rgb.B
		e.screenDataRGBA[offset+3] = 0xff
		offset += 4
	}

	screen.ReplacePixels(e.screenDataRGBA)
}

func (e *ebitengineGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return scrWidth, scrHeight
}
