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
	ColorTables[0] = defaultColorTable()
	ColorTables[1] = identityColorTable()
}

func Pal(fromColor, toColor Color) {
	ColorTables[0][fromColor] = useOneColorNoMatterWhatIsTheTargetColor(toColor)
}

func Palt(color Color, t bool) {
	if t {
		ColorTables[0][color] = transparent()
	} else {
		ColorTables[0][color] = useOneColorNoMatterWhatIsTheTargetColor(color)
	}
}

func defaultColorTable() (p ColorTable) {
	p[0] = transparent()
	for sourceColor := Color(1); sourceColor < MaxColors; sourceColor++ {
		p[sourceColor] = useOneColorNoMatterWhatIsTheTargetColor(sourceColor)
	}
	return p
}

// drawing has no effect
func identityColorTable() (p [MaxColors][MaxColors]Color) {
	for sourceColor := Color(0); sourceColor < MaxColors; sourceColor++ {
		for targetColor := Color(0); targetColor < MaxColors; targetColor++ {
			p[sourceColor][targetColor] = targetColor
		}
	}
	return p
}

func transparent() (colors [MaxColors]Color) {
	for i := Color(0); i < MaxColors; i++ {
		colors[i] = i
	}
	return
}

func useOneColorNoMatterWhatIsTheTargetColor(c Color) (colors [MaxColors]Color) {
	for targetColor := 0; targetColor < MaxColors; targetColor++ {
		colors[targetColor] = c
	}
	return
}
