// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package main

import (
	_ "embed"
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/picofont"
	"github.com/elgopher/pi/piebiten"
	"github.com/elgopher/pi/pipad"
)

//go:embed "gamepad.png"
var gamepadPNG []byte

func main() {
	pi.SetScreenSize(85, 60)
	pi.Palette = pi.DecodePalette(gamepadPNG)
	gamepad := pi.DecodeCanvas(gamepadPNG)

	buttonSprites := map[pipad.Button]pi.Sprite{
		pipad.X:      pi.SpriteFrom(gamepad, 48, 18, 9, 9),
		pipad.Y:      pi.SpriteFrom(gamepad, 58, 10, 9, 9),
		pipad.B:      pi.SpriteFrom(gamepad, 68, 18, 9, 9),
		pipad.A:      pi.SpriteFrom(gamepad, 58, 26, 9, 9),
		pipad.Left:   pi.SpriteFrom(gamepad, 11, 19, 8, 8),
		pipad.Right:  pi.SpriteFrom(gamepad, 25, 19, 8, 8),
		pipad.Top:    pi.SpriteFrom(gamepad, 19, 13, 6, 8),
		pipad.Bottom: pi.SpriteFrom(gamepad, 19, 25, 6, 8),
	}

	pi.Draw = func() {
		pi.Cls()
		pi.Blit(gamepad, 0, 0)

		for button, sprite := range buttonSprites {
			if pipad.Duration(button) > 0 { // duration is > 0 when button is pressed
				pi.Spr(sprite, sprite.X, sprite.Y+1) // draw pressed button
			}
		}

		picofont.Print("PRESS BTN ON GAMEPAD", 3, 50)
	}

	piebiten.Run()
}
