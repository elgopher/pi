// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elgopher/pi"
)

func WindowAutoSize(monitor *ebiten.MonitorType) (w, h, minW, minH int) {
	monitorWidth, monitorHeight := monitor.Size() // already scaled
	deviceScaleFactor := monitor.DeviceScaleFactor()

	screen := pi.Screen()
	// so we also scale the screen size
	screenWidth := float64(screen.W()) / deviceScaleFactor
	screenHeight := float64(screen.H()) / deviceScaleFactor

	widthRatio := float64(monitorWidth) / screenWidth
	heightRatio := float64(monitorHeight) / screenHeight

	ratio := math.Ceil(min(widthRatio, heightRatio) / 2.0)

	w = int(screenWidth * ratio)
	h = int(screenHeight * ratio)

	minW = int(screenWidth)
	minH = int(screenHeight)

	w = adjustSize(w, deviceScaleFactor, float64(screen.W())*ratio)
	h = adjustSize(h, deviceScaleFactor, float64(screen.H())*ratio)
	minW = adjustSize(minW, deviceScaleFactor, float64(screen.W()))
	minH = adjustSize(minH, deviceScaleFactor, float64(screen.H()))

	return
}

// Ebitengine workaround - Ebitengine does not allow setting the screen size in real pixels,
// but only in device-independent pixels. These are scaled pixels, and unfortunately,
// I cannot pass them as floats. Therefore, I search for the closest integer value
// which, when multiplied by deviceScaleFactor, will give the expected window size.
func adjustSize(size int, deviceScaleFactor float64, screenSize float64) int {
	for ; float64(size)*deviceScaleFactor < screenSize; size += 1 {
	}
	return size
}
