// Example showing practical use of palette swapping. Drawing same sprite with different palette
// generates tens of different sprites.
package main

import (
	"embed"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/ebitengine"
)

const (
	eyes  = 1 // color number of eyes in sprite-sheet
	skin  = 7
	mouth = 8
	hair  = 5
)

var (
	//go:embed sprite-sheet.png
	resources embed.FS

	eyeColors   = [...]byte{0, 1, 3, 4, 12}
	skinColors  = [...]byte{7, 5, 15}
	hairColors  = [...]byte{0, 4, 5, 6, 7, 9, 10}
	mouthColors = [...]byte{2, 8}
)

func main() {
	pi.Load(resources)
	pi.Draw = draw
	ebitengine.MustRun()
}

func draw() {
	pi.Cls()
	x, y := 0, 0

	for _, eyeColor := range eyeColors {
		pi.Pal[eyes] = eyeColor

		for _, skinColor := range skinColors {
			pi.Pal[skin] = skinColor

			for _, hairColor := range hairColors {
				pi.Pal[hair] = hairColor

				for _, mouthColor := range mouthColors {
					pi.Pal[mouth] = mouthColor
					// draw the sprite with swapped colors:
					pi.Spr(0, x, y)

					x += 8
					if x >= 128 {
						// go to next line
						x = 0
						y += 8
					}
				}
			}
		}
	}
}
