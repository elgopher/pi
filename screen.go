// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

// Screen-specific data
var (
	Color byte = 6 // Color is a currently used color in draw state. Used by Pset.

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

	ScreenData[y*scrWidth+x] = Color
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
// TODO Flipping is not implemented yet
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

	spriteX := n % spritesInLine
	spriteY := n / spritesInLine

	if spriteX+int(w) >= spritesInLine {
		w = float64(spritesInLine - spriteX)
	}

	if spriteY+int(h) >= spritesRows {
		h = float64(spritesRows - spriteY)
	}

	spriteSheetOffset := spriteY*ssWidth*SpriteHeight + spriteX*SpriteWidth

	width := int(SpriteWidth * w)
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

	height := int(SpriteHeight * h)
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

	for i := 0; i < height; i++ {
		spriteSheetLine := SpriteSheetData[spriteSheetOffset : spriteSheetOffset+width]
		if flipX {
			for j := 0; j < len(spriteSheetLine); j++ {
				ScreenData[screenOffset+j] = spriteSheetLine[len(spriteSheetLine)-1-j]
			}
		} else {
			copy(ScreenData[screenOffset:screenOffset+width], spriteSheetLine)
		}
		screenOffset += scrWidth
		spriteSheetOffset += spriteSheetStep
	}
}
