// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build !js

package fuzz_test

import (
	"testing"

	"github.com/elgopher/pi"
)

func FuzzRect(f *testing.F) {
	pi.ScreenWidth = 16
	pi.ScreenHeight = 16
	pi.BootOrPanic()
	f.Fuzz(func(t *testing.T, x0, y0, x1, y1 int) {
		pi.Rect(x0, y0, x1, y1, color)
	})
}

func FuzzRectFill(f *testing.F) {
	pi.ScreenWidth = 16
	pi.ScreenHeight = 16
	pi.BootOrPanic()
	f.Fuzz(func(t *testing.T, x0, y0, x1, y1 int) {
		pi.RectFill(x0, y0, x1, y1, color)
	})
}

func FuzzLine(f *testing.F) {
	pi.ScreenWidth = 16
	pi.ScreenHeight = 16
	pi.BootOrPanic()
	f.Fuzz(func(t *testing.T, x0, y0, x1, y1 int) {
		pi.Line(x0, y0, x1, y1, color)
	})
}

func FuzzCirc(f *testing.F) {
	pi.ScreenWidth = 16
	pi.ScreenHeight = 16
	pi.BootOrPanic()
	f.Fuzz(func(t *testing.T, x, y, r int) {
		pi.Circ(x, y, r, color)
	})
}

func FuzzCircFill(f *testing.F) {
	pi.ScreenWidth = 16
	pi.ScreenHeight = 16
	pi.BootOrPanic()
	f.Fuzz(func(t *testing.T, x, y, r int) {
		pi.CircFill(x, y, r, color)
	})
}
