// Example showing how to use PixMap struct, which is used to store screen
// and sprite-sheet pixels.
package main

import (
	"embed"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/ebitengine"
)

//go:embed sprite-sheet.png
var resources embed.FS

func main() {
	pi.Load(resources)

	// copy from sprite-sheet to sprite-sheet:
	src := pi.SprSheet().WithClip(10, 0, 100, 100)
	src.Copy(pi.SprSheet(), 0, 0)

	// draw a filled rectangle directly to sprite-sheet:
	pi.SprSheet().RectFill(60, 30, 70, 40, 7)

	// merge from sprite-sheet to screen using custom merge function, which merges two lines
	src = pi.SprSheet().WithClip(-1, -1, 103, 70)
	src.Merge(pi.Scr(), -1, -1, func(dst, src []byte) {
		for x := 0; x < len(dst); x++ {
			dst[x] += pi.Pal[src[x]] + 1
		}
	})

	// update each line in a loop:
	src = pi.Scr().WithClip(10, 10, 16, 16)
	src.Foreach(func(x, y int, line []byte) {
		for i := 0; i < len(line); i++ {
			line[i] = byte(i)
		}
	})

	ebitengine.MustRun()
}
