// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package icons

import (
	_ "embed"
	"fmt"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/font"
)

const (
	Pointer      = 0
	MeasureTool  = 1
	SetTool      = 2
	LineTool     = 3
	RectTool     = 4
	RectFillTool = 5
	CircTool     = 6
	CircFillTool = 7
	SprTool      = 8
)

//go:embed icons.png
var iconsPng []byte

var icons = pi.Font{
	Width:        4,
	SpecialWidth: 8,
	Height:       4,
}

func init() {
	var err error
	icons.Data, err = font.Load(iconsPng)
	if err != nil {
		panic(fmt.Sprintf("problem loading devtools icons %s", err))
	}
}

func Draw(x, y int, color byte, icon ...byte) {
	for _, i := range icon {
		text := string(rune(i))
		x = icons.Print(text, x, y, color)
	}
}
