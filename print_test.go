// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	_ "embed"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi"
)

func TestSetCustomFontWidth(t *testing.T) {
	tests := map[int]int{
		-1: 0,
		0:  0,
		1:  1,
		8:  8,
		9:  8,
	}

	for width, expectedWidth := range tests {
		pi.SetCustomFontWidth(width)
		assert.Equal(t, expectedWidth, pi.CustomFont().Width)
	}

	for width, expectedWidth := range tests {
		pi.SetCustomFontSpecialWidth(width)
		assert.Equal(t, expectedWidth, pi.CustomFont().SpecialWidth)
	}
}

func TestSetCustomFontHeight(t *testing.T) {
	tests := map[int]int{
		-1: 0,
		0:  0,
		1:  1,
	}
	for height, expectedHeight := range tests {
		pi.SetCustomFontHeight(height)
		assert.Equal(t, expectedHeight, pi.CustomFont().Height)
	}
}

func TestPrint(t *testing.T) {
	reset := func() {
		pi.Reset()
		pi.SetScreenSize(16, 24)
	}

	const color = 7

	t.Run("should print chars using color on the top-left corner", func(t *testing.T) {
		chars := []string{`!`, `A`, `b`, `AB`, `ABCD`}
		for _, char := range chars {
			t.Run(char, func(t *testing.T) {
				reset()
				// when
				pi.Print(char, 0, 0, color)
				// then
				assertScreenEqual(t, "internal/testimage/print/"+char+".png")
			})
		}
	})

	t.Run("should print question mark for characters > 255", func(t *testing.T) {
		reset()
		pi.Print("\u0100", 0, 0, color)
		assertScreenEqual(t, "internal/testimage/print/unknown.png")
	})

	t.Run("should print special character", func(t *testing.T) {
		reset()
		pi.Print("\u0080", 0, 0, color)
		assertScreenEqual(t, "internal/testimage/print/special.png")
	})

	t.Run("should print 2 special characters", func(t *testing.T) {
		reset()
		pi.Print("\u0080\u0081", 0, 0, color)
		assertScreenEqual(t, "internal/testimage/print/special-2chars.png")
	})

	t.Run("should go to next line", func(t *testing.T) {
		reset()
		pi.Print("0L\n1L", 0, 0, color)
		assertScreenEqual(t, "internal/testimage/print/two-lines.png")
	})

	t.Run("should print at position", func(t *testing.T) {
		tests := map[string]struct {
			x, y int
			file string
		}{
			"1,2":   {x: 1, y: 2, file: "internal/testimage/print/two-lines-at-1.2.png"},
			"-1,-2": {x: -1, y: -2, file: "internal/testimage/print/two-lines-at--1.-2.png"},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				reset()
				pi.Print("0L\n1L", test.x, test.y, color)
				assertScreenEqual(t, test.file)
			})
		}
	})

	t.Run("should print moved by camera position", func(t *testing.T) {
		reset()
		pi.Camera(-1, -2)
		pi.Print("0L\n1L", 0, 0, color)
		assertScreenEqual(t, "internal/testimage/print/two-lines-at-1.2.png")
	})

	t.Run("should clip text", func(t *testing.T) {
		tests := map[string]struct {
			x, y, w, h       int
			posX, posY       int
			cameraX, cameraY int
			file             string
		}{
			"clip left":             {x: 1, y: 0, w: 16, h: 16, file: "clip-left.png"},
			"clip right":            {x: 0, y: 0, w: 6, h: 16, file: "clip-right.png"},
			"clip top":              {x: 0, y: 1, w: 16, h: 16, file: "clip-top.png"},
			"clip bottom":           {x: 0, y: 0, w: 16, h: 4, file: "clip-bottom.png"},
			"clip left, posX set":   {x: 2, y: 0, w: 16, h: 16, posX: 1, file: "clip-left-posx.png"},
			"clip right, posX set":  {x: 0, y: 0, w: 7, h: 16, posX: 1, file: "clip-right-posx.png"},
			"clip top, posY set":    {x: 0, y: 2, w: 16, h: 16, posY: 1, file: "clip-top-posy.png"},
			"clip bottom, posY set": {x: 0, y: 0, w: 16, h: 5, posY: 1, file: "clip-bottom-posy.png"},
			"camerax, clip left":    {x: 2, y: 0, w: 16, h: 16, cameraX: -1, file: "clip-left-posx.png"},
			"camerax, clip right":   {x: 0, y: 0, w: 7, h: 16, cameraX: -1, file: "clip-right-posx.png"},
			"cameray, clip top":     {x: 0, y: 2, w: 16, h: 16, cameraY: -1, file: "clip-top-posy.png"},
			"cameray, clip bottom":  {x: 0, y: 0, w: 16, h: 5, cameraY: -1, file: "clip-bottom-posy.png"},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				reset()
				pi.Camera(test.cameraX, test.cameraY)
				pi.Clip(test.x, test.y, test.w, test.h)
				pi.Print("\u0080", test.posX, test.posY, color)
				assertScreenEqual(t, "internal/testimage/print/"+test.file)
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
				reset()
				// when
				x := pi.Print(test.text, 0, 0, color)
				assert.Equal(t, test.expectedX, x)
			})
		}
	})
}

func TestFont_String(t *testing.T) {
	t.Run("should convert font with small data to string", func(t *testing.T) {
		font := pi.Font{
			Data:         make([]byte, 2), // invalid data, but still possible to initialize by hand
			Width:        1,
			SpecialWidth: 2,
			Height:       3,
		}
		actual := font.String()
		assert.Equal(t, "{width: 1, specialWidth: 2, height: 3, data: [0 0]}", actual)
	})

	t.Run("should convert font to string", func(t *testing.T) {
		font := pi.Font{
			Data:         make([]byte, 2048),
			Width:        1,
			SpecialWidth: 2,
			Height:       3,
		}
		actual := font.String()
		fmt.Println(len(actual))
		assert.True(t, len(actual) < 1500)
	})
}
