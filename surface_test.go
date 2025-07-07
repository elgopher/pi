// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi"
)

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

func BenchmarkBlit(b *testing.B) {
	dst := pi.NewCanvas(320, 180)
	pi.SetDrawTarget(dst)

	src := pi.NewCanvas(32, 32)
	src.Clear(7)

	for b.Loop() {
		pi.Blit(src, 130, 130)
	}
}
