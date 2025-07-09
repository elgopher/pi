// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	_ "embed"
	"github.com/elgopher/pi"
)

//go:embed "icons.png"
var iconsPNG []byte

var icons = struct {
	AlignTop    pi.Sprite
	AlignBottom pi.Sprite
	Screen      pi.Sprite
	Palette     pi.Sprite
	ColorTables pi.Sprite
	Variables   pi.Sprite
	Paint       pi.Sprite
	Separator   pi.Sprite
	Snap        pi.Sprite
	Prev        pi.Sprite
	Pause       pi.Sprite
	Play        pi.Sprite
	Next        pi.Sprite
	Exit        pi.Sprite
}{}

func init() {
	// TODO Decode palette
	iconsSheet := pi.DecodeCanvas(iconsPNG)
	icons.AlignTop = pi.SpriteFrom(iconsSheet, 0, 0, 8, 8)
	icons.AlignBottom = pi.SpriteFrom(iconsSheet, 8, 0, 8, 8)
	icons.Screen = pi.SpriteFrom(iconsSheet, 16, 0, 8, 8)
	icons.Palette = pi.SpriteFrom(iconsSheet, 24, 0, 8, 8)
	icons.ColorTables = pi.SpriteFrom(iconsSheet, 32, 0, 8, 8)
	icons.Variables = pi.SpriteFrom(iconsSheet, 40, 0, 8, 8)
	icons.Paint = pi.SpriteFrom(iconsSheet, 48, 0, 8, 8)
	icons.Separator = pi.SpriteFrom(iconsSheet, 58, 0, 4, 8)
	icons.Snap = pi.SpriteFrom(iconsSheet, 64, 0, 8, 8)
	icons.Prev = pi.SpriteFrom(iconsSheet, 74, 0, 5, 8)
	icons.Pause = pi.SpriteFrom(iconsSheet, 82, 0, 5, 8)
	icons.Play = pi.SpriteFrom(iconsSheet, 90, 0, 5, 8)
	icons.Next = pi.SpriteFrom(iconsSheet, 98, 0, 5, 8)
	icons.Exit = pi.SpriteFrom(iconsSheet, 112, 0, 8, 8)
}
