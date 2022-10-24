// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"github.com/elgopher/pi/mem"
)

// Screen-specific data
var (
	lineOfScreenWidth []byte
	zeroScreenData    []byte
)

// Cls cleans the entire screen with color 0. It does not take into account any draw state parameters such as clipping region or camera.
func Cls() {
	cls()
	ClipReset()
}

func cls() {
	copy(mem.ScreenData, zeroScreenData)
}

// ClsCol cleans the entire screen with specified color. It does not take into account any draw state parameters such as clipping region or camera.
func ClsCol(col byte) {
	clsCol(col)
	ClipReset()
}

func clsCol(col byte) {
	for i := 0; i < len(lineOfScreenWidth); i++ {
		lineOfScreenWidth[i] = col
	}

	offset := 0
	for y := 0; y < mem.ScreenHeight; y++ {
		copy(mem.ScreenData[offset:offset+mem.ScreenWidth], lineOfScreenWidth)
		offset += mem.ScreenWidth
	}
}

// Pset sets a pixel color on the screen.
func Pset(x, y int, color byte) {
	pset(x-mem.Camera.X, y-mem.Camera.Y, color)
}

// pset sets a pixel color on the screen **without** taking camera position into account.
func pset(x, y int, color byte) {
	if x < mem.ClippingRegion.X {
		return
	}
	if y < mem.ClippingRegion.Y {
		return
	}
	if x >= mem.ClippingRegion.X+mem.ClippingRegion.W {
		return
	}
	if y >= mem.ClippingRegion.Y+mem.ClippingRegion.H {
		return
	}

	mem.ScreenData[y*mem.ScreenWidth+x] = mem.DrawPalette[color]
}

// Pget gets a pixel color on the screen.
func Pget(x, y int) byte {
	x -= mem.Camera.X
	y -= mem.Camera.Y

	if x < mem.ClippingRegion.X {
		return 0
	}
	if y < mem.ClippingRegion.Y {
		return 0
	}
	if x >= mem.ClippingRegion.X+mem.ClippingRegion.W {
		return 0
	}
	if y >= mem.ClippingRegion.Y+mem.ClippingRegion.H {
		return 0
	}

	return mem.ScreenData[y*mem.ScreenWidth+x]
}

// Clip sets the clipping region in the form of rectangle. All screen drawing operations will not affect any pixels outside the region.
//
// Clip returns previous clipping region.
func Clip(x, y, w, h int) (prevX, prevY, prevW, prevH int) {
	prev := mem.ClippingRegion

	if x < 0 {
		w += x
		x = 0
	}

	if y < 0 {
		h += y
		y = 0
	}

	if x+w > mem.ScreenWidth {
		w = mem.ScreenWidth - x
	}

	if y+h > mem.ScreenHeight {
		h = mem.ScreenHeight - y
	}

	mem.ClippingRegion.X = x
	mem.ClippingRegion.Y = y
	mem.ClippingRegion.W = w
	mem.ClippingRegion.H = h

	return prev.X, prev.Y, prev.W, prev.H
}

// ClipReset resets the clipping region, which means that entire screen will be clipped.
func ClipReset() (prevX, prevY, prevW, prevH int) {
	return Clip(0, 0, mem.ScreenWidth, mem.ScreenHeight)
}

//func ClipPrev(x, y, w, h int) {}

// Camera sets the camera offset used for all subsequent draw operations.
func Camera(x, y int) (prevX, prevY int) {
	prev := mem.Camera
	mem.Camera.X = x
	mem.Camera.Y = y
	return prev.X, prev.Y
}

// CameraReset resets the camera offset to origin (0,0).
func CameraReset() (prevX, prevY int) {
	return Camera(0, 0)
}

// Spr draws a sprite with specified number on the screen.
// Sprites are counted from left to right, top to bottom. Sprite 0 is on top-left corner, sprite 1 is to the right and so on.
func Spr(n, x, y int) {
	SprSize(n, x, y, 1.0, 1.0)
}

