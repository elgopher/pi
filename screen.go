package pi

import "github.com/elgopher/pi/image"

var (
	Color byte = 6 // Color is a currently used color in draw state.

	// Palette has all colors available in the game. Up to 256. Palette is taken from sprite sheet PNG file with indexed color mode.
	//
	// Can be freely read and updated. Changes will be visible immediately.
	Palette [256]image.RGB

	// ScreenData contains pixel colors for screen visible by the player. Each pixel is one byte. It is initialized during pi.Boot.
	//
	// Can be freely read and updated. Useful when you want to use your own functions for pixel manipulation. Pi will panic if you try to change the length of the slice.
	ScreenData []byte

	scrWidth, scrHeight int
	lineOfScreenWidth   []byte
	zeroScreenData      []byte
	clippingRegion      rect
)

func Cls() {
	copy(ScreenData, zeroScreenData)
}

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

func ClipReset() (prevX, prevY, prevW, prevH int) {
	return Clip(0, 0, scrWidth, scrHeight)
}

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

//func ClipPrev(x, y, w, h int) {}

type pos struct {
	x, y int
}

var camera pos

func Camera(x, y int) (prevX, prevY int) {
	prev := camera
	camera.x = x
	camera.y = y
	return prev.x, prev.y
}

func CameraReset() (prevX, prevY int) {
	return Camera(0, 0)
}

func Spr(n, x, y int) {
	SprSize(n, x, y, 1.0, 1.0)
}

func SprSize(n, x, y int, w, h float64) {
	SprSizeFlip(n, x, y, w, h, false, false)
}

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

	for i := 0; i < height; i++ {
		copy(ScreenData[screenOffset:screenOffset+width], SpriteSheetData[spriteSheetOffset:spriteSheetOffset+width])
		screenOffset += scrWidth
		spriteSheetOffset += ssWidth
	}
}