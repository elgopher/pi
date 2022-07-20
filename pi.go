package pi

import (
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"time"
)

const (
	defaultSpriteSheetWidth  = 128
	defaultSpriteSheetHeight = 128
	defaultScreenWidth       = 128
	defaultScreenHeight      = 128
)

// User parameters. Will be used on Boot
var (
	Update = func() {} // Update is executed every frame
	Draw   = func() {} // Draw is executed every frame, if he can.

	Resources fs.ReadFileFS // Resources contains files like sprite-sheet.png

	// SpriteSheetWidth will be used if sprite-sheet.png was not found.
	SpriteSheetWidth = defaultSpriteSheetWidth
	// SpriteSheetHeight will be used if sprite-sheet.png was not found.
	SpriteSheetHeight = defaultSpriteSheetHeight

	ScreenWidth  = defaultScreenWidth
	ScreenHeight = defaultScreenHeight
)

// Run opens the window and run the game. It returns error when something terrible happened during initialization
// or the game panicked during execution.
func Run() error {
	if err := Boot(); err != nil {
		return fmt.Errorf("booting game failed: %w", err)
	}

	timeStarted = time.Now()

	return run()
}

// RunOrPanic opens the window and run he game. It panics when something terrible happened. Useful for quick and dirty
// testing.
func RunOrPanic() {
	if err := Run(); err != nil {
		panic(fmt.Sprintf("Something terrible happened! Pi cannot be run: %v\n", err))
	}
}

// Reset resets all user parameters to default values. Useful in unit tests.
func Reset() {
	Update = nil
	Draw = nil
	Resources = nil
	SpriteSheetWidth = defaultSpriteSheetWidth
	SpriteSheetHeight = defaultSpriteSheetHeight
	ScreenWidth = defaultScreenWidth
	ScreenHeight = defaultScreenHeight
	Color = 6
}

func Boot() error {
	if SpriteSheetWidth%8 != 0 || SpriteSheetWidth == 0 {
		return fmt.Errorf("sprite sheet width %d is not a multiplcation of 8", SpriteSheetWidth)
	}
	if SpriteSheetHeight%8 != 0 || SpriteSheetHeight == 0 {
		return fmt.Errorf("sprite sheet height %d is not a multiplcation of 8", SpriteSheetHeight)
	}
	ssWidth = SpriteSheetWidth
	ssHeight = SpriteSheetHeight
	numberOfSprites = (ssWidth * ssHeight) / (SpriteWidth * SpriteHeight)

	if Resources == nil {
		Resources = embed.FS{}
	}

	if err := loadResources(Resources); err != nil {
		return err
	}

	spritesInLine = ssWidth / SpriteWidth
	spritesRows = ssHeight / SpriteHeight

	scrWidth = ScreenWidth
	scrHeight = ScreenHeight
	screenSize := scrWidth * scrHeight
	ScreenData = make([]byte, screenSize)
	zeroScreenData = make([]byte, screenSize)
	lineOfScreenWidth = make([]byte, scrWidth)

	Clip(0, 0, scrWidth, scrHeight)
	Camera(0, 0)

	return nil
}

func BootOrPanic() {
	if err := Boot(); err != nil {
		panic("init failed " + err.Error())
	}
}
