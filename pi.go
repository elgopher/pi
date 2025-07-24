// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package pi provides the core Pi functions for the game loop,
// screen, color palette, and drawing pixels, shapes, and sprites.
//
// This package and all other pi* packages are not thread-safe.
// This is an intentional design choice to significantly improve performance.
// You should not call Pi's API from any goroutine other than the one
// running your `pi.Update` and `pi.Draw` functions. You can still
// create your own goroutines in your game, but they must not call
// any Pi functions or access Pi state.
package pi

import "strconv"

// MaxColors is the maximum number of colors that can be used simultaneously on the screen.
const MaxColors = 64

var ticksPerSecond = 30

// SetTPS sets the number of ticks per second. This determines how often
// the Update and Draw functions are called.
//
// tps must be one of: 60, 30, or 15.
//
// You can change the TPS while the game is running.
func SetTPS(tps int) {
	if tps != 60 && tps != 30 && tps != 15 {
		panic("tps must be one of: 60, 30, 15, but was: " + strconv.Itoa(tps) + "")
	}
	ticksPerSecond = tps
}

// TPS returns the number of ticks per second. The default value is 30.
func TPS() int {
	return ticksPerSecond
}

var (
	// Frame is the current game frame number.
	//
	// It is automatically incremented by the backend at the start of each game frame.
	Frame int

	// Time is the current game time in seconds.
	//
	// It is automatically incremented by the backend at the start of each frame.
	Time float64
)

// Camera is the camera offset applied to all subsequent draw operations.
var Camera Position

var (
	drawColor  Color = 7
	drawTarget Canvas
	clip       IntArea
)

// Color represents a pixel color value in the range 0..63 (first 6 bits).
// Bits 6 and 7 specify the ColorTable index.
type Color = uint8

// Number describes any numeric type in Go.
//
// Includes signed integers, unsigned integers, and floating-point types.
type Number interface {
	~int | ~float64 |
		~int8 | ~int16 | ~int32 | ~int64 |
		~float32 |
		~uint | ~byte | ~uint16 | ~uint32 | ~uint64
}

// SetDrawTarget sets c as the target Canvas for all subsequent drawing,
// including functions like DrawSprite, SetPixel, Line, etc.
//
// This function also automatically sets the clip region to cover the entire area of c.
func SetDrawTarget(c Canvas) (prev Canvas) {
	prev = drawTarget
	drawTarget = c
	SetClip(c.EntireArea())
	return
}

func DrawTarget() Canvas {
	return drawTarget
}

// SetColor sets the current draw color.
//
// Returns the previous color.
func SetColor(c Color) (prev Color) {
	prev = drawColor
	drawColor = c
	return
}

// GetColor returns the current draw color.
func GetColor() Color {
	return drawColor
}

// SetClip sets the clipping region to the specified area.
//
// Returns the previous clipping area.
func SetClip(area IntArea) (prev IntArea) {
	prev = clip
	clip, _, _ = area.ClippedBy(IntArea{W: drawTarget.width, H: drawTarget.height})
	return
}

func Clip() IntArea {
	return clip
}

func setPixelWithColor(x, y int, draw Color) {
	x -= Camera.X
	y -= Camera.Y

	if x < clip.X {
		return
	}
	if y < clip.Y {
		return
	}
	if x >= clip.X+clip.W {
		return
	}
	if y >= clip.Y+clip.H {
		return
	}

	idx := y*drawTarget.width + x
	target := drawTarget.data[idx] & ShapeTargetMask

	drawTarget.data[idx] = ColorTables[(draw|target)>>6][drawColor&(MaxColors-1)][target&(MaxColors-1)]
}

// SetPixel sets the draw color at the given coordinates.
//
// It takes into account the camera position, clipping region,
// color tables, and masks.
func SetPixel(x, y int) {
	setPixelWithColor(x, y, drawColor&ReadMask)
}

// GetPixel returns the color at the given coordinates.
//
// It takes into account the camera position and the clipping region.
func GetPixel(x, y int) (color Color) {
	x -= Camera.X
	y -= Camera.Y

	if x < clip.X {
		return
	}
	if y < clip.Y {
		return
	}
	if x >= clip.X+clip.W {
		return
	}
	if y >= clip.Y+clip.H {
		return
	}

	return drawTarget.data[y*drawTarget.width+x]
}
