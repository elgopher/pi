// Example showing how to print text to screen.
package main

import (
	"github.com/elgopher/pi"
)

func main() {
	pi.Draw = func() {
		pi.CursorReset()   // set cursor to 0,0
		pi.Color(9)        // change to yellow
		pi.Cursor(50, 58)  // set cursor position
		pi.Print("HELLO,") // print text and go to next line
		pi.Color(12)
		pi.Print("GOPHER!")
	}
	pi.RunOrPanic()
}
