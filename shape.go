// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import "math"

// RectFill draws a filled rectangle between points x0,y0 and x1,y1 (inclusive).
//
// RectFill takes into consideration: current color, camera position,
// clipping region and draw palette.
func RectFill(x0, y0, x1, y1 int) {
	xmin, xmax := x0-camera.x, x1-camera.x
	if xmin > xmax {
		xmin, xmax = xmax, xmin
	}
	ymin, ymax := y0-camera.y, y1-camera.y
	if ymin > ymax {
		ymin, ymax = ymax, ymin
	}

	if xmin >= clippingRegion.x+clippingRegion.w {
		return
	}

	if xmax < clippingRegion.x {
		return
	}

	if ymin >= clippingRegion.y+clippingRegion.h {
		return
	}

	if ymax < clippingRegion.y {
		return
	}

	if xmin < clippingRegion.x {
		xmin = clippingRegion.x
	}

	if xmax >= clippingRegion.x+clippingRegion.w {
		xmax = clippingRegion.x + clippingRegion.w - 1
	}

	if ymin < clippingRegion.y {
		ymin = clippingRegion.y
	}

	if ymax >= clippingRegion.y+clippingRegion.h {
		ymax = clippingRegion.y + clippingRegion.h - 1
	}

	w := xmax - xmin + 1
	col := drawPalette[color]
	for i := 0; i < w; i++ {
		lineOfScreenWidth[i] = col
	}
	line := lineOfScreenWidth[:w]

	for y := ymin; y <= ymax; y++ {
		copy(ScreenData[y*scrWidth+xmin:], line)
	}
}

// Rect draws a rectangle between points x0,y0 and x1,y1 (inclusive).
//
// Rect takes into consideration: current color, camera position,
// clipping region and draw palette.
func Rect(x0, y0, x1, y1 int) {
	xmin, xmax := x0-camera.x, x1-camera.x
	if xmin > xmax {
		xmin, xmax = xmax, xmin
	}
	ymin, ymax := y0-camera.y, y1-camera.y
	if ymin > ymax {
		ymin, ymax = ymax, ymin
	}

	if xmin >= clippingRegion.x+clippingRegion.w {
		return
	}

	if xmax < clippingRegion.x {
		return
	}

	if ymin >= clippingRegion.y+clippingRegion.h {
		return
	}

	if ymax < clippingRegion.y {
		return
	}

	drawLeftLine := true
	drawRightLine := true

	if xmin < clippingRegion.x {
		xmin = clippingRegion.x
		drawLeftLine = false
	}

	if xmax >= clippingRegion.x+clippingRegion.w {
		xmax = clippingRegion.x + clippingRegion.w - 1
		drawRightLine = false
	}

	col := drawPalette[color]

	w := xmax - xmin + 1
	for i := 0; i < w; i++ {
		lineOfScreenWidth[i] = col
	}
	line := lineOfScreenWidth[:w]

	if ymin < clippingRegion.y {
		ymin = clippingRegion.y
	} else {
		copy(ScreenData[ymin*scrWidth+xmin:], line)
	}

	if ymax >= clippingRegion.y+clippingRegion.h {
		ymax = clippingRegion.y + clippingRegion.h - 1
	} else {
		copy(ScreenData[ymax*scrWidth+xmin:], line)
	}

	if drawLeftLine {
		for y := ymin; y <= ymax; y++ {
			ScreenData[y*scrWidth+xmin] = col
		}
	}

	if drawRightLine {
		for y := ymin; y <= ymax; y++ {
			ScreenData[y*scrWidth+xmax] = col
		}
	}
}

// Line draws a line between points x0,y0 and x1,y1 (inclusive).
//
// Line takes into account the current color, camera position,
// clipping region and draw palette.
func Line(x0, y0, x1, y1 int) {
	x0 -= camera.x
	x1 -= camera.x
	y0 -= camera.y
	y1 -= camera.y

	// Bresenham algorithm: https://www.youtube.com/watch?v=IDFB5CDpLDE
	run := float64(x1 - x0)
	if run == 0 {
		verticalLine(x0, y0, y1)
		return
	}

	rise := float64(y1 - y0)
	if rise == 0 {
		horizontalLine(y0, x0, x1)
		return
	}

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
			pset(x, y)

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
			pset(x, y)

			offset += delta
			if offset >= threshold {
				x += adjust
				threshold += 1
			}
		}
	}
}

// verticalLine draws a vertical line between y0-y1 inclusive
func verticalLine(x, y0, y1 int) {
	if y0 > y1 {
		y0, y1 = y1, y0
	}

	if x < clippingRegion.x {
		return
	}

	if x >= clippingRegion.x+clippingRegion.w {
		return
	}

	if y0 < clippingRegion.y {
		y0 = clippingRegion.y
	}

	if y1 >= clippingRegion.y+clippingRegion.h {
		y1 = clippingRegion.y + clippingRegion.h - 1
	}

	for y := y0; y <= y1; y++ {
		ScreenData[y*scrWidth+x] = drawPalette[color]
	}
}

// horizontalLine draws a vertical line between x0-x1 inclusive
func horizontalLine(y, x0, x1 int) {
	if y < clippingRegion.y {
		return
	}

	if y >= clippingRegion.y+clippingRegion.h {
		return
	}

	if x0 > x1 {
		x0, x1 = x1, x0
	}

	if x0 < clippingRegion.x {
		x0 = clippingRegion.x
	}

	if x1 >= clippingRegion.x+clippingRegion.w {
		x1 = clippingRegion.x + clippingRegion.w - 1
	}

	for x := x0; x <= x1; x++ {
		ScreenData[y*scrWidth+x] = drawPalette[color]
	}
}
