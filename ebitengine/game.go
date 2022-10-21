// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitengine

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elgopher/pi/vm"
)

type game struct {
	screenDataRGBA     []byte // reused RGBA pixels
	screenChanged      bool
	shouldSkipNextDraw bool
}

func (e *game) Update() error {
	updateStartedTime := time.Now()

	updateTime()
	updateController()
	updateMouse()
	updateKeyDuration()

	if vm.Update != nil {
		vm.Update()
	}

	if vm.GameLoopStopped {
		return gameStoppedErr
	}

	// Ebitengine treats Draw differently than π. In π Draw must be executed at most 30 times per second.
	// That's why π runs Draw() from inside Ebitengine's Update().
	if vm.Draw != nil {
		if e.shouldSkipNextDraw {
			e.shouldSkipNextDraw = false
			return nil
		}

		vm.Draw()

		elapsed := time.Since(updateStartedTime)
		if elapsed.Seconds() > 1/float64(tps) {
			e.shouldSkipNextDraw = true
		}
	}

	e.screenChanged = true

	return nil
}

func (e *game) Draw(screen *ebiten.Image) {
	// Ebitengine executes Draw based on display frequency.
	// But the screen is changed at most 30 times per second.
	// That's why there is no need to write pixels more often
	// than 30 times per second.
	if e.screenChanged {
		e.writeScreenPixels(screen)
		e.screenChanged = false
	}
}

func (e *game) writeScreenPixels(screen *ebiten.Image) {
	if e.screenDataRGBA == nil || len(e.screenDataRGBA)/4 != len(vm.ScreenData) {
		e.screenDataRGBA = make([]byte, len(vm.ScreenData)*4)
	}

	offset := 0
	for _, col := range vm.ScreenData {
		rgb := vm.Palette[vm.DisplayPalette[col]]
		e.screenDataRGBA[offset] = rgb.R
		e.screenDataRGBA[offset+1] = rgb.G
		e.screenDataRGBA[offset+2] = rgb.B
		e.screenDataRGBA[offset+3] = 0xff
		offset += 4
	}

	screen.WritePixels(e.screenDataRGBA)
}

func (e *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return vm.ScreenWidth, vm.ScreenHeight
}