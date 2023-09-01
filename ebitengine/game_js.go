// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitengine

import (
	"bytes"
	_ "embed"
	"image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed "play.png"
var playButton []byte

var playButtonImage *ebiten.Image

func init() {
	img, err := png.Decode(bytes.NewReader(playButton))
	if err != nil {
		panic("decoding play.png failed: " + err.Error())
	}
	playButtonImage = ebiten.NewImageFromImage(img)
}

// in web browser user action is needed to initialize audio - such as clicking the mouse or hitting the keyboard.
// To inform the user that his action is needed drawNotReady draws a play button.
func (e *game) drawNotReady(screen *ebiten.Image) {
	screenSize := screen.Bounds().Max
	screenWidth := screenSize.X
	screenHeight := screenSize.Y

	imageSize := playButtonImage.Bounds()
	imageWidth := imageSize.Max.X
	imageHeight := imageSize.Max.Y

	m := ebiten.GeoM{}
	m.Translate(float64(screenWidth-imageWidth)/2.0, float64(screenHeight-imageHeight)/2.0)

	screen.DrawImage(playButtonImage, &ebiten.DrawImageOptions{GeoM: m})
}
