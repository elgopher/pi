// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import "math"

// Rect draws the outline of a rectangle between (x0, y0) and (x1, y1), inclusive.
//
// It takes into account the camera position, clipping region,
// color tables, and masks.
func Rect(x0, y0, x1, y1 int) {
	// Optimize - run vertical and horizontal line functions directly
	Line(x0, y0, x1, y0) // horizontal line top
	Line(x0, y1, x1, y1) // horizontal line bottom
	Line(x0, y0, x0, y1) // vertical line left
	Line(x1, y0, x1, y1) // vertical line right
}

// RectFill draws a filled rectangle between (x0, y0) and (x1, y1), inclusive.
//
// It takes into account the camera position, clipping region,
// color tables, and masks.
func RectFill(x0 int, y0 int, x1 int, y1 int) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}

	area := IntArea{
		X: x0,
		Y: y0,
		W: x1 - x0 + 1,
		H: y1 - y0 + 1,
	}

	area = area.MovedBy(-Camera.X, -Camera.Y)
	area, _, _ = area.ClippedBy(Clip())

	currentColor := GetColor() & ReadMask

	for _, line := range DrawTarget().LinesIterator(area) {
		for i := 0; i < len(line); i++ {
			target := line[i]
			line[i] = ColorTables[(currentColor|target)>>6][currentColor&(MaxColors-1)][target&(MaxColors-1)]
		}
	}
}

// Line draws a line on the screen between (x0, y0) and (x1, y1), inclusive.
//
// It takes into account the camera position, clipping region,
// color tables, and masks.
func Line(x0, y0, x1, y1 int) {
	draw := drawColor & ReadMask
	// Optimize - add vertical and horizontal line functions

	// Bresenham algorithm: https://www.youtube.com/watch?v=IDFB5CDpLDE
	run := float64(x1 - x0)

	rise := float64(y1 - y0)

	slope := rise / run

	adjust := 1
	if slope < 0 {
		adjust = -1
	}

	offset := 0.0    // performance could be better if offset was an integer instead
	threshold := 0.5 // performance could be better if threshold was an integer instead

	if slope >= -1 && slope <= 1 {
		delta := math.Abs(slope)
		y := y0
		if x1 < x0 {
			x0, x1 = x1, x0
			y = y1
		}

		for x := x0; x <= x1; x++ {
			setPixelWithColor(x, y, draw)

			offset += delta
			if offset >= threshold {
				y += adjust
				threshold += 1
			}
		}
	} else {
		delta := math.Abs(run / rise)
		x := x0
		if y0 > y1 {
			y0, y1 = y1, y0
			x = x1
		}

		for y := y0; y <= y1; y++ {
			setPixelWithColor(x, y, draw)

			offset += delta
			if offset >= threshold {
				x += adjust
				threshold += 1
			}
		}
	}
}

// Circ draws the outline of a circle with center at (cx, cy) and radius r.
//
// It takes into account the camera position, clipping region,
// color tables, and masks.
func Circ(cx, cy, r int) {
	draw := drawColor & ReadMask

	x := 0
	y := r
	d := 3 - 2*r

	for x <= y {
		if x == 0 {
			setPixelWithColor(cx+y, cy, draw)
			setPixelWithColor(cx-y, cy, draw)
			setPixelWithColor(cx, cy+y, draw)
			setPixelWithColor(cx, cy-y, draw)
		} else {
			setPixelWithColor(cx+x, cy+y, draw)
			setPixelWithColor(cx-x, cy+y, draw)
			setPixelWithColor(cx+x, cy-y, draw)
			setPixelWithColor(cx-x, cy-y, draw)

			if x != y {
				setPixelWithColor(cx+y, cy+x, draw)
				setPixelWithColor(cx-y, cy+x, draw)
				setPixelWithColor(cx+y, cy-x, draw)
				setPixelWithColor(cx-y, cy-x, draw)
			}
		}

		if d <= 0 {
			d += 4*x + 6
		} else {
			d += 4*(x-y) + 10
			y--
		}
		x++
	}
}

func horizontalLine(x0, x1, y int) {
	RectFill(x0, y, x1, y)
}

// CircFill draws a filled circle with center at (centerX, centerY) and the given radius.
//
// It takes into account the camera position, clipping region,
// color tables, and masks.
func CircFill(centerX, centerY, radius int) {
	// Algorithm designed by https://stackoverflow.com/users/3797048/colinday
	//
	// Details: https://stackoverflow.com/questions/10878209/midpoint-circle-algorithm-for-filled-circles#answer-24527943
	x := radius
	y := 0
	radiusError := 1 - x

	// iterate to the circle diagonal
	for x >= y {
		// use symmetry to draw the two horizontal lines at this Y with a special case to draw
		// only one line at the centerY where y == 0
		startX := -x + centerX
		endX := x + centerX
		horizontalLine(startX, endX, y+centerY)
		if y != 0 {
			horizontalLine(startX, endX, -y+centerY)
		}

		// move Y one line
		y++

		// calculate or maintain new x
		if radiusError < 0 {
			radiusError += 2*y + 1
		} else {
			// we're about to move x over one, this means we completed a column of X values, use
			// symmetry to draw those complete columns as horizontal lines at the top and bottom of the circle
			// beyond the diagonal of the main loop
			if x >= y {
				startX = -y + 1 + centerX
				endX = y - 1 + centerX
				horizontalLine(startX, endX, x+centerY)
				horizontalLine(startX, endX, -x+centerY)
			}
			x--
			radiusError += 2 * (y - x + 1)
		}
	}
}
