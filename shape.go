// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import "math"

// RectFill draws a filled rectangle between points x0,y0 and x1,y1 (inclusive).
//
// RectFill takes into account camera position, clipping region and draw palette.
func RectFill(x0, y0, x1, y1 int, color byte) {
	xmin, xmax := x0-screen.Camera.X, x1-screen.Camera.X
	if xmin > xmax {
		xmin, xmax = xmax, xmin
	}
	ymin, ymax := y0-screen.Camera.Y, y1-screen.Camera.Y
	if ymin > ymax {
		ymin, ymax = ymax, ymin
	}

	if xmin >= screen.Clip.X+screen.Clip.W {
		return
	}

	if xmax < screen.Clip.X {
		return
	}

	if ymin >= screen.Clip.Y+screen.Clip.H {
		return
	}

	if ymax < screen.Clip.Y {
		return
	}

	if xmin < screen.Clip.X {
		xmin = screen.Clip.X
	}

	if xmax >= screen.Clip.X+screen.Clip.W {
		xmax = screen.Clip.X + screen.Clip.W - 1
	}

	if ymin < screen.Clip.Y {
		ymin = screen.Clip.Y
	}

	if ymax >= screen.Clip.Y+screen.Clip.H {
		ymax = screen.Clip.Y + screen.Clip.H - 1
	}

	w := xmax - xmin + 1
	col := DrawPalette[color]
	for i := 0; i < w; i++ {
		screen.lineOfScreenWidth[i] = col
	}
	line := screen.lineOfScreenWidth[:w]

	for y := ymin; y <= ymax; y++ {
		copy(screen.Pix[y*screen.W+xmin:], line)
	}
}

// Rect draws a rectangle between points x0,y0 and x1,y1 (inclusive).
//
// Rect takes into account camera position, clipping region and draw palette.
func Rect(x0, y0, x1, y1 int, color byte) {
	xmin, xmax := x0-screen.Camera.X, x1-screen.Camera.X
	if xmin > xmax {
		xmin, xmax = xmax, xmin
	}
	ymin, ymax := y0-screen.Camera.Y, y1-screen.Camera.Y
	if ymin > ymax {
		ymin, ymax = ymax, ymin
	}

	if xmin >= screen.Clip.X+screen.Clip.W {
		return
	}

	if xmax < screen.Clip.X {
		return
	}

	if ymin >= screen.Clip.Y+screen.Clip.H {
		return
	}

	if ymax < screen.Clip.Y {
		return
	}

	drawLeftLine := true
	drawRightLine := true

	if xmin < screen.Clip.X {
		xmin = screen.Clip.X
		drawLeftLine = false
	}

	if xmax >= screen.Clip.X+screen.Clip.W {
		xmax = screen.Clip.X + screen.Clip.W - 1
		drawRightLine = false
	}

	col := DrawPalette[color]

	w := xmax - xmin + 1
	for i := 0; i < w; i++ {
		screen.lineOfScreenWidth[i] = col
	}
	line := screen.lineOfScreenWidth[:w]

	if ymin < screen.Clip.Y {
		ymin = screen.Clip.Y
	} else {
		copy(screen.Pix[ymin*screen.W+xmin:], line)
	}

	if ymax >= screen.Clip.Y+screen.Clip.H {
		ymax = screen.Clip.Y + screen.Clip.H - 1
	} else {
		copy(screen.Pix[ymax*screen.W+xmin:], line)
	}

	if drawLeftLine {
		for y := ymin; y <= ymax; y++ {
			screen.Pix[y*screen.W+xmin] = col
		}
	}

	if drawRightLine {
		for y := ymin; y <= ymax; y++ {
			screen.Pix[y*screen.W+xmax] = col
		}
	}
}

// Line draws a line between points x0,y0 and x1,y1 (inclusive).
//
// Line takes into account camera position, clipping region and draw palette.
func Line(x0, y0, x1, y1 int, color byte) {
	camera := screen.Camera
	x0 -= camera.X
	x1 -= camera.X
	y0 -= camera.Y
	y1 -= camera.Y

	// Bresenham algorithm: https://www.youtube.com/watch?v=IDFB5CDpLDE
	run := float64(x1 - x0)
	if run == 0 {
		verticalLine(x0, y0, y1, color)
		return
	}

	rise := float64(y1 - y0)
	if rise == 0 {
		horizontalLine(y0, x0, x1, color)
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
			pset(x, y, color)

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
			pset(x, y, color)

			offset += delta
			if offset >= threshold {
				x += adjust
				threshold += 1
			}
		}
	}
}

// verticalLine draws a vertical line between y0-y1 inclusive
func verticalLine(x, y0, y1 int, color byte) {
	if y0 > y1 {
		y0, y1 = y1, y0
	}

	if x < screen.Clip.X {
		return
	}

	if x >= screen.Clip.X+screen.Clip.W {
		return
	}

	if y0 < screen.Clip.Y {
		y0 = screen.Clip.Y
	}

	if y1 >= screen.Clip.Y+screen.Clip.H {
		y1 = screen.Clip.Y + screen.Clip.H - 1
	}

	for y := y0; y <= y1; y++ {
		screen.Pix[y*screen.W+x] = DrawPalette[color]
	}
}

// horizontalLine draws a vertical line between x0-x1 inclusive
func horizontalLine(y, x0, x1 int, color byte) {
	if y < screen.Clip.Y {
		return
	}

	if y >= screen.Clip.Y+screen.Clip.H {
		return
	}

	if x0 > x1 {
		x0, x1 = x1, x0
	}

	if x0 < screen.Clip.X {
		x0 = screen.Clip.X
	}

	if x1 >= screen.Clip.X+screen.Clip.W {
		x1 = screen.Clip.X + screen.Clip.W - 1
	}

	offset := y * screen.W

	for x := x0; x <= x1; x++ {
		screen.Pix[offset+x] = DrawPalette[color]
	}
}

// Circ draws a circle
//
// Circ takes into account camera position, clipping region and draw palette.
func Circ(centerX, centerY, radius int, color byte) {
	centerX = centerX - screen.Camera.X
	centerY = centerY - screen.Camera.Y
	// Code based on Frédéric Goset work: http://fredericgoset.ovh/mathematiques/courbes/en/bresenham_circle.html
	x := 0
	y := radius
	m := 5 - 4*radius

	for x <= y {
		pset(centerX+x, centerY+y, color)
		pset(centerX+x, centerY-y, color)
		pset(centerX-x, centerY+y, color)
		pset(centerX-x, centerY-y, color)
		pset(centerX+y, centerY+x, color)
		pset(centerX+y, centerY-x, color)
		pset(centerX-y, centerY+x, color)
		pset(centerX-y, centerY-x, color)

		if m > 0 {
			y--
			m -= 8 * y
		}

		x++

		m += 8*x + 4
	}
}

// CircFill draws a filled circle
//
// CircFill takes into account camera position, clipping region and draw palette.
func CircFill(centerX, centerY, radius int, color byte) {
	// Code based on Frédéric Goset work: http://fredericgoset.ovh/mathematiques/courbes/en/filled_circle.html
	x := 0
	y := radius
	m := 5 - 4*radius

	for x <= y {
		RectFill(centerX-y, centerY-x, centerX+y, centerY-x, color)
		RectFill(centerX-y, centerY+x, centerX+y, centerY+x, color)

		if m > 0 {
			RectFill(centerX-x, centerY-y, centerX+x, centerY-y, color)
			RectFill(centerX-x, centerY+y, centerX+x, centerY+y, color)
			y--
			m -= 8 * y
		}

		x++

		m += 8*x + 4
	}
}
