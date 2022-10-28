// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package rgb

import (
	"github.com/elgopher/pi/image"
)

func BrightnessIndex(rgb image.RGB) byte {
	return byte((299*int(rgb.R) + 587*int(rgb.G) + 114*int(rgb.B)) / 1000)
}

func BrightnessDelta(rgb1, rgb2 image.RGB) byte {
	i1 := BrightnessIndex(rgb1)
	i2 := BrightnessIndex(rgb2)
	if i1 > i2 {
		return i1 - i2
	}

	return i2 - i1
}
