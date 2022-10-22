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
	Pointer = "\u0000"
)

//go:embed icons.png
var iconsPng []byte

var icons = pi.Font{
	Width:        4,
	WidthSpecial: 8,
	Height:       8,
}

func init() {
	if err := font.Load(iconsPng, icons.Data[:]); err != nil {
		panic(fmt.Sprintf("problem loading devtools icons %s", err))
	}
}

func Draw(icon string, x, y int, color byte) {
	icons.Print(icon, x, y, color)
}
