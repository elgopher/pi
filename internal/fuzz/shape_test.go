// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package fuzz_test

import (
	"testing"

	"github.com/elgopher/pi"
)

func FuzzRect(f *testing.F) {
	pi.ScreenWidth = 16
	pi.ScreenHeight = 16
	pi.BootOrPanic()
	f.Add(0, 0, 0, 0)
	f.Fuzz(func(t *testing.T, x0, y0, x1, y1 int) {
		pi.Rect(x0, y0, x1, y1)
	})
}

func FuzzRectFill(f *testing.F) {
	pi.ScreenWidth = 16
	pi.ScreenHeight = 16
	pi.BootOrPanic()
	f.Add(0, 0, 0, 0)
	f.Fuzz(func(t *testing.T, x0, y0, x1, y1 int) {
		pi.RectFill(x0, y0, x1, y1)
	})
}
