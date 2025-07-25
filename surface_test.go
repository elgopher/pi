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

func TestSurface_Set(t *testing.T) {
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
	t.Run("inside surface", func(t *testing.T) {
		surface := pi.NewSurface[rune](4, 3)
		area := pi.IntArea{X: 1, Y: 1, W: 2, H: 2}
		// when
		surface.SetArea(area,
			'a', 'b',
			'c', 'd',
		)
		// then
		expected := pi.NewSurface[rune](4, 3)
		expected.SetMany(1, 1, 'a', 'b')
		expected.SetMany(1, 2, 'c', 'd')

		pitest.AssertSurfaceEqual(t, expected, surface)
	})

	t.Run("clipped area", func(t *testing.T) {
		surface := pi.NewSurface[rune](2, 2)
		area := pi.IntArea{X: -1, Y: -1, W: 3, H: 3}
		// when
		surface.SetArea(area,
			'a', 'b', 'c',
			'd', 'e', 'f',
			'g', 'h', 'i',
		)
		// then
		expected := pi.NewSurface[rune](2, 2)
		expected.SetMany(0, 0, 'e', 'f')
		expected.SetMany(0, 1, 'h', 'i')

		pitest.AssertSurfaceEqual(t, expected, surface)
	})

	t.Run("area outside surface", func(t *testing.T) {
		surface := pi.NewSurface[rune](2, 2)
		original := surface.Clone()
		area := pi.IntArea{X: 3, Y: 3, W: 2, H: 2}
		// when
		surface.SetArea(area,
			'a', 'b',
			'c', 'd')
		// then
		pitest.AssertSurfaceEqual(t, original, surface)
	})

	t.Run("panic on too few values", func(t *testing.T) {
		surface := pi.NewSurface[rune](2, 2)
		area := pi.IntArea{W: 2, H: 2}

		require.Panics(t, func() {
			surface.SetArea(area, 'a', 'b', 'c')
		})
	})
}

func BenchmarkSurface_SetArea(b *testing.B) {
	surface := pi.NewSurface[int](320, 180)
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
	src := pi.NewSurface[rune](2, 2)
	src.SetAll(
		'a', 'b',
		'c', 'd',
	)

	t.Run("inside destination", func(t *testing.T) {
		dst := pi.NewSurface[rune](4, 4)
		// when
		dst.SetSurface(1, 1, src)
		// then
		expected := pi.NewSurface[rune](4, 4)
		expected.SetMany(1, 1, 'a', 'b')
		expected.SetMany(1, 2, 'c', 'd')
		pitest.AssertSurfaceEqual(t, expected, dst)
	})

	t.Run("partially outside top-left", func(t *testing.T) {
		dst := pi.NewSurface[rune](4, 4)
		// when
		dst.SetSurface(-1, -1, src)
		// then
		expected := pi.NewSurface[rune](4, 4)
		expected.Set(0, 0, 'd')
		pitest.AssertSurfaceEqual(t, expected, dst)
	})

	t.Run("partially outside bottom-right", func(t *testing.T) {
		dst := pi.NewSurface[rune](4, 4)
		// when
		dst.SetSurface(3, 3, src)
		// then
		expected := pi.NewSurface[rune](4, 4)
		expected.Set(3, 3, 'a')
		pitest.AssertSurfaceEqual(t, expected, dst)
	})

	t.Run("completely outside", func(t *testing.T) {
		dst := pi.NewSurface[rune](4, 4)
		// when
		dst.SetSurface(4, 4, src)
		// then
		expected := pi.NewSurface[rune](4, 4)
		pitest.AssertSurfaceEqual(t, expected, dst)
	})
}

func TestSurface_LinesIterator(t *testing.T) {
	surface := pi.NewSurface[rune](4, 3)
	surface.SetAll(
		'a', 'b', 'c', 'd',
		'e', 'f', 'g', 'h',
		'i', 'j', 'k', 'l',
	)

	t.Run("iterate over area", func(t *testing.T) {
		area := pi.IntArea{X: 1, Y: 1, W: 2, H: 2}
		expectedPos := []pi.Position{
			{1, 1},
			{1, 2},
		}
		expectedLines := [][]rune{
			{'f', 'g'},
			{'j', 'k'},
		}

		var i int
		for pos, line := range surface.LinesIterator(area) {
			require.Less(t, i, len(expectedPos))
			assert.Equal(t, expectedPos[i], pos)
			assert.Equal(t, expectedLines[i], line)
			i++
		}
		assert.Equal(t, len(expectedPos), i)
	})

	t.Run("modifies underlying data", func(t *testing.T) {
		area := pi.IntArea{X: 1, Y: 1, W: 2, H: 1}
		for _, line := range surface.LinesIterator(area) {
			line[0] = 'z'
		}
		assert.Equal(t, 'z', surface.Get(1, 1))
	})

	t.Run("panic on area outside surface", func(t *testing.T) {
		area := pi.IntArea{X: -1, Y: 0, W: 2, H: 1}
		require.Panics(t, func() {
			for range surface.LinesIterator(area) {
			}
		})
	})
}

func BenchmarkSurface_LinesIterator(b *testing.B) {
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
	t.Run("inside surface", func(t *testing.T) {
		surface := pi.NewSurface[rune](3, 2)
		// when
		surface.SetMany(1, 0, 'a', 'b', 'c', 'd')
		// then
		expected := pi.NewSurface[rune](3, 2)
		expected.Set(1, 0, 'a')
		expected.Set(2, 0, 'b')
		expected.Set(0, 1, 'c')
		expected.Set(1, 1, 'd')

		pitest.AssertSurfaceEqual(t, expected, surface)
	})

	t.Run("partially outside left", func(t *testing.T) {
		surface := pi.NewSurface[rune](3, 2)
		// when
		surface.SetMany(-1, 0, 'a', 'b', 'c')
		// then
		expected := pi.NewSurface[rune](3, 2)
		expected.Set(0, 0, 'b')
		expected.Set(1, 0, 'c')

		pitest.AssertSurfaceEqual(t, expected, surface)
	})

	t.Run("partially outside top", func(t *testing.T) {
		surface := pi.NewSurface[rune](2, 2)
		// when
		surface.SetMany(0, -1, 'a', 'b', 'c')
		// then
		expected := pi.NewSurface[rune](2, 2)
		expected.Set(0, 0, 'c')

		pitest.AssertSurfaceEqual(t, expected, surface)
	})

	t.Run("start far left with too few values", func(t *testing.T) {
		surface := pi.NewSurface[rune](2, 2)
		original := surface.Clone()
		// when
		surface.SetMany(-3, 0, 'a')
		// then
		pitest.AssertSurfaceEqual(t, original, surface)
	})

	t.Run("outside bottom-right", func(t *testing.T) {
		surface := pi.NewSurface[rune](2, 2)
		original := surface.Clone()
		// when
		surface.SetMany(2, 2, 'a', 'b')
		// then
		pitest.AssertSurfaceEqual(t, original, surface)
	})

	t.Run("truncate at surface end", func(t *testing.T) {
		surface := pi.NewSurface[rune](2, 2)
		// when
		surface.SetMany(1, 1, 'a', 'b')
		// then
		expected := pi.NewSurface[rune](2, 2)
		expected.Set(1, 1, 'a')

		pitest.AssertSurfaceEqual(t, expected, surface)
	})
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
