// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"fmt"
	"strings"
)

// ColorTable defines the color mapping rules used during drawing.
//
// The first index is the draw (source) color index, and the second is the
// target (background) color index.
//
// For example:
//
//	colorTable[7][0] = 6
//
// means drawing color 7 over color 0 will result in color 6.
//
// See more about color tables:
//
//	https://www.lexaloffle.com/bbs/?tid=149249
type ColorTable [MaxColors][MaxColors]Color

func (p ColorTable) String() string {
	var s strings.Builder
	for i := 0; i < len(p); i++ {
		s.WriteString(fmt.Sprintf("%d :%v\n", i, p[i]))
	}
	return s.String()
}

// ColorTables defines 4 different ColorTable entries,
// selected based on bits 6 and 7 in the Color value.
//
// When drawing a draw (source) color over a target (background) color,
// the following operation is performed:
//
//	(source | target) >> 6
//
// The result of this operation determines the index of the ColorTable to use.
var ColorTables [4]ColorTable

func init() {
	ResetColorTables()
}

func ResetColorTables() {
	ColorTables[0] = opaqueColorTable
	ColorTables[0][0] = transparentColor
	ColorTables[1] = identityColorTable
}

// RemapColor changes how the color from is rendered by replacing it with the color to.
//
// This affects all future drawing operations (sprites, shapes, text, etc.).
// It does not modify the original image or sprite data — only how colors appear on screen.
//
// For example, calling:
//
//	RemapColor(1, 15)
//
// causes all pixels with color index 1 to be drawn using color 15 instead.
//
// This function updates the color tables. To reset the changes, use ResetColorTables.
func RemapColor(from, to Color) {
	ColorTables[0][from] = opaqueColorTable[to]
}

// SetTransparency sets whether the given color is treated as transparent.
//
// When transparency is enabled for a color, pixels using that color will not be drawn.
// This affects all future drawing operations (sprites, shapes, text, etc.).
// It does not modify the original image or sprite data — only how colors appear on screen.
//
// For example, calling:
//
//	SetTransparency(0, true)
//
// makes color 0 transparent, meaning all pixels with color index 0 will be skipped during drawing.
//
// To disable transparency for a color, pass false as the second argument.
//
// This function updates the color tables. To reset the changes, use ResetColorTables.
func SetTransparency(color Color, t bool) {
	if t {
		ColorTables[0][color] = transparentColor
	} else {
		ColorTables[0][color] = opaqueColorTable[color]
	}
}

var opaqueColorTable = func() (table [MaxColors][MaxColors]Color) {
	for draw := Color(0); draw < MaxColors; draw++ {
		for target := 0; target < MaxColors; target++ {
			table[draw][target] = draw
		}
	}
	return
}()

var transparentColor = func() (colors [MaxColors]Color) {
	for target := Color(0); target < MaxColors; target++ {
		colors[target] = target
	}
	return
}()

// drawing has no effect
var identityColorTable = func() (table [MaxColors][MaxColors]Color) {
	for draw := Color(0); draw < MaxColors; draw++ {
		for target := Color(0); target < MaxColors; target++ {
			table[draw][target] = target
		}
	}
	return table
}()
