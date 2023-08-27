package audio

import "github.com/elgopher/pi/audio/internal"

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
