// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/elgopher/pi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrint(t *testing.T) {
	pi.ScreenWidth = 16
	pi.ScreenHeight = 24

	t.Run("should print chars using color on the top-left corner", func(t *testing.T) {
		chars := []string{`!`, `A`, `b`, `AB`, `ABCD`}
		for _, char := range chars {
			t.Run(char, func(t *testing.T) {
				pi.BootOrPanic()
				// when
				pi.Color(7)
				pi.Print(char)
				// then
				assertScreenEqual(t, "internal/testimage/print/"+char+".png")
			})
		}
	})

	t.Run("should print question mark for characters > 255", func(t *testing.T) {
		pi.BootOrPanic()
		pi.Color(7)
		pi.Print("\u0100")
		assertScreenEqual(t, "internal/testimage/print/unknown.png")
	})

	t.Run("should print special character", func(t *testing.T) {
		pi.BootOrPanic()
		pi.Color(7)
		pi.Print("\u0080")
		assertScreenEqual(t, "internal/testimage/print/special.png")
	})

	t.Run("should print 2 special characters", func(t *testing.T) {
		pi.BootOrPanic()
		pi.Color(7)
		pi.Print("\u0080\u0081")
		assertScreenEqual(t, "internal/testimage/print/special-2chars.png")
	})

	t.Run("should go to next line", func(t *testing.T) {
		pi.BootOrPanic()
		pi.Color(7)
		pi.Print("0L")
		pi.Print("1L")
		assertScreenEqual(t, "internal/testimage/print/two-lines.png")
	})

	t.Run("should print at cursor position", func(t *testing.T) {
		tests := map[string]struct {
			x, y int
			file string
		}{
			"1,2":   {x: 1, y: 2, file: "internal/testimage/print/two-lines-at-1.2.png"},
			"-1,-2": {x: -1, y: -2, file: "internal/testimage/print/two-lines-at--1.-2.png"},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pi.BootOrPanic()
				pi.Color(7)
				pi.Cursor(test.x, test.y)
				pi.Print("0L")
				pi.Print("1L")
				assertScreenEqual(t, test.file)
			})
		}
	})

	t.Run("should print moved by camera position", func(t *testing.T) {
		pi.BootOrPanic()
		pi.Color(7)
		pi.Camera(-1, -2)
		pi.Print("0L")
		pi.Print("1L")
		assertScreenEqual(t, "internal/testimage/print/two-lines-at-1.2.png")
	})

	t.Run("should clip text", func(t *testing.T) {
		tests := map[string]struct {
			x, y, w, h       int
			cursorX, cursorY int
			cameraX, cameraY int
			file             string
		}{
			"clip left":                {x: 1, y: 0, w: 16, h: 16, file: "clip-left.png"},
			"clip right":               {x: 0, y: 0, w: 6, h: 16, file: "clip-right.png"},
			"clip top":                 {x: 0, y: 1, w: 16, h: 16, file: "clip-top.png"},
			"clip bottom":              {x: 0, y: 0, w: 16, h: 4, file: "clip-bottom.png"},
			"clip left, cursorX set":   {x: 2, y: 0, w: 16, h: 16, cursorX: 1, file: "clip-left-cursorx.png"},
			"clip right, cursorX set":  {x: 0, y: 0, w: 7, h: 16, cursorX: 1, file: "clip-right-cursorx.png"},
			"clip top, cursorY set":    {x: 0, y: 2, w: 16, h: 16, cursorY: 1, file: "clip-top-cursory.png"},
			"clip bottom, cursorY set": {x: 0, y: 0, w: 16, h: 5, cursorY: 1, file: "clip-bottom-cursory.png"},
			"camerax, clip left":       {x: 2, y: 0, w: 16, h: 16, cameraX: -1, file: "clip-left-cursorx.png"},
			"camerax, clip right":      {x: 0, y: 0, w: 7, h: 16, cameraX: -1, file: "clip-right-cursorx.png"},
			"cameray, clip top":        {x: 0, y: 2, w: 16, h: 16, cameraY: -1, file: "clip-top-cursory.png"},
			"cameray, clip bottom":     {x: 0, y: 0, w: 16, h: 5, cameraY: -1, file: "clip-bottom-cursory.png"},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pi.BootOrPanic()
				pi.Camera(test.cameraX, test.cameraY)
				pi.Clip(test.x, test.y, test.w, test.h)
				pi.Color(7)
				pi.Cursor(test.cursorX, test.cursorY)
				pi.Print("\u0080")
				assertScreenEqual(t, "internal/testimage/print/"+test.file)
			})
		}
	})

	t.Run("should reset cursor", func(t *testing.T) {
		tests := map[string]func(){
			"Cls()":         func() { pi.Cls() },
			"CursorReset()": func() { pi.CursorReset() },
			"Cursor(0,0)":   func() { pi.Cursor(0, 0) },
		}
		for name, function := range tests {
			t.Run(name, func(t *testing.T) {
				pi.BootOrPanic()
				pi.Color(7)
				pi.Cursor(20, 20)
				// when
				function()
				pi.Print("\u0080")
				// then
				assertScreenEqual(t, "internal/testimage/print/special.png")
			})
		}
	})

	t.Run("should scroll screen when there is no space", func(t *testing.T) {
		const charHeight = 6

		tests := map[string]struct {
			cursorY      int
			expectedFile string
		}{
			"no scrolling": {
				cursorY:      pi.ScreenHeight - charHeight*2,
				expectedFile: "internal/testimage/print/no-scrolling.png",
			},
			"print, then scroll 1": {
				cursorY:      pi.ScreenHeight - charHeight,
				expectedFile: "internal/testimage/print/print-then-scroll.png",
			},
			"print, then scroll 2": {
				cursorY:      pi.ScreenHeight - charHeight - 1,
				expectedFile: "internal/testimage/print/print-then-scroll2.png",
			},
			"scroll, print and scroll": {
				cursorY:      pi.ScreenHeight - charHeight + 1,
				expectedFile: "internal/testimage/print/scroll-print-scroll.png",
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pi.BootOrPanic()
				pi.ClsCol(3)
				pi.Color(7)
				pi.Cursor(0, test.cursorY)
				// when
				pi.Print("\u0080")
				// then
				assertScreenEqual(t, test.expectedFile)
			})
		}
	})

	t.Run("should scroll screen (clipping)", func(t *testing.T) {
		const charHeight = 6

		tests := map[string]struct {
			cursorY      int
			clipY, clipH int
			clipX, clipW int
			expectedFile string
		}{
			"clip h, print, then scroll": {
				clipH:        pi.ScreenHeight - 1,
				clipW:        pi.ScreenWidth,
				cursorY:      pi.ScreenHeight - charHeight,
				expectedFile: "internal/testimage/print/clip-print-then-scroll.png",
			},
			"clip yh, print, then scroll": {
				clipY:        1,
				clipH:        pi.ScreenHeight - 2,
				clipW:        pi.ScreenWidth,
				cursorY:      pi.ScreenHeight - charHeight,
				expectedFile: "internal/testimage/print/clipyh-print-then-scroll.png",
			},
			"clip all, print, then scroll": {
				clipX:        1,
				clipY:        1,
				clipW:        5,
				clipH:        pi.ScreenHeight - 2,
				cursorY:      pi.ScreenHeight - charHeight,
				expectedFile: "internal/testimage/print/clipall-print-then-scroll.png",
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pi.BootOrPanic()
				pi.ScreenData = decodePNG(t, "internal/testimage/print/multicolor.png").Pixels
				pi.Clip(test.clipX, test.clipY, test.clipW, test.clipH)
				pi.Color(7)
				pi.Cursor(0, test.cursorY)
				// when
				pi.Print("\u0080")
				// then
				assertScreenEqual(t, test.expectedFile)
			})
		}
	})

	t.Run("should return the right-most x position that occurred while printing", func(t *testing.T) {
		tests := map[string]struct {
			text      string
			expectedX int
		}{
			"empty":              {text: "", expectedX: 0},
			"normal char":        {text: "a", expectedX: 4},
			"wide char":          {text: "\u0080", expectedX: 8},
			"two normal chars":   {text: "ab", expectedX: 8},
			"normal + wide char": {text: "\u0080b", expectedX: 12},
			"100 chars":          {text: strings.Repeat("a", 100), expectedX: 400},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pi.BootOrPanic()
				// when
				x := pi.Print(test.text)
				assert.Equal(t, test.expectedX, x)
			})
		}
	})
}

