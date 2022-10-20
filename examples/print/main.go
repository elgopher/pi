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
	pi.CustomFont.Width = 6 // set the width of all characters below 128 (ascii)
	pi.Draw = func() {
		pi.Print("HELLO,\nMY NAME IS", 45, 58, 9)   // print two lines of yellow text using system font
		pi.CustomFont.Print("PI\u0082", 45, 70, 12) // print blue text with special character using custom font
	}
	pi.MustRun()
}
