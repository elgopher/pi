package pi

const (
	SpriteWidth, SpriteHeight = 8, 8
)

var (
	// SpriteSheetData contains pixel colors for the entire sprite sheet. Each pixel is one byte. It is initialized during pi.Boot.
	//
	// Can be freely read and updated. Useful when you want to use your own functions for pixel manipulation.
	// Pi will panic if you try to change the length of the slice.
	SpriteSheetData []byte

	ssWidth, ssHeight int
	numberOfSprites   int
	spritesInLine     int
	spritesRows       int
)

// Sset does not change the global Color
func Sset(x, y int, color byte) {
	if x < 0 {
		return
	}
	if y < 0 {
		return
	}
	if x >= ssWidth {
		return
	}
	if y >= ssHeight {
		return
	}

	SpriteSheetData[y*ssWidth+x] = color
}

func Sget(x, y int) byte {
	if x < 0 {
		return 0
	}
	if y < 0 {
		return 0
	}
	if x >= ssWidth {
		return 0
	}
	if y >= ssHeight {
		return 0
	}

	return SpriteSheetData[y*ssWidth+x]
}
