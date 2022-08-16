// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package fuzz_test

import (
	"testing"

	"github.com/elgopher/pi"
)

func FuzzPrint(f *testing.F) {
	pi.ScreenWidth = 16
	pi.ScreenHeight = 24
	pi.BootOrPanic()
	f.Fuzz(func(t *testing.T, x, y int) {
		pi.Cursor(x, y)
		pi.Print("A", color)
	})
}