func assertScreenEqual(t *testing.T, file string) {
	expected := decodePNG(t, file).Pixels
	if !assert.Equal(t, expected, pi.ScreenData) {
		screenshot, err := pi.Snap()
		require.NoError(t, err)
		fmt.Println("Screenshot taken", screenshot)
	}
}

func TestCursor(t *testing.T) {
	t.Run("should return default cursor position", func(t *testing.T) {
		pi.BootOrPanic()
		x, y := pi.Cursor(1, 1)
		assert.Zero(t, x)
		assert.Zero(t, y)
	})

	t.Run("should return previous cursor position", func(t *testing.T) {
		pi.BootOrPanic()
		prevX, prevY := 1, 2
		pi.Cursor(prevX, prevY)
		// when
		x, y := pi.Cursor(3, 4)
		assert.Equal(t, prevX, x)
		assert.Equal(t, prevY, y)
	})
}

func TestCursorReset(t *testing.T) {
	t.Run("should return default cursor position", func(t *testing.T) {
		pi.BootOrPanic()
		x, y := pi.CursorReset()
		assert.Zero(t, x)
		assert.Zero(t, y)
	})

	t.Run("should return previous cursor position", func(t *testing.T) {
		pi.BootOrPanic()
		prevX, prevY := 1, 2
		pi.Cursor(prevX, prevY)
		// when
		x, y := pi.CursorReset()
		assert.Equal(t, prevX, x)
		assert.Equal(t, prevY, y)
	})
}
