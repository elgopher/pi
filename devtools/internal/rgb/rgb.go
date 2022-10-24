package rgb

import "github.com/elgopher/pi/mem"

func BrightnessIndex(rgb mem.RGB) byte {
	return byte((299*int(rgb.R) + 587*int(rgb.G) + 114*int(rgb.B)) / 1000)
}

func BrightnessDelta(rgb1, rgb2 mem.RGB) byte {
	i1 := BrightnessIndex(rgb1)
	i2 := BrightnessIndex(rgb2)
	if i1 > i2 {
		return i1 - i2
	}

	return i2 - i1
}
