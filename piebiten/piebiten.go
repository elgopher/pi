// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package piebiten enables running your game using the [Ebitengine] backend.
//
// Ebitengine is a cross-platform game engine that supports Windows, macOS,
// Linux, FreeBSD, web browsers, Android, iOS, and even Nintendo Switch.
//
// To launch your game, use [Run] or [RunOrErr].
//
// This package also provides advanced functions for integrating Pi
// with your own Ebitengine-based game, such as [CopyCanvasToEbitenImage].
//
// [Ebitengine]: https://ebitengine.org
package piebiten

import (
	"errors"
	"github.com/elgopher/pi/piaudio"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"strconv"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/piebiten/internal"
)

// RememberWindow determines whether the game should open
// at its last window position, size, and monitor when set to true
var RememberWindow = false

// Run starts the Ebitengine backend. It panics if something goes wrong.
//
// If you want to handle errors gracefully, use [RunOrErr] instead.
//
// This function must be called from the first goroutine (the main thread).
func Run() {
	if err := RunOrErr(); err != nil {
		panic("piebiten.Run failed: " + err.Error())
	}
}

// RunOrErr starts the Ebitengine backend and returns an error if something goes wrong.
//
// This function must be called from the first goroutine (the main thread).
func RunOrErr() error {
	if internal.CurrentGoroutineID() != 1 {
		return errors.New("must be run from main goroutine 1")
	}
	internal.RememberWindow = RememberWindow
	return internal.RunOrErr() //nolint:wrapcheck
}

// CopyCanvasToEbitenImage copies the canvas to dst using the current
// palette in pi.Palette and the palette mapping in pi.PaletteMapping.
func CopyCanvasToEbitenImage(canvas pi.Canvas, dst *ebiten.Image) {
	internal.CopyCanvasToEbitenImage(canvas, dst)
}

// StartAudioBackend starts the audio backend with the given Ebitengine audio.Context.
// Use if you want only piaudio functionality without Pi's graphics.
//
// audio.Context must have a sample rate of 44100.
func StartAudioBackend(ctx *audio.Context) Audio {
	if ctx.SampleRate() != internal.CtxSampleRate {
		panic("piebiten.StartAudioBackend: audio.Context must have " + strconv.Itoa(internal.CtxSampleRate) + " sample rate")
	}
	return internal.StartAudioBackend(ctx)
}

type Audio interface {
	piaudio.BackendInterface
	// OnBeforeUpdate must be called at the start of Ebitengine's Update function.
	OnBeforeUpdate()
	// OnAfterUpdate must be called at the end of Ebitengine's Update function.
	OnAfterUpdate()
}
