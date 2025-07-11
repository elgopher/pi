// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"fmt"
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/picofont"
	"github.com/elgopher/pi/pidebug"
	"github.com/elgopher/pi/pievent"
	"github.com/elgopher/pi/pigui"
	"github.com/elgopher/pi/piloop"
	"github.com/elgopher/pi/pimouse"
)

var bgColor, fgColor *pi.Color

var consoleMode, pauseOnNextFrame bool

// Start launches the developer tools.
//
// Obecnie piscope wymaga, żeby gra miała rozdzielczość conajmniej 128
// pikseli w poziomie oraz 16 pikseli w pionie. Dodatkowa paleta gry musi
// używać conajmniej 2 kolorów.
//
// Pressing Ctrl+Shift+I will activate the tools in the game
func Start(backgroundColor, foregroundColor *pi.Color) {
	bgColor = backgroundColor
	fgColor = foregroundColor

	// TODO Handle screen size change event and redraw entire gui.

	smallFont := picofont.Sheet

	registerShortcuts()

	gui := pigui.New()
	attachToolbar(gui)

	piloop.DebugTarget().Subscribe(piloop.EventUpdate, func(piloop.Event, pievent.Handler) {
		if consoleMode {
			gui.Update()

			if !pidebug.Paused() {
				theScreenRecorder.Save()
			}

			if pauseOnNextFrame {
				pidebug.SetPaused(true)
				pauseOnNextFrame = false
			}

			handleInputInConsoleMode()
		}
	})

	piloop.DebugTarget().Subscribe(piloop.EventLateDraw, func(piloop.Event, pievent.Handler) {
		if consoleMode {
			gui.Draw()

			screen := pi.Screen()

			y := screen.H() - smallFont.Height - 1

			prev := pi.SetColor(*bgColor)
			defer pi.SetColor(prev)

			pixelColor := pi.GetPixel(pimouse.Position.X, pimouse.Position.Y)
			if pixelColor != *bgColor {
				pi.SetColor(pixelColor)
			} else {
				pi.SetColor(1)
			}
			msg := fmt.Sprintf("%d(%d,%d)", pixelColor, pimouse.Position.X, pimouse.Position.Y)
			smallFont.Print(msg, 50, y+2)
		}
	})
}
