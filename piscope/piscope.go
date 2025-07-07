// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package piscope provides developer tools.
package piscope

import (
	"log"
	"os"
	"runtime/pprof"
	"strconv"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/picofont"
	"github.com/elgopher/pi/pidebug"
	"github.com/elgopher/pi/pievent"
	"github.com/elgopher/pi/pikey"
	"github.com/elgopher/pi/piloop"
	"github.com/elgopher/pi/pimouse"
)

const bgColor = 11

// Start launches the developer tools.
//
// Pressing Ctrl+Shift+I will activate the tools in the game
func Start() {
	smallFont := picofont.Sheet
	var consoleMode = false
	var mouseTracing = false
	var keyTracing = false

	onCtrlShiftI := func() {
		if !consoleMode {
			log.Println("Entering console")

			prev := pi.SetColor(bgColor)
			pi.Rect(0, 0, pi.Screen().W()-1, pi.Screen().H()-1)
			pi.SetColor(prev)
		} else {
			log.Println("Exiting console")
			theScreenRecorder.Reset()

			pi.Screen().Clear(0)
		}
		consoleMode = !consoleMode
	}
	pikey.RegisterShortcut(onCtrlShiftI, pikey.Control, pikey.Shift, pikey.I)

	pauseOnNextFrame := false

	piloop.DebugTarget().Subscribe(piloop.EventUpdate, func(piloop.Event, pievent.Handler) {
		if consoleMode {
			if !pidebug.Paused() {
				theScreenRecorder.Save()
			}
			if pauseOnNextFrame {
				pidebug.SetPaused(true)
				pauseOnNextFrame = false
			}
			if pikey.Duration(pikey.P) == 1 {
				pidebug.SetPaused(!pidebug.Paused())
			}

			right := pikey.Duration(pikey.Right)
			if right > 0 {
				if right == 1 || right > 10 {
					if !theScreenRecorder.ShowNext() {
						pidebug.SetPaused(false)
						pauseOnNextFrame = true
					}
				}
			} else {
				left := pikey.Duration(pikey.Left)
				if left == 1 || left > 10 {
					pidebug.SetPaused(true)
					theScreenRecorder.ShowPrev()
				}
			}

			if pikey.Duration(pikey.D) == 1 { // memory dump
				f, _ := os.Create("heap.prof")
				if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
					log.Println("Error writing heap profile:", err)
				}
				_ = f.Close()
				log.Println("Memory dumped to heap.prof")
			}

			if pikey.Duration(pikey.M) == 1 {
				mouseTracing = !mouseTracing
				if mouseTracing {
					log.Println("Mouse tracing enabled")
				} else {
					log.Println("Mouse tracing disabled")
				}
				pimouse.ButtonTarget().SetTracing(mouseTracing)
			}
			if pikey.Duration(pikey.K) == 1 {
				keyTracing = !keyTracing
				if keyTracing {
					log.Println("Keyboard tracing enabled")
				} else {
					log.Println("Keyboard tracing disabled")
				}
				pikey.Target().SetTracing(keyTracing)
			}
		}
	})

	piloop.DebugTarget().Subscribe(piloop.EventLateDraw, func(piloop.Event, pievent.Handler) {
		if consoleMode {
			prev := pi.SetColor(bgColor)
			defer pi.SetColor(prev)

			y := pi.DrawTarget().H() - smallFont.Height - 2
			pi.RectFill(0, y, pi.Screen().W()-1, pi.Screen().H()-1)

			pixelColor := pi.GetPixel(pimouse.Position.X, pimouse.Position.Y)
			if pixelColor != bgColor {
				pi.SetColor(pixelColor)
			} else {
				pi.SetColor(0)
			}
			smallFont.Print(strconv.Itoa(int(pixelColor)), pi.DrawTarget().W()-30, y)
		}
	})
}
