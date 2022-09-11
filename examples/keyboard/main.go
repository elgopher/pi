// Example showing how to use virtual keyboard. Useful for writing tools
// or games leveraging mouse + keyboard input.
package main

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/key"
)

func main() {
	pi.Draw = func() {
		pi.Cls()
		drawKeyboard(6, 34)
	}
	pi.MustRun()
}

func drawKeyboard(x, y int) {
	drawKey(x, y, 6, key.Esc)
	drawFunctionKeys(x, y)

	drawLineOfKeys(x, y+10, "`1234567890-=")
	drawKey(x+104, y+10, 10, key.Back)

	drawKey(x, y+20, 10, key.Tab)
	drawLineOfKeys(x+12, y+20, "QWERTYUIOP[]\\")

	drawKey(x, y+30, 14, key.Cap)
	drawLineOfKeys(x+16, y+30, "ASDFGHJKL;'")
	drawKey(x+104, y+30, 10, key.Enter)

	drawKey(x, y+40, 18, key.Shift)
	drawLineOfKeys(x+20, y+40, "ZXCVBNM,./")
	drawKey(x+100, y+40, 6, key.Up)

	drawKey(x, y+50, 18, key.Ctrl)
	drawKey(x+20, y+50, 10, key.Alt)
	drawKey(x+32, y+50, 58, key.Space)
	drawKey(x+92, y+50, 6, key.Left)
	drawKey(x+100, y+50, 6, key.Down)
	drawKey(x+108, y+50, 6, key.Right)
}

func drawFunctionKeys(x int, y int) {
	for i := 0; i < 12; i++ {
		offset := x + 8 + (i * 9)
		button := key.Button(int(key.F1) + i)
		drawKey(offset, y, 7, button)
	}
}

// draw line of printable keys
func drawLineOfKeys(x, y int, s string) {
	for _, char := range s {
		// printable char (such as 'a') can be directly converted to Button
		button := key.Button(char)
		drawKey(x, y, 6, button)
		x += 8
	}
}

func drawKey(x, y, w int, b key.Button) {
	// key.Btn returns true if button is pressed
	if key.Btn(b) {
		pi.RectFill(x, y, x+w, y+8, 12)
	} else {
		pi.Rect(x, y, x+w, y+8, 12)
	}

	str := b.String()
	if len(str) == 1 {
		pi.Cursor(x+2, y+2)
		pi.Print(str, 15)
	}

}
