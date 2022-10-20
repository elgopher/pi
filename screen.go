// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"fmt"
	"image"
	stdcolor "image/color"
	"image/png"
	"os"
	"runtime"

	"github.com/elgopher/pi/vm"
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
	copy(vm.ScreenData, zeroScreenData)
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
	for y := 0; y < vm.ScreenHeight; y++ {
		copy(vm.ScreenData[offset:offset+vm.ScreenWidth], lineOfScreenWidth)
		offset += vm.ScreenWidth
	}
}

// Pset sets a pixel color on the screen.
func Pset(x, y int, color byte) {
	pset(x-vm.Camera.X, y-vm.Camera.Y, color)
}

// pset sets a pixel color on the screen **without** taking camera position into account.
func pset(x, y int, color byte) {
	if x < vm.ClippingRegion.X {
		return
	}
	if y < vm.ClippingRegion.Y {
		return
	}
	if x >= vm.ClippingRegion.X+vm.ClippingRegion.W {
		return
	}
	if y >= vm.ClippingRegion.Y+vm.ClippingRegion.H {
		return
	}

	vm.ScreenData[y*vm.ScreenWidth+x] = vm.DrawPalette[color]
}

// Pget gets a pixel color on the screen.
func Pget(x, y int) byte {
	x -= vm.Camera.X
	y -= vm.Camera.Y

	if x < vm.ClippingRegion.X {
		return 0
	}
	if y < vm.ClippingRegion.Y {
		return 0
	}
	if x >= vm.ClippingRegion.X+vm.ClippingRegion.W {
		return 0
	}
	if y >= vm.ClippingRegion.Y+vm.ClippingRegion.H {
		return 0
	}

	return vm.ScreenData[y*vm.ScreenWidth+x]
}

// Clip sets the clipping region in the form of rectangle. All screen drawing operations will not affect any pixels outside the region.
//
// Clip returns previous clipping region.
func Clip(x, y, w, h int) (prevX, prevY, prevW, prevH int) {
	prev := vm.ClippingRegion

	if x < 0 {
		w += x
		x = 0
	}

	if y < 0 {
		h += y
		y = 0
	}

	if x+w > vm.ScreenWidth {
		w = vm.ScreenWidth - x
	}

	if y+h > vm.ScreenHeight {
		h = vm.ScreenHeight - y
	}

	vm.ClippingRegion.X = x
	vm.ClippingRegion.Y = y
	vm.ClippingRegion.W = w
	vm.ClippingRegion.H = h

	return prev.X, prev.Y, prev.W, prev.H
}

// ClipReset resets the clipping region, which means that entire screen will be clipped.
func ClipReset() (prevX, prevY, prevW, prevH int) {
	return Clip(0, 0, vm.ScreenWidth, vm.ScreenHeight)
}

//func ClipPrev(x, y, w, h int) {}

// Camera sets the camera offset used for all subsequent draw operations.
func Camera(x, y int) (prevX, prevY int) {
	prev := vm.Camera
	vm.Camera.X = x
	vm.Camera.Y = y
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

	x -= vm.Camera.X
	y -= vm.Camera.Y

	screenOffset := y*vm.ScreenWidth + x

	spriteX := (n % spritesInLine) * SpriteWidth
	spriteY := (n / spritesInLine) * SpriteHeight

	width := int(SpriteWidth * w)
	height := int(SpriteHeight * h)

	if spriteX+width > vm.SpriteSheetWidth {
		width = vm.SpriteSheetWidth - spriteX
	}

	if spriteY+height > vm.SpriteSheetHeight {
		height = vm.SpriteSheetHeight - spriteY
	}

	spriteSheetOffset := spriteY*vm.SpriteSheetWidth + spriteX

	if x < vm.ClippingRegion.X {
		dx := vm.ClippingRegion.X - x
		width -= dx
		screenOffset += dx
		spriteSheetOffset += dx
	} else if x+width > vm.ClippingRegion.W {
		width = vm.ClippingRegion.W - x
	}

	if width <= 0 {
		return
	}

	if y < vm.ClippingRegion.Y {
		dy := vm.ClippingRegion.Y - y
		height -= dy
		screenOffset += dy * vm.ScreenWidth
		spriteSheetOffset += dy * vm.SpriteSheetWidth
	} else if y+height > vm.ClippingRegion.H {
		height = vm.ClippingRegion.H - y
	}

	spriteSheetStep := vm.SpriteSheetWidth

	if flipY {
		spriteSheetOffset += (height - 1) * vm.SpriteSheetWidth
		spriteSheetStep = -vm.SpriteSheetWidth
	}

	startingPixel := 0
	step := 1

	if flipX {
		startingPixel = width - 1
		step = -1
	}

	for i := 0; i < height; i++ {
		spriteSheetLine := vm.SpriteSheetData[spriteSheetOffset : spriteSheetOffset+width]

		for j := 0; j < len(spriteSheetLine); j++ {
			col := spriteSheetLine[startingPixel+(step*j)]
			if vm.ColorTransparency[col] {
				continue
			}

			vm.ScreenData[screenOffset+j] = vm.DrawPalette[col]
		}
		screenOffset += vm.ScreenWidth
		spriteSheetOffset += spriteSheetStep
	}
}

// Palt sets color transparency. If true then the color will not be drawn
// for next drawing operations.
//
// Color transparency is used by Spr, SprSize and SprSizeFlip.
func Palt(color byte, transparent bool) {
	vm.ColorTransparency[color] = transparent
}

var defaultTransparency = [256]bool{true}

// PaltReset sets all transparent colors to false and makes color 0 transparent.
func PaltReset() {
	vm.ColorTransparency = defaultTransparency
}

// Pal replaces color with another one for all subsequent drawings (it is changing
// the so-called draw palette).
//
// Affected functions are Pset, Spr, SprSize, SprSizeFlip, Rect and RectFill.
func Pal(color byte, replacementColor byte) {
	vm.DrawPalette[color] = replacementColor
}

// PalDisplay replaces color with another one for the whole screen at the end of a frame
// (it is changing the so-called display palette).
func PalDisplay(color byte, replacementColor byte) {
	vm.DisplayPalette[color] = replacementColor
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
	vm.DrawPalette = notSwappedPalette
	vm.DisplayPalette = notSwappedPalette
}

// Snap takes a screenshot and saves it to temp dir.
//
// Snap returns a filename. If something went wrong error is returned.
func Snap() (string, error) {
	if runtime.GOOS == "js" {
		return "", fmt.Errorf("storing files does not work on js")
	}

	var palette stdcolor.Palette
	for _, col := range vm.DisplayPalette {
		rgb := vm.Palette[col]
		rgba := &stdcolor.NRGBA{R: rgb.R, G: rgb.G, B: rgb.B, A: 255}
		palette = append(palette, rgba)
	}

	size := image.Rectangle{Max: image.Point{X: vm.ScreenWidth, Y: vm.ScreenHeight}}
	img := image.NewPaletted(size, palette)

	copy(img.Pix, vm.ScreenData)

	file, err := os.CreateTemp("", "pi-screenshot")
	if err != nil {
		return "", fmt.Errorf("error creating temp file for screenshot: %w", err)
	}

	if err = png.Encode(file, img); err != nil {
		return "", fmt.Errorf("error encoding screenshot into PNG file: %w", err)
	}

	if err = file.Close(); err != nil {
		return "", fmt.Errorf("error closing file: %w", err)
	}

	return file.Name(), nil
}
