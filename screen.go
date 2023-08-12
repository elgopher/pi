// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import "fmt"

var screen = NewPixMap(defaultScreenWidth, defaultScreenHeight)

// Cls cleans the entire screen with color 0. It does not take into account any global state such as clipping region or camera.
// Cls also resets clipping region.
func Cls() {
	screen.Clear()
	ClipReset()
}

// ClsCol cleans the entire screen with specified color. It does not take into account any global state such as clipping region or camera.
// ClsCol also resets clipping region.
func ClsCol(col byte) {
	screen.ClearCol(col)
	ClipReset()
}

// Pset sets a pixel color on the screen. It takes into account camera and draw palette.
func Pset(x, y int, color byte) {
	screen.Set(x-ScreenCamera.X, y-ScreenCamera.Y, DrawPalette[color])
}

// Pget gets a pixel color on the screen.
func Pget(x, y int) byte {
	x -= ScreenCamera.X
	y -= ScreenCamera.Y

	return screen.Get(x, y)
}

// Clip sets the clipping region in the form of rectangle. All screen drawing operations will not affect any pixels outside the region.
//
// Clip returns previous clipping region.
func Clip(x, y, w, h int) (prevX, prevY, prevW, prevH int) {
	prev := screen.clip
	screen = screen.WithClip(x, y, w, h)
	return prev.X, prev.Y, prev.W, prev.H
}

// ClipReset resets the clipping region, which means that entire screen will be clipped.
func ClipReset() (prevX, prevY, prevW, prevH int) {
	return Clip(0, 0, screen.Width(), screen.Height())
}

//func ClipPrev(x, y, w, h int) {}

// Camera sets the camera offset used for all subsequent draw operations.
func Camera(x, y int) (prevX, prevY int) {
	prev := ScreenCamera
	ScreenCamera.X = x
	ScreenCamera.Y = y
	return prev.X, prev.Y
}

var ScreenCamera Position

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
	if n >= sprSheet.numberOfSprites {
		return
	}

	x -= ScreenCamera.X
	y -= ScreenCamera.Y

	screenWidth := screen.Width()

	screenOffset := y*screenWidth + x

	spriteX := (n % sprSheet.spritesInLine) * SpriteWidth
	spriteY := (n / sprSheet.spritesInLine) * SpriteHeight

	width := int(SpriteWidth * w)
	height := int(SpriteHeight * h)

	sprSheetWidth := sprSheet.clip.W
	sprSheetHeight := sprSheet.clip.H

	if spriteX+width > sprSheetWidth {
		width = sprSheetWidth - spriteX
	}

	if spriteY+height > sprSheetHeight {
		height = sprSheetHeight - spriteY
	}

	spriteSheetOffset := spriteY*sprSheetWidth + spriteX

	if x < screen.clip.X {
		dx := screen.clip.X - x
		width -= dx
		screenOffset += dx
		spriteSheetOffset += dx
	} else if x+width > screen.clip.W {
		width = screen.clip.W - x
	}

	if width <= 0 {
		return
	}

	if y < screen.clip.Y {
		dy := screen.clip.Y - y
		height -= dy
		screenOffset += dy * screenWidth
		spriteSheetOffset += dy * sprSheetWidth
	} else if y+height > screen.clip.H {
		height = screen.clip.H - y
	}

	spriteSheetStep := sprSheetWidth

	if flipY {
		spriteSheetOffset += (height - 1) * sprSheetWidth
		spriteSheetStep = -sprSheetWidth
	}

	startingPixel := 0
	step := 1

	if flipX {
		startingPixel = width - 1
		step = -1
	}

	for i := 0; i < height; i++ {
		spriteSheetLine := sprSheet.pix[spriteSheetOffset : spriteSheetOffset+width]

		for j := 0; j < len(spriteSheetLine); j++ {
			col := spriteSheetLine[startingPixel+(step*j)]
			if ColorTransparency[col] {
				continue
			}

			screen.pix[screenOffset+j] = DrawPalette[col]
		}
		screenOffset += screenWidth
		spriteSheetOffset += spriteSheetStep
	}
}

// Palt sets color transparency. If true then the color will not be drawn
// for next drawing operations.
//
// Color transparency is used by Spr, SprSize and SprSizeFlip.
func Palt(color byte, transparent bool) {
	ColorTransparency[color] = transparent
}

var defaultTransparency = [256]bool{true}

// PaltReset sets all transparent colors to false and makes color 0 transparent.
func PaltReset() {
	ColorTransparency = defaultTransparency
}

// Pal replaces color with another one for all subsequent drawings (it is changing
// the so-called draw palette).
//
// Affected functions are Pset, Spr, SprSize, SprSizeFlip, Rect and RectFill.
func Pal(color byte, replacementColor byte) {
	DrawPalette[color] = replacementColor
}

// PalDisplay replaces color with another one for the whole screen at the end of a frame
// (it is changing the so-called display palette).
func PalDisplay(color byte, replacementColor byte) {
	DisplayPalette[color] = replacementColor
}

// PalSecondary

var notSwappedPalette [256]byte

func init() {
	for i := 0; i < 256; i++ {
		c := byte(i)
		notSwappedPalette[i] = c
		DrawPalette[i] = c
		DisplayPalette[i] = c
	}
}

// PalReset resets all swapped colors for all palettes.
func PalReset() {
	DrawPalette = notSwappedPalette
	DisplayPalette = notSwappedPalette
}

// Scr returns the Screen PixMap
func Scr() PixMap {
	return screen
}

type Region struct {
	X, Y, W, H int
}

type Position struct {
	X, Y int
}

const maxScreenSize = 1024 * 64

// SetScreenSize sets the screen size to specified resolution. The maximum number of pixels is 65536 (64KB).
// Will panic if screen size is too big or width/height are <= 0.
func SetScreenSize(width, height int) {
	if width <= 0 {
		panic(fmt.Sprintf("screen width %d is not greather than 0", width))
	}
	if height <= 0 {
		panic(fmt.Sprintf("screen height %d is not greather than 0", height))
	}

	if width*height > maxScreenSize {
		panic(fmt.Sprintf("number of pixels for screen resolution %dx%d is higher than maximum %d. Please use smaller screen.", width, height, maxScreenSize))
	}

	screen = NewPixMap(width, height)
}
