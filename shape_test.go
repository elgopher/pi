// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	_ "embed"
	"testing"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/pisnap"
	"github.com/elgopher/pi/pitest"
	"github.com/stretchr/testify/require"
)

//go:embed "internal/test/shapes.png"
var shapes []byte

func TestShapes(t *testing.T) {
	t.Run("should draw shapes in the same way as in shapes.png", func(t *testing.T) {
		pi.ResetPaletteMapping()
		pi.ResetColorTables()
		pi.Palette = pi.DecodePalette(shapes)
		shapesSheet := pi.DecodeCanvas(shapes)

		tests := map[string]struct {
			draw     func()
			col, row int
		}{
			"horizontal line": {
				draw: line(4, 15, 27, 15),
				col:  0, row: 0,
			},
			"horizontal line reversed": {
				draw: line(27, 15, 4, 15),
				col:  0, row: 0,
			},
			"vertical line": {
				draw: line(47-32, 2, 47-32, 27),
				col:  1, row: 0,
			},
			"diagonal line 1": {
				draw: line(67-64, 2, 92-64, 27),
				col:  2, row: 0,
			},
			"diagonal line 2": {
				draw: line(100-96, 29, 122-96, 7),
				col:  3, row: 0,
			},
			"diagonal line 3": {
				draw: line(136-128, 4, 147-128, 27),
				col:  4, row: 0,
			},
			"diagonal line 4": {
				draw: line(165-160, 21, 187-160, 16),
				col:  5, row: 0,
			},
			"line dot": {
				draw: line(208-192, 16, 208-192, 16),
				col:  6, row: 0,
			},
			"rect dot": {
				draw: rect(16, 47-32, 16, 47-32),
				col:  0, row: 1,
			},
			"rect horizontal line": {
				draw: rect(40-32, 47-32, 54-32, 47-32),
				col:  1, row: 1,
			},
			"rect vertical line": {
				draw: rect(78-64, 38-32, 78-64, 56-32),
				col:  2, row: 1,
			},
			"square": {
				draw: rect(102-96, 39-32, 119-96, 56-32),
				col:  3, row: 1,
			},
			"square reversed": {
				draw: rect(119-96, 56-32, 102-96, 39-32),
				col:  3, row: 1,
			},
			"rectangle wide": {
				draw: rect(130-128, 45-32, 155-128, 51-32),
				col:  4, row: 1,
			},
			"rectangle toll": {
				draw: rect(170-160, 38-32, 179-160, 58-32),
				col:  5, row: 1,
			},
			"rectfill dot": {
				draw: rectfill(16, 80-64, 16, 80-64),
				col:  0, row: 2,
			},
			"rectfill horizontal line": {
				draw: rectfill(40-32, 80-64, 54-32, 80-64),
				col:  1, row: 2,
			},
			"rectfill vertical line": {
				draw: rectfill(78-64, 71-64, 78-64, 89-64),
				col:  2, row: 2,
			},
			"rectfill square": {
				draw: rectfill(102-96, 72-64, 119-96, 89-64),
				col:  3, row: 2,
			},
			"rectfill square reversed": {
				draw: rectfill(119-96, 89-64, 102-96, 72-64),
				col:  3, row: 2,
			},
			"rectfill wide": {
				draw: rectfill(130-128, 78-64, 155-128, 84-64),
				col:  4, row: 2,
			},
			"rectfill toll": {
				draw: rectfill(170-160, 71-64, 179-160, 91-64),
				col:  5, row: 2,
			},
			"circ": {
				draw: circ(19, 112-96, 9),
				col:  0, row: 3,
			},
			"tiny circ": {
				draw: circ(49-32, 112-96, 1),
				col:  1, row: 3,
			},
			"small circ": {
				draw: circ(79-64, 110-96, 2),
				col:  2, row: 3,
			},
			"big circ": {
				draw: circ(110-96, 113-96, 10),
				col:  3, row: 3,
			},
			"ugly circ": {
				draw: circ(144-128, 111-96, 4),
				col:  4, row: 3,
			},
			"dot circ": {
				draw: circ(175-160, 110-96, 0),
				col:  5, row: 3,
			},
			"circfill": {
				draw: circfill(19, 144-128, 9),
				col:  0, row: 4,
			},
			"tiny circfill": {
				draw: circfill(49-32, 144-128, 1),
				col:  1, row: 4,
			},
			"small circfill": {
				draw: circfill(79-64, 142-128, 2),
				col:  2, row: 4,
			},
			"big circfill": {
				draw: circfill(110-96, 145-128, 10),
				col:  3, row: 4,
			},
			"ugly circfill": {
				draw: circfill(144-128, 143-128, 4), // TODO this circfill is not as ugly as circ :)
				col:  4, row: 4,
			},
			"dot circfill": {
				draw: circfill(175-160, 142-128, 0),
				col:  5, row: 4,
			},
		}

		for testName, shape := range tests {
			t.Run(testName, func(t *testing.T) {
				pi.SetScreenSize(32, 32)
				pi.Cls()
				pi.SetColor(1)
				// when
				shape.draw()
				// then
				shapeArea := pi.IntArea{X: shape.col * 32, Y: shape.row * 32, W: 32, H: 32}
				expected := shapesSheet.CloneArea(shapeArea)
				equal := pitest.AssertSurfaceEqual(t, expected, pi.Screen())
				if !equal {
					f, err := pisnap.CaptureOrErr()
					require.NoError(t, err)
					t.Log("INVALID SNAPSHOT STORED TO", f)
				}
			})
		}
	})
}

func line(x0, y0, x1, y1 int) func() {
	return func() {
		pi.Line(x0, y0, x1, y1)
	}
}

func rect(x0, y0, x1, y1 int) func() {
	return func() {
		pi.Rect(x0, y0, x1, y1)
	}
}

func rectfill(x0, y0, x1, y1 int) func() {
	return func() {
		pi.RectFill(x0, y0, x1, y1)
	}
}

func circ(cx, cy, r int) func() {
	return func() {
		pi.Circ(cx, cy, r)
	}
}

func circfill(cx, cy, r int) func() {
	return func() {
		pi.CircFill(cx, cy, r)
	}
}
