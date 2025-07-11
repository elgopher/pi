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

func Pal(fromColor, toColor Color) {
	ColorTables[0][fromColor] = opaqueColorTable[toColor]
}

func Palt(color Color, t bool) {
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
