// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import "github.com/elgopher/pi/mem"

const (
	defaultSpriteSheetWidth  = 128
	defaultSpriteSheetHeight = 128
	defaultScreenWidth       = 128
	defaultScreenHeight      = 128
)

var (
	//nolint:govet
	defaultPalette = [256]mem.RGB{
		{0, 0, 0},          // 0 - black
		{0x1D, 0x2B, 0x53}, // 1 - dark blue
		{0x7E, 0x25, 0x53}, // 2 - dark purple
		{0x00, 0x87, 0x51}, // 3 - dark green
		{0xAB, 0x52, 0x36}, // 4 - brown
		{0x5F, 0x57, 0x4F}, // 5 - dark gray
		{0xC2, 0xC3, 0xC7}, // 6 - light gray
		{0xff, 0xf1, 0xe8}, // 7 - white
		{0xFF, 0x00, 0x4D}, // 8 - red
		{0xFF, 0xA3, 0x00}, // 9 - orange
		{0xFF, 0xEC, 0x27}, // 10 - yellow
		{0x00, 0xE4, 0x36}, // 11 - green
		{0x29, 0xAD, 0xFF}, // 12 - blue
		{0x83, 0x76, 0x9C}, // 13 - indigo
		{0xFF, 0x77, 0xA8}, // 14 - pink
		{0xFF, 0xCC, 0xAA}, // 15 - peach
	}
)
