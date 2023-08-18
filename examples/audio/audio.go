package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/ebitengine" // engine capable of rendering the game on multiple operating systems
	"github.com/elgopher/pi/internal/audio"
)

const tau = 2 * math.Pi

func main() {
	var osc audio.Oscillator
	osc.SetDuration(time.Hour)
	freq := uint16(415)
	osc.SetFrequency(freq)

	forms := []audio.WaveForm{
		{
			Name: "triangle",
			F:    triangle,
		},
		{
			Name: "square",
			F:    square,
		},
		{
			Name: "noise",
			F:    noise,
		},
		{
			Name: "organ",
			F:    organ,
		},
		{
			Name: "tilted saw",
			F:    tiltedSaw,
		},
		{
			Name: "sin",
			F:    math.Sin,
		},
		{
			Name: "sawtooth",
			F:    sawtooth,
		},
		{
			Name: "pulse",
			F:    pulse,
		},

		{
			Name: "upwSawtooth",
			F:    upwSawtooth,
		},
		{
			Name: "downSawtooth",
			F:    downSawtooth,
		},
	}

	osc.SetWaveForm(forms[0])

	ebitengine.AudioStream = &osc

	currentForm := 0

	pi.Update = func() {
		d := uint16(20)
		if pi.Btn(pi.O) {
			d = 200
		}
		if pi.Btnp(pi.Up) {
			freq += d
			osc.SetFrequency(freq)
			fmt.Println(freq)
		}
		if pi.Btnp(pi.Down) {
			freq -= d
			osc.SetFrequency(freq)
			fmt.Println(freq)
		}
		if pi.Btnp(pi.Left) {
			currentForm -= 1
			if currentForm < 0 {
				currentForm = len(forms) - 1
			}
			f := forms[currentForm]
			osc.SetWaveForm(f)
			fmt.Println(f.Name)
		}
		if pi.Btnp(pi.Right) {
			currentForm += 1
			if currentForm >= len(forms) {
				currentForm = 0
			}
			f := forms[currentForm]
			osc.SetWaveForm(f)
			fmt.Println(f.Name)
		}
		if pi.Btn(pi.X) {
			tiltFactor += 0.05
			fmt.Println(tiltFactor)
		}

	}
	//devtools.Export("osc", &osc)
	//devtools.ExportType[audio.WaveForm]()
	// Run game with devtools (Hit F12 to show screen inspector)
	//devtools.MustRun(ebitengine.Run)
	ebitengine.MustRun()
}

func sawtooth(t float64) float64 {
	return ((t / tau) * 2) - 1
}

func square(time float64) float64 {
	if time < math.Pi {
		return 1
	} else {
		return -1
	}
}

func pulse(phase float64) float64 {
	if phase < math.Pi+math.Pi*2/3 {
		return 1
	} else {
		return -1
	}
}

func triangle(phase float64) float64 {
	val := 2.0*(phase*(1.0/tau)) - 1.0
	if val < 0.0 {
		val = -val
	}
	val = 2.0 * (val - 0.5)
	return val
}

func upwSawtooth(phase float64) float64 {
	val := 2.0*(phase*(1.0/tau)) - 1.0
	return val
}

func downSawtooth(phase float64) float64 {
	val := 1.0 - 2.0*(phase*(1.0/tau))
	return val
}

// TO W OGOLE NIE DZIALA
func organ(phase float64) float64 {
	v := triangle(phase)
	if phase > math.Pi {
		v = v / 2
	}

	return v
}

func noise(phase float64) float64 {
	return math.Sin(phase) * ((rand.Float64() * 2) - 1)
}

var tiltFactor = 0.1

// pochylona piła
func tiltedSaw(f float64) float64 {
	return (1-tiltFactor)*downSawtooth(f) + tiltFactor*triangle(f)
}
