package audio

import (
	"math"

	"github.com/elgopher/pi/audio/internal"
)

func silence(float64) float64 {
	return 0
}

func oscillatorFunc(instrument Instrument) func(float64) float64 {
	switch instrument {
	case InstrumentTriangle:
		return internal.Triangle
	case InstrumentTiltedSaw:
		return internal.TiltedSaw
	case InstrumentSaw:
		return internal.Saw
	case InstrumentSquare:
		return internal.Square
	case InstrumentPulse:
		return internal.Pulse
	case InstrumentOrgan:
		return internal.Organ
	default:
		return silence
	}
}

// pitchToFreq returns frequency in Hz.
func pitchToFreq(pitch Pitch) float64 {
	diff := int(pitch) - int(PitchA2)
	multiplier := math.Pow(2, float64(diff)/12)
	return 440 * multiplier
}

func singleNoteSamples(speed byte) int {
	if speed == 0 {
		speed = 1
	}

	return int(speed) * 183
}
