// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/elgopher/pi"
)

func CopyCanvasToEbitenImage(canvas pi.Canvas, dst *ebiten.Image) {
	pixels := canvas.Data()
	if buffer == nil || len(buffer)/4 < len(pixels) {
		buffer = make([]byte, len(pixels)*4)
	}

	buff := buffer[0 : len(pixels)*4]

	offset := 0
	for _, col := range pixels {
		col &= pi.MaxColors - 1
		rgba := pi.Palette[pi.PaletteMapping[col]&(pi.MaxColors-1)]
		buff[offset] = byte(rgba >> 16)
		buff[offset+1] = byte(rgba >> 8)
		buff[offset+2] = byte(rgba)
		buff[offset+3] = 0xFF
		offset += 4
	}

	dst.WritePixels(buff)
}

var buffer []byte
