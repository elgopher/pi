// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"github.com/elgopher/pi/piaudio"
	"github.com/elgopher/pi/piebiten/internal/audio"
	ebitenaudio "github.com/hajimehoshi/ebiten/v2/audio"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/pidebug"
	"github.com/elgopher/pi/pievent"
	"github.com/elgopher/pi/piloop"
)

func RunEbitenGame() *EbitenGame {
	screen := pi.Screen()

	ctx := ebitenaudio.NewContext(audio.CtxSampleRate)
	theAudioBackend := audio.StartAudioBackend(ctx)
	piaudio.Backend = theAudioBackend

	game := &EbitenGame{
		piScreen:       screen,
		ebitenScreen:   ebiten.NewImage(screen.W(), screen.H()),
		drawScreenOpts: &ebiten.DrawImageOptions{},
		audioBackend:   theAudioBackend,
	}

	pidebug.Target().SubscribeAll(game.onPidebugEvent)

	return game
}

// get the monitor once per second and cache it.
//
// This is not a strong optimization, since Ebitengine itself
// performs syscalls on each call to ebiten's Draw.
type cachedMonitor struct {
	monitor       *ebiten.MonitorType
	lastCheckTime time.Time
}

func (c *cachedMonitor) Get() *ebiten.MonitorType {
	if time.Since(c.lastCheckTime) > time.Second {
		c.monitor = ebiten.Monitor()
		c.lastCheckTime = time.Now()
	}
	return c.monitor
}

// TODO split into multiple objects to reduce complexity
type EbitenGame struct {
	piScreen       pi.Canvas
	ebitenScreen   *ebiten.Image
	drawScreenOpts *ebiten.DrawImageOptions
	keys           []ebiten.Key
	mousePosition  pi.Position
	cachedMonitor  cachedMonitor
	started        bool

	scale     float64
	left, top float64

	// When true, indicates that the frame was rendered by pi.Draw
	// and should be displayed on the next Draw call.
	dirty bool

	// When true, indicates that the last update+draw cycle
	// exceeded the tick duration of 1/TPS (e.g., 33 ms for TPS=30).
	skipNextDraw bool

	paused bool

	gamepads     ebitenGamepads
	windowState  windowState
	audioBackend *audio.Backend

	ebitenFrame int // frame incremented on each Ebiten tick
}

func (g *EbitenGame) Update() error {
	if ebiten.IsWindowBeingClosed() {
		piloop.Target().Publish(piloop.EventWindowClose)
		return ebiten.Termination
	}

	g.windowState.store()

	g.audioBackend.OnBeforeUpdate()

	started := time.Now()

	if !g.started {
		if pi.Init != nil {
			pi.Init()
		}
		piloop.Target().Publish(piloop.EventInit)
	}
	g.started = true

	if g.ebitenFrame%(ebitenTPS/pi.TPS()) == 0 {
		if !g.paused {
			piloop.Target().Publish(piloop.EventFrameStart)
		}
		piloop.DebugTarget().Publish(piloop.EventFrameStart)
	}

	g.updateMouse()
	g.updateKeyboard()
	g.gamepads.update()

	if g.ebitenFrame%(ebitenTPS/pi.TPS()) == (ebitenTPS/pi.TPS())-1 {
		if !g.paused {
			pi.Update()
			piloop.Target().Publish(piloop.EventUpdate)
		}
		piloop.DebugTarget().Publish(piloop.EventUpdate)

		if !g.paused {
			piloop.Target().Publish(piloop.EventLateUpdate)
		}
		piloop.DebugTarget().Publish(piloop.EventLateUpdate)

		if !g.skipNextDraw {
			if !g.paused {
				pi.Draw()
				piloop.Target().Publish(piloop.EventDraw)
			}
			piloop.DebugTarget().Publish(piloop.EventDraw)

			if !g.paused {
				piloop.Target().Publish(piloop.EventLateDraw)
			}
			piloop.DebugTarget().Publish(piloop.EventLateDraw)

			g.dirty = true
		} else {
			g.skipNextDraw = false
		}

		if time.Since(started).Seconds() > 1/float64(pi.TPS()) {
			g.skipNextDraw = true // game is too slow. Try to keep up by discarding next pi.Draw()
		}

		pi.Time += 1.0 / float64(pi.TPS())
		pi.Frame++
	}

	g.audioBackend.OnAfterUpdate()

	g.ebitenFrame++

	return nil
}

func (g *EbitenGame) Draw(screen *ebiten.Image) {
	if g.dirty { // draw only when needed to avoid CPU load on monitors >30 Hz
		g.dirty = false

		CopyCanvasToEbitenImage(pi.Screen(), g.ebitenScreen)

		screen.DrawImage(g.ebitenScreen, g.drawScreenOpts)
	}
}

func (g *EbitenGame) LayoutF(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64) {
	piScrW, piScrH := float64(g.piScreen.W()), float64(g.piScreen.H())

	monitor := g.cachedMonitor.Get()
	deviceScaleFactor := monitor.DeviceScaleFactor()
	realWith := outsideWidth * deviceScaleFactor
	realHeight := outsideHeight * deviceScaleFactor
	widthRatio := realWith / piScrW
	heightRatio := realHeight / piScrH
	scale := math.Floor(min(widthRatio, heightRatio))

	screenWidth = realWith
	screenHeight = realHeight

	g.scale = scale
	// center on screen:
	g.left = (realWith - piScrW*scale) / 2.0
	g.top = (realHeight - piScrH*scale) / 2.0

	g.drawScreenOpts.GeoM.Reset()
	g.drawScreenOpts.GeoM.Scale(g.scale, g.scale)
	g.drawScreenOpts.GeoM.Translate(g.left, g.top)

	return
}

// this method is not executed because LayoutF is
func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return
}

func (g *EbitenGame) Resize() {
	screen := pi.Screen()
	g.piScreen = screen
	g.ebitenScreen = ebiten.NewImage(screen.W(), screen.H())
}

func (g *EbitenGame) onPidebugEvent(event pidebug.Event, _ pievent.Handler) {
	switch event {
	case pidebug.EventPause:
		g.paused = true
	case pidebug.EventResume:
		g.paused = false
	}
}
