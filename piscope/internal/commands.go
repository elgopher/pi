// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/pidebug"
	"github.com/elgopher/pi/pisnap"
	"log"
)

func enterConsoleMode() {
	log.Println("Entering console")
	prev := pi.SetColor(*bgColor)
	pi.Rect(0, 0, pi.Screen().W()-1, pi.Screen().H()-1)
	pi.SetColor(prev)
	consoleMode = true
}

func exitConsoleMode() {
	log.Println("Exiting console")
	theScreenRecorder.ShowPrev()
	theScreenRecorder.Reset()
	consoleMode = false
	pidebug.SetPaused(false)
}

func captureSnapshot() {
	f, err := pisnap.CaptureOrErr()
	if err != nil {
		log.Println("Error capturing screenshot:", err)
	} else {
		log.Println("Screenshot saved to", f)
	}
}

func showPrevSnapshot() {
	pidebug.SetPaused(true)
	theScreenRecorder.ShowPrev()
}

func showNextSnapshot() {
	if !theScreenRecorder.ShowNext() {
		pidebug.SetPaused(false)
		pauseOnNextFrame = true
	}
}

func pauseOrResume() {
	if consoleMode {
		theScreenRecorder.GoToLast()

		pidebug.SetPaused(!pidebug.Paused())
	}
}
