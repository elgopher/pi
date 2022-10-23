package rgb

import "github.com/elgopher/pi/vm"

func BrightnessIndex(rgb vm.RGB) byte {
	return byte((299*int(rgb.R) + 587*int(rgb.G) + 114*int(rgb.B)) / 1000)
}

func BrightnessDelta(rgb1, rgb2 vm.RGB) byte {
	i1 := BrightnessIndex(rgb1)
	i2 := BrightnessIndex(rgb2)
	if i1 > i2 {
		return i1 - i2
	}

	return i2 - i1
}
