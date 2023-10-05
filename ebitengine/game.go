// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitengine

import (
	_ "embed"
	"sync/atomic"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/image"
	"github.com/elgopher/pi/key"
)

type game struct {
	ready              atomic.Bool
	shouldSkipNextDraw bool
	screenFrame        screenFrame
}

type screenFrame struct {
	changed        bool
	pix            []byte
	palette        [256]image.RGB
	pald           pi.PalMapping
	screenDataRGBA []byte // reused RGBA pixels
}

func (e *game) Update() error {
	if !e.ready.Load() {
		return nil
	}

	updateStartedTime := time.Now()

	updateTime()
	updateController()
	updateMouse()
	updateKeyDuration()
	handleKeyboardShortcuts()

	if pi.Update != nil {
		pi.Update()
	}

	if pi.GameLoopStopped {
		return gameStoppedErr
	}

	// Ebitengine treats Draw differently than π. In π Draw must be executed at most 30 times per second.
	// That's why π runs Draw() from inside Ebitengine's Update().
	if pi.Draw != nil {
		if e.shouldSkipNextDraw {
			e.shouldSkipNextDraw = false
			return nil
		}

		pi.Draw()

		elapsed := time.Since(updateStartedTime)
		if elapsed.Seconds() > 1/float64(tps) {
			e.shouldSkipNextDraw = true
		}
	}

	e.screenFrame.update()

	return nil
}

func (f *screenFrame) update() {
	scrPix := pi.Scr().Pix()
	screenChanged := !slicesEqual(f.pix, scrPix) || f.palette != pi.Palette || f.pald != pi.Pald
	if screenChanged {
		f.changed = true

		if len(f.pix) != len(scrPix) {
			f.pix = make([]byte, len(scrPix))
		}
		copy(f.pix, scrPix)
		f.pald = pi.Pald
		f.palette = pi.Palette
	}
}

// replace with slices.Equal after upgrading to go 1.21
func slicesEqual[S ~[]E, E comparable](s1, s2 S) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func handleKeyboardShortcuts() {
	f11 := key.Duration[key.F11] == 1
	altEnter := key.Duration[key.Enter] == 1 && key.Duration[key.Alt] > 0
	if f11 || altEnter {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
}

func (e *game) Draw(screen *ebiten.Image) {
	if !e.ready.Load() {
		e.drawNotReady(screen)
		return
	}

	e.screenFrame.writeScreenPixels(screen)
}

func (f *screenFrame) writeScreenPixels(screen *ebiten.Image) {
	if f.changed {
		f.changed = false

		pix := pi.Scr().Pix()
		if f.screenDataRGBA == nil || len(f.screenDataRGBA)/4 != len(pix) {
			f.screenDataRGBA = make([]byte, len(pix)*4)
		}

		offset := 0
		for _, col := range pix {
			rgb := pi.Palette[pi.Pald[col]]
			f.screenDataRGBA[offset] = rgb.R
			f.screenDataRGBA[offset+1] = rgb.G
			f.screenDataRGBA[offset+2] = rgb.B
			f.screenDataRGBA[offset+3] = 0xff
			offset += 4
		}

		screen.WritePixels(f.screenDataRGBA)
	}
}

func (e *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	scr := pi.Scr()
	return scr.Width(), scr.Height()
}
