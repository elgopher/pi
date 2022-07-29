// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

// Screen-specific data
var (
	// ScreenData contains pixel colors for the screen visible by the player.
	// Each pixel is one byte. It is initialized during pi.Boot.
	//
	// Pixels on the screen are organized from left to right,
	// top to bottom. Slice element number 0 has pixel located
	// in the top-left corner. Slice element number 1 has pixel color
	// on the right and so on.
	//
	// Can be freely read and updated. Useful when you want to use your own
	// functions for pixel manipulation.
	// Pi will panic if you try to change the length of the slice.
	ScreenData []byte

	scrWidth, scrHeight int
	lineOfScreenWidth   []byte
	zeroScreenData      []byte
	clippingRegion      rect
	color               = defaultColor // Color is a currently used color in draw state. Used by Pset.
)

// Cls cleans the entire screen with color 0. It does not take into account any draw state parameters such as clipping region or camera.
func Cls() {
	copy(ScreenData, zeroScreenData)
}

// ClsCol cleans the entire screen with specified color. It does not take into account any draw state parameters such as clipping region or camera.
func ClsCol(col byte) {
	for i := 0; i < len(lineOfScreenWidth); i++ {
		lineOfScreenWidth[i] = col
	}

	offset := 0
	for y := 0; y < scrHeight; y++ {
		copy(ScreenData[offset:offset+scrWidth], lineOfScreenWidth)
		offset += scrWidth
	}
}

// Color sets the color used by Pset.
//
// Color returns previously used color.
func Color(col byte) (prevCol byte) {
	prevCol = color
	color = col
	return
}

// ColorReset resets the color to default value which is 6.
//
// ColorReset returns previously used color.
func ColorReset() (prevCol byte) {
	return Color(defaultColor)
}

// Pset sets a pixel color on the screen to Color.
func Pset(x, y int) {
	x -= camera.x
	y -= camera.y

	if x < 0 {
		return
	}
	if y < 0 {
		return
	}
	if x >= scrWidth {
		return
	}
	if y >= scrHeight {
		return
	}
	if x < clippingRegion.x {
		return
	}
	if y < clippingRegion.y {
		return
	}
	if x > clippingRegion.w {
		return
	}
	if y > clippingRegion.h {
		return
	}

	ScreenData[y*scrWidth+x] = drawPalette[color]
}

// Pget gets a pixel color on the screen.
func Pget(x, y int) byte {
	x -= camera.x
	y -= camera.y

	if x < 0 {
		return 0
	}
	if y < 0 {
		return 0
	}
	if x >= scrWidth {
		return 0
	}
	if y >= scrHeight {
		return 0
	}
	if x < clippingRegion.x {
		return 0
	}
	if y < clippingRegion.y {
		return 0
	}
	if x > clippingRegion.w {
		return 0
	}
	if y > clippingRegion.h {
		return 0
	}

	return ScreenData[y*scrWidth+x]
}

type rect struct {
	x, y, w, h int
}

// Clip sets the clipping region in the form of rectangle. All screen drawing operations will not affect any pixels outside the region.
//
// Clip returns previous clipping region.
func Clip(x, y, w, h int) (prevX, prevY, prevW, prevH int) {
	prev := clippingRegion

	if x < 0 {
		w += x
		x = 0
	}

	if y < 0 {
		h += y
		y = 0
	}

	if x+w > scrWidth {
		w = scrWidth - x
	}

	if y+h > scrHeight {
		h = scrHeight - y
	}

	clippingRegion.x = x
	clippingRegion.y = y
	clippingRegion.w = w
	clippingRegion.h = h

	return prev.x, prev.y, prev.w, prev.h
}

// ClipReset resets the clipping region, which means that entire screen will be clipped.
func ClipReset() (prevX, prevY, prevW, prevH int) {
	return Clip(0, 0, scrWidth, scrHeight)
}

//func ClipPrev(x, y, w, h int) {}

type pos struct {
	x, y int
}

var camera pos

// Camera sets the camera offset used for all subsequent draw operations.
func Camera(x, y int) (prevX, prevY int) {
	prev := camera
	camera.x = x
	camera.y = y
	return prev.x, prev.y
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

	x -= camera.x
	y -= camera.y

	screenOffset := y*scrWidth + x

	spriteX := (n % spritesInLine) * SpriteWidth
	spriteY := (n / spritesInLine) * SpriteHeight

	width := int(SpriteWidth * w)
	height := int(SpriteHeight * h)

	if spriteX+width > ssWidth {
		width = ssWidth - spriteX
	}

	if spriteY+height > ssHeight {
		height = ssHeight - spriteY
	}

	spriteSheetOffset := spriteY*ssWidth + spriteX

	if x < clippingRegion.x {
		dx := clippingRegion.x - x
		width -= dx
		screenOffset += dx
		spriteSheetOffset += dx
	} else if x+width > clippingRegion.w {
		width = clippingRegion.w - x
	}

	if width <= 0 {
		return
	}

	if y < clippingRegion.y {
		dy := clippingRegion.y - y
		height -= dy
		screenOffset += dy * scrWidth
		spriteSheetOffset += dy * ssWidth
	} else if y+height > clippingRegion.h {
		height = clippingRegion.h - y
	}

	spriteSheetStep := ssWidth

	if flipY {
		spriteSheetOffset += (height - 1) * ssWidth
		spriteSheetStep = -ssWidth
	}

	startingPixel := 0
	step := 1

	if flipX {
		startingPixel = width - 1
		step = -1
	}

	for i := 0; i < height; i++ {
		spriteSheetLine := SpriteSheetData[spriteSheetOffset : spriteSheetOffset+width]

		for j := 0; j < len(spriteSheetLine); j++ {
			col := spriteSheetLine[startingPixel+(step*j)]
			if colorIsTransparent[col] {
				continue
			}

			ScreenData[screenOffset+j] = drawPalette[col]
		}
		screenOffset += scrWidth
		spriteSheetOffset += spriteSheetStep
	}
}

var colorIsTransparent [256]bool

// Palt sets color transparency. If true then the color will not be drawn
// for next drawing operations.
//
// Color transparency is used by Spr, SprSize and SprSizeFlip.
func Palt(color byte, transparent bool) {
	colorIsTransparent[color] = transparent
}

var defaultTransparency = [256]bool{true}

// PaltReset sets all transparent colors to false and makes color 0 transparent.
func PaltReset() {
	colorIsTransparent = defaultTransparency
}

var drawPalette [256]byte

// Pal replaces color with another one for all subsequent drawings (it is changing
// the so-called draw palette).
//
// Affected functions are Pset, Spr, SprSize and SprSizeFlip.
func Pal(color byte, replacementColor byte) {
	drawPalette[color] = replacementColor
}

var displayPalette [256]byte

// PalDisplay replaces color with another one for the whole screen at the end of a frame
// (it is changing the so-called display palette).
func PalDisplay(color byte, replacementColor byte) {
	displayPalette[color] = replacementColor
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
	drawPalette = notSwappedPalette
	displayPalette = notSwappedPalette
}
