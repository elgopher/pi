// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"fmt"
	"strings"
)

// PaletteMapping defines the color mapping for display.
//
// Colors are remapped at the end of each frame, just before rendering
// the contents of the Screen to the game window.
//
// The array index is the original color, and the value is the mapped color.
var PaletteMapping PaletteMap = notRemappedPalette

// PaletteMap defines the color mapping for display.
//
// The array index is the original color, and the value is the mapped color.
type PaletteMap [MaxColors]Color

func (p PaletteMap) String() string {
	var s strings.Builder
	s.WriteString("{")
	for i := 0; i < len(p); i++ {
		s.WriteString(fmt.Sprintf("%d:%v, ", i, p[i]))
	}
	s.WriteString("}")
	return s.String()
}

func ResetPaletteMapping() {
	PaletteMapping = notRemappedPalette
}

var notRemappedPalette = func() PaletteMap {
	var m PaletteMap
	for i := 0; i < MaxColors; i++ {
		m[i] = Color(i)
	}
	return m
}()
