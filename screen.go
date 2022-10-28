// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import "fmt"

var screen = newScreen(defaultScreenWidth, defaultScreenHeight)

// Cls cleans the entire screen with color 0. It does not take into account any draw state parameters such as clipping region or camera.
func Cls() {
	cls()
	ClipReset()
}

func cls() {
	copy(screen.Pix, screen.zeroScreenData)
}

// ClsCol cleans the entire screen with specified color. It does not take into account any draw state parameters such as clipping region or camera.
func ClsCol(col byte) {
	clsCol(col)
	ClipReset()
}

func clsCol(col byte) {
	for i := 0; i < len(screen.lineOfScreenWidth); i++ {
		screen.lineOfScreenWidth[i] = col
	}

	offset := 0
	for y := 0; y < screen.H; y++ {
		copy(screen.Pix[offset:offset+screen.W], screen.lineOfScreenWidth)
		offset += screen.W
	}
}

// Pset sets a pixel color on the screen.
func Pset(x, y int, color byte) {
	pset(x-screen.Camera.X, y-screen.Camera.Y, color)
}

// pset sets a pixel color on the screen **without** taking camera position into account.
func pset(x, y int, color byte) {
	if x < screen.Clip.X {
		return
	}
	if y < screen.Clip.Y {
		return
	}
	if x >= screen.Clip.X+screen.Clip.W {
		return
	}
	if y >= screen.Clip.Y+screen.Clip.H {
		return
	}

	screen.Pix[y*screen.W+x] = DrawPalette[color]
}

// Pget gets a pixel color on the screen.
func Pget(x, y int) byte {
	x -= screen.Camera.X
	y -= screen.Camera.Y

	if x < screen.Clip.X {
		return 0
	}
	if y < screen.Clip.Y {
		return 0
	}
	if x >= screen.Clip.X+screen.Clip.W {
		return 0
	}
	if y >= screen.Clip.Y+screen.Clip.H {
		return 0
	}

	return screen.Pix[y*screen.W+x]
}

// Clip sets the clipping region in the form of rectangle. All screen drawing operations will not affect any pixels outside the region.
//
// Clip returns previous clipping region.
func Clip(x, y, w, h int) (prevX, prevY, prevW, prevH int) {
	prev := screen.Clip

	if x < 0 {
		w += x
		x = 0
	}

	if y < 0 {
		h += y
		y = 0
	}

	if x+w > screen.W {
		w = screen.W - x
	}

	if y+h > screen.H {
		h = screen.H - y
	}

	screen.Clip.X = x
	screen.Clip.Y = y
	screen.Clip.W = w
	screen.Clip.H = h

	return prev.X, prev.Y, prev.W, prev.H
}

// ClipReset resets the clipping region, which means that entire screen will be clipped.
func ClipReset() (prevX, prevY, prevW, prevH int) {
	return Clip(0, 0, screen.W, screen.H)
}

//func ClipPrev(x, y, w, h int) {}

// Camera sets the camera offset used for all subsequent draw operations.
func Camera(x, y int) (prevX, prevY int) {
	prev := screen.Camera
	screen.Camera.X = x
	screen.Camera.Y = y
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
	if n >= sprSheet.numberOfSprites {
		return
	}

	x -= screen.Camera.X
	y -= screen.Camera.Y

	screenOffset := y*screen.W + x

	spriteX := (n % sprSheet.spritesInLine) * SpriteWidth
	spriteY := (n / sprSheet.spritesInLine) * SpriteHeight

	width := int(SpriteWidth * w)
	height := int(SpriteHeight * h)

	if spriteX+width > sprSheet.W {
		width = sprSheet.W - spriteX
	}

	if spriteY+height > sprSheet.H {
		height = sprSheet.H - spriteY
	}

	spriteSheetOffset := spriteY*sprSheet.W + spriteX

	if x < screen.Clip.X {
		dx := screen.Clip.X - x
		width -= dx
		screenOffset += dx
		spriteSheetOffset += dx
	} else if x+width > screen.Clip.W {
		width = screen.Clip.W - x
	}

	if width <= 0 {
		return
	}

	if y < screen.Clip.Y {
		dy := screen.Clip.Y - y
		height -= dy
		screenOffset += dy * screen.W
		spriteSheetOffset += dy * sprSheet.W
	} else if y+height > screen.Clip.H {
		height = screen.Clip.H - y
	}

	spriteSheetStep := sprSheet.W

	if flipY {
		spriteSheetOffset += (height - 1) * sprSheet.W
		spriteSheetStep = -sprSheet.W
	}

	startingPixel := 0
	step := 1

	if flipX {
		startingPixel = width - 1
		step = -1
	}

	for i := 0; i < height; i++ {
		spriteSheetLine := sprSheet.Pix[spriteSheetOffset : spriteSheetOffset+width]

		for j := 0; j < len(spriteSheetLine); j++ {
			col := spriteSheetLine[startingPixel+(step*j)]
			if ColorTransparency[col] {
				continue
			}

			screen.Pix[screenOffset+j] = DrawPalette[col]
		}
		screenOffset += screen.W
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

func Scr() Screen {
	return screen
}

type Screen struct {
	// Width and height in pixels
	W, H int

	// Pix contains pixel colors for the screen visible by the player.
	// Each pixel is one byte. It is initialized during pi.Boot.
	//
	// Pix on the screen are organized from left to right,
	// top to bottom. Slice element number 0 has pixel located
	// in the top-left corner. Slice element number 1 has pixel color
	// on the right and so on.
	//
	// Can be freely read and updated. Useful when you want to use your own
	// functions for pixel manipulation.
	Pix []byte

	Clip Region

	Camera Position

	zeroScreenData    []byte
	lineOfScreenWidth []byte
}

type Region struct {
	X, Y, W, H int
}

type Position struct {
	X, Y int
}

func newScreen(w, h int) Screen {
	screenSize := w * h

	return Screen{
		W:                 w,
		H:                 h,
		Pix:               make([]byte, screenSize),
		zeroScreenData:    make([]byte, screenSize),
		Clip:              Region{W: w, H: h},
		lineOfScreenWidth: make([]byte, w),
	}
}

func SetScreenSize(w, h int) {
	if w <= 0 {
		panic(fmt.Sprintf("screen width %d is not greather than 0", w))
	}
	if h <= 0 {
		panic(fmt.Sprintf("screen height %d is not greather than 0", h))
	}

	screen = newScreen(w, h)
}
