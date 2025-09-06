// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package main

import (
	_ "embed"
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/piebiten"
	"github.com/elgopher/pi/pikey"
	"github.com/elgopher/pi/piscope"
)

//go:embed "gamepad.png"
var gamepadPNG []byte

func main() {
	pi.SetScreenSize(100, 60)
	pi.Palette = pi.DecodePalette(gamepadPNG)
	canvas := pi.DecodeCanvas(gamepadPNG)
	spr := pi.SpriteFrom(canvas, 58, 26, 9, 9)
	piscope.Start()
	posx := 0
	posy := 0
	pi.Update = func() {
		if pikey.Duration(pikey.Left) > 0 {
			posx--
		}
		if pikey.Duration(pikey.Right) > 0 {
			posx++
		}
		if pikey.Duration(pikey.Up) > 0 {
			posy--
		}
		if pikey.Duration(pikey.Down) > 0 {
			posy++
		}
	}
	pi.Draw = func() {
		pi.DrawSprite(spr, posx, posy)
	}
	piebiten.Run()
}
