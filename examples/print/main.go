// Example showing how to print text to screen.
package main

import (
	"embed"

	"github.com/elgopher/pi"
)

//go:embed custom-font.png
var resources embed.FS

func main() {
	pi.Resources = resources
	pi.Draw = func() {
		pi.Print("HELLO,\nWORLD!", 50, 58, 9) // print two lines of yellow text using system font
	}
	pi.MustRun()
}
