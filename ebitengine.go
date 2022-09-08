// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var gameStoppedErr = errors.New("game stopped")

const tps = 30

func run() error {
	ebiten.SetTPS(tps)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(scrWidth*scale(), scrHeight*scale())
	ebiten.SetWindowSizeLimits(scrWidth, scrHeight, -1, -1)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
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
	screenDataRGBA     []byte // reused RGBA pixels
	screenChanged      bool
	shouldSkipNextDraw bool
}

func (e *ebitengineGame) Update() error {
	updateStartedTime := time.Now()

	updateTime()
	updateController()
	updateMouse()

	if Update != nil {
		Update()
	}

	if gameLoopStopped {
		return gameStoppedErr
	}

	// Ebitengine treats Draw differently than π. In π Draw must be executed at most 30 times per second.
	// That's why π runs Draw() from inside Ebitengine's Update().
	if Draw != nil {
		if e.shouldSkipNextDraw {
			e.shouldSkipNextDraw = false
			return nil
		}

		Draw()

		elapsed := time.Since(updateStartedTime)
		if elapsed.Seconds() > 1/float64(tps) {
			e.shouldSkipNextDraw = true
		}
	}

	e.screenChanged = true

	return nil
}

func (e *ebitengineGame) Draw(screen *ebiten.Image) {
	// Ebitengine executes Draw based on display frequency.
	// But the screen is changed at most 30 times per second.
	// That's why there is no need to write pixels more often
	// than 30 times per second.
	if e.screenChanged {
		e.writeScreenPixels(screen)
		e.screenChanged = false
	}
}

func (e *ebitengineGame) writeScreenPixels(screen *ebiten.Image) {
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

	screen.WritePixels(e.screenDataRGBA)
}

func (e *ebitengineGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return scrWidth, scrHeight
}
