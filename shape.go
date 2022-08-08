// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

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