// SprSize draws a range of sprites on the screen.
//
// n is a sprite number in the top-left corner.
//
// Non-integer w or h may be used to draw partial sprites.
func SprSize(n, x, y int, w, h float64) {
	SprSizeFlip(n, x, y, w, h, false, false)
}

// SprSizeFlip draws a range of sprites on the screen.
//
// If flipX is true then sprite is flipped horizontally.
// If flipY is true then sprite is flipped vertically.
func SprSizeFlip(n, x, y int, w, h float64, flipX, flipY bool) {
	if n < 0 {
		return
	}
	if n >= numberOfSprites {
		return
	}

	x -= mem.Camera.X
	y -= mem.Camera.Y

	screenOffset := y*mem.ScreenWidth + x

	spriteX := (n % spritesInLine) * SpriteWidth
	spriteY := (n / spritesInLine) * SpriteHeight

	width := int(SpriteWidth * w)
	height := int(SpriteHeight * h)

	if spriteX+width > mem.SpriteSheetWidth {
		width = mem.SpriteSheetWidth - spriteX
	}

	if spriteY+height > mem.SpriteSheetHeight {
		height = mem.SpriteSheetHeight - spriteY
	}

	spriteSheetOffset := spriteY*mem.SpriteSheetWidth + spriteX

	if x < mem.ClippingRegion.X {
		dx := mem.ClippingRegion.X - x
		width -= dx
		screenOffset += dx
		spriteSheetOffset += dx
	} else if x+width > mem.ClippingRegion.W {
		width = mem.ClippingRegion.W - x
	}

	if width <= 0 {
		return
	}

	if y < mem.ClippingRegion.Y {
		dy := mem.ClippingRegion.Y - y
		height -= dy
		screenOffset += dy * mem.ScreenWidth
		spriteSheetOffset += dy * mem.SpriteSheetWidth
	} else if y+height > mem.ClippingRegion.H {
		height = mem.ClippingRegion.H - y
	}

	spriteSheetStep := mem.SpriteSheetWidth

	if flipY {
		spriteSheetOffset += (height - 1) * mem.SpriteSheetWidth
		spriteSheetStep = -mem.SpriteSheetWidth
	}

	startingPixel := 0
	step := 1

	if flipX {
		startingPixel = width - 1
		step = -1
	}

	for i := 0; i < height; i++ {
		spriteSheetLine := mem.SpriteSheetData[spriteSheetOffset : spriteSheetOffset+width]

		for j := 0; j < len(spriteSheetLine); j++ {
			col := spriteSheetLine[startingPixel+(step*j)]
			if mem.ColorTransparency[col] {
				continue
			}

			mem.ScreenData[screenOffset+j] = mem.DrawPalette[col]
		}
		screenOffset += mem.ScreenWidth
		spriteSheetOffset += spriteSheetStep
	}
}

// Palt sets color transparency. If true then the color will not be drawn
// for next drawing operations.
//
// Color transparency is used by Spr, SprSize and SprSizeFlip.
func Palt(color byte, transparent bool) {
	mem.ColorTransparency[color] = transparent
}

var defaultTransparency = [256]bool{true}

// PaltReset sets all transparent colors to false and makes color 0 transparent.
func PaltReset() {
	mem.ColorTransparency = defaultTransparency
}

// Pal replaces color with another one for all subsequent drawings (it is changing
// the so-called draw palette).
//
// Affected functions are Pset, Spr, SprSize, SprSizeFlip, Rect and RectFill.
func Pal(color byte, replacementColor byte) {
	mem.DrawPalette[color] = replacementColor
}

// PalDisplay replaces color with another one for the whole screen at the end of a frame
// (it is changing the so-called display palette).
func PalDisplay(color byte, replacementColor byte) {
	mem.DisplayPalette[color] = replacementColor
}

// PalSecondary

var notSwappedPalette [256]byte

func init() {
	for i := 0; i < 256; i++ {
		notSwappedPalette[i] = byte(i)
	}
}

// PalReset resets all swapped colors for all palettes.
func PalReset() {
	mem.DrawPalette = notSwappedPalette
	mem.DisplayPalette = notSwappedPalette
}
