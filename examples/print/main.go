// Example showing how to print text to screen.
package main

import (
	"github.com/elgopher/pi"
)

func main() {
	pi.Draw = func() {
		pi.CursorReset()      // set cursor to 0,0
		pi.Cursor(50, 58)     // set cursor position
		pi.Print("HELLO,", 9) // print yellow text and go to next line
		pi.Print("GOPHER!", 12)
	}
	pi.RunOrPanic()
}
