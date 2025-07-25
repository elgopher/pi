// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	_ "embed"
	"github.com/elgopher/pi/pitest"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi"
)

var (
	//go:embed "internal/test/decode/indexed.png"
	indexedPNG []byte
	//go:embed "internal/test/decode/rgb.png"
	rgbPNG []byte
	//go:embed "internal/test/decode/indexed-brighter.png"
	brighterIndexedPNG []byte
	//go:embed "internal/test/decode/rgb-brighter.png"
	rgbBrighterPNG []byte
)

func TestDecodeCanvasOrErr(t *testing.T) {
	tests := map[string][]byte{
		"indexed png when palette is the same": indexedPNG,
		"RGB png":                              rgbPNG,
		"indexed png when palette is slightly brighter": brighterIndexedPNG,
		"RGB png when palette is slightly brighter":     rgbBrighterPNG,
	}

	for testName, png := range tests {
		t.Run(testName, func(t *testing.T) {
			pi.Palette = pi.DecodePalette(indexedPNG)
			// when
			canvas, err := pi.DecodeCanvasOrErr(png)
			// then
			require.NoError(t, err)
			expected := pi.NewCanvas(4, 4)
			expected.SetAll(
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			)
			pitest.AssertSurfaceEqual(t, expected, canvas)
		})
	}
}

func TestSet(t *testing.T) {
	t.Run("should be noop when outside surface", func(t *testing.T) {
		width := 2
		height := 3

		tests := map[string]struct {
			x, y int
		}{
			"left":   {x: -1, y: 0},
			"right":  {x: width, y: 0},
			"top":    {x: 0, y: -1},
			"bottom": {x: 0, y: height},
		}
		for testName, test := range tests {
			t.Run(testName, func(t *testing.T) {
				original := pi.NewSurface[int](width, height)
				s := original.Clone()
				// when
				s.Set(test.x, test.y, 1)
				// then
				assert.Equal(t, original, s)
			})
		}
	})

	t.Run("should set value", func(t *testing.T) {
		surface := pi.NewSurface[int](3, 3)
		surface.Set(1, 2, 1)
		expected := pi.NewSurface[int](3, 3)
		expected.SetAll(0, 0, 0, 0, 0, 0, 0, 1, 0)
		assert.Equal(t, expected, surface)
	})
}

func TestSurface_SetArea(t *testing.T) {
	// temporary test
	surface := pi.NewSurface[int](2, 2)
	area := pi.IntArea{X: -1, Y: 0, W: 2, H: 2}
	surface.SetArea(area, 1, 2, 3, 4)
	surface.SetArea(area, []int{1, 2, 3, 4}...)

	area2 := pi.IntArea{X: -2, Y: -2, W: 2, H: 2}
	surface.SetArea(area2, []int{1, 2, 3, 4}...)
}

func BenchmarkSurface_SetArea(b *testing.B) {
	// temporary test
	surface := pi.NewSurface[int](1920, 1080)
	area := pi.IntArea{X: -1, Y: 0, W: 32, H: 32}
	slice := make([]int, area.Size())
	for i := 0; i < len(slice); i++ {
		slice[i] = rand.Int()
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		surface.SetArea(area, slice...)
	}
}

func TestSurface_SetSurface(t *testing.T) {
	// temporary test
	surface := pi.NewSurface[int](4, 4)
	src := pi.NewSurface[int](2, 2)

	surface.SetSurface(-1, 0, src)
	surface.SetSurface(0, -1, src)
	surface.SetSurface(3, 3, src)
	surface.SetSurface(2, 2, src)
	surface.SetSurface(4, 4, src)
}

func TestSurface_LinesIterator(t *testing.T) {
	// temporary test
	surface := pi.NewSurface[int](2, 2)
	for range surface.LinesIterator(pi.IntArea{W: 2, H: 2}) {
	}
}

func BenchmarkSurface_LinesIterator(b *testing.B) {
	// temporary test
	surface := pi.NewSurface[int](2, 2)
	area := pi.IntArea{W: 2, H: 2}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for coords, line := range surface.LinesIterator(area) {
			coords.X += line[0]
		}
	}
}

func TestSurface_SetMany(t *testing.T) {
	// temporary test
	canvas := pi.NewCanvas(4, 4)
	canvas.SetMany(1, 1,
		2, 3, 4, 5)
	canvas.SetMany(-1, 0,
		2, 3, 4, 5)
	canvas.SetMany(0, -1,
		2, 3, 4, 5)
	canvas.SetMany(3, 3,
		2, 3, 4, 5)
	canvas.SetMany(4, 4,
		2, 3, 4, 5)
}

func TestSurface_Clear(t *testing.T) {
	surface := pi.NewSurface[int](2, 3)
	surface.Clear(7)
	assert.Equal(t, []int{7, 7, 7, 7, 7, 7}, surface.Data())
}

func BenchmarkSurface_Clear(b *testing.B) {
	surface := pi.NewSurface[int](320, 180)
	for b.Loop() {
		surface.Clear(7)
	}
}

func TestSurface_Get(t *testing.T) {
	surface := pi.NewSurface[int](2, 3)
	surface.Clear(7)
	surface.Set(1, 1, 5)
	assert.Equal(t, 5, surface.Get(1, 1))
	assert.Equal(t, 7, surface.Get(0, 0))
	assert.Equal(t, 0, surface.Get(-1, 1))
	assert.Equal(t, 0, surface.Get(1, -1))
	assert.Equal(t, 0, surface.Get(2, 2))
	assert.Equal(t, 0, surface.Get(1, 3))
}

func BenchmarkSurface_Get(b *testing.B) {
	surface := pi.NewSurface[int](2, 3)
	surface.Clear(7)
	surface.Set(1, 1, 5)

	for b.Loop() {
		surface.Get(1, 1)
	}
}

func BenchmarkDrawCanvas(b *testing.B) {
	dst := pi.NewCanvas(320, 180)
	pi.SetDrawTarget(dst)

	src := pi.NewCanvas(32, 32)
	src.Clear(7)

	for b.Loop() {
		pi.DrawCanvas(src, 130, 130)
	}
}

func TestSurface_GetLine(t *testing.T) {
	surface := pi.NewSurface[int](3, 2)
	surface.SetAll(1, 2, 3, 4, 5, 6)

	assert.Equal(t, []int{1, 2, 3}, surface.GetLine(0))
	assert.Equal(t, []int{4, 5, 6}, surface.GetLine(1))
	assert.Nil(t, surface.GetLine(-1))
	assert.Nil(t, surface.GetLine(2))
}
