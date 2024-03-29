// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi"
)

const delta = 0.000000000000001

func TestSin(t *testing.T) {
	tests := map[float64]float64{
		-2:   0,
		-1.5: 0,
		-1:   0,
		-0.9: -0.5877852522924732,
		-0.8: -0.9510565162951536,
		-0.7: -0.9510565162951535,
		-0.6: -0.587785252292473,
		-0.5: 0,
		-0.4: 0.5877852522924732,
		-0.3: 0.9510565162951536,
		-0.2: 0.9510565162951535,
		-0.1: 0.5877852522924731,
		0:    0,
		0.1:  -0.5877852522924731,
		0.2:  -0.9510565162951535,
		0.3:  -0.9510565162951536,
		0.4:  -0.5877852522924732,
		0.5:  0,
		0.6:  0.587785252292473,
		0.7:  0.9510565162951535,
		0.8:  0.9510565162951536,
		0.9:  0.5877852522924732,
		1:    0,
		1.5:  0,
		2:    0,
	}

	for angle, expected := range tests {
		name := fmt.Sprintf("%f", angle)
		t.Run(name, func(t *testing.T) {
			actual := pi.Sin(angle)
			assert.InDelta(t, expected, actual, delta)
		})
	}
}

func TestCos(t *testing.T) {
	tests := map[float64]float64{
		-2:   1,
		-1.5: -1,
		-1:   1,
		-0.9: 0.8090169943749473,
		-0.8: 0.30901699437494723,
		-0.7: -0.30901699437494756,
		-0.6: -0.8090169943749473,
		-0.5: -1,
		-0.4: -0.8090169943749475,
		-0.3: -0.30901699437494734,
		-0.2: 0.30901699437494734,
		-0.1: 0.8090169943749473,
		0:    1,
		0.1:  0.8090169943749473,
		0.2:  0.30901699437494723,
		0.3:  -0.30901699437494723,
		0.4:  -0.8090169943749473,
		0.5:  -1,
		0.6:  -0.8090169943749473,
		0.7:  -0.30901699437494723,
		0.8:  0.30901699437494723,
		0.9:  0.8090169943749473,
		1:    1,
		1.5:  -1,
		2:    1,
	}

	for angle, expected := range tests {
		name := fmt.Sprintf("%f", angle)
		t.Run(name, func(t *testing.T) {
			actual := pi.Cos(angle)
			assert.InDelta(t, expected, actual, delta)
		})
	}
}

func TestAtan2(t *testing.T) {
	type params struct{ dx, dy float64 }
	tests := map[params]float64{
		{0, 0}:         0.75, // TODO In Pico-8 this is 0.25
		{0.001, 0.001}: 0.875,
		{0, 1}:         0.75,
		{1, 0}:         0,
		{1, 1}:         0.875,
		{-1, -1}:       0.375,
		{-1, 1}:        0.625,
		{1, -1}:        0.125,
		{1, 2}:         0.8237918088252166,
		{1, -0.3}:      0.04638678953887121,
	}
	for p, expected := range tests {
		name := fmt.Sprintf("%+v", p)
		t.Run(name, func(t *testing.T) {
			actual := pi.Atan2(p.dx, p.dy)
			assert.InDelta(t, expected, actual, delta)
		})
	}
}

func TestMinInt(t *testing.T) {
	assert.Equal(t, 0, pi.MinInt(0, 0))
	assert.Equal(t, 1, pi.MinInt(1, 2))
	assert.Equal(t, 1, pi.MinInt(1, 1))
	assert.Equal(t, 1, pi.MinInt(2, 1))
	assert.Equal(t, -2, pi.MinInt(-1, -2))
	assert.Equal(t, -2, pi.MinInt(-2, 2))
}

func TestMaxInt(t *testing.T) {
	assert.Equal(t, 0, pi.MaxInt(0, 0))
	assert.Equal(t, 2, pi.MaxInt(2, 1))
	assert.Equal(t, 1, pi.MaxInt(1, 1))
	assert.Equal(t, 2, pi.MaxInt(1, 2))
	assert.Equal(t, -1, pi.MaxInt(-1, -2))
	assert.Equal(t, 2, pi.MaxInt(-2, 2))
}

func TestMidInt(t *testing.T) {
	assert.Equal(t, 0, pi.MidInt(0, 0, 0))
	assert.Equal(t, 1, pi.MidInt(0, 1, 2))
	assert.Equal(t, 1, pi.MidInt(2, 1, 0))
	assert.Equal(t, 1, pi.MidInt(1, 0, 2))
	assert.Equal(t, 1, pi.MidInt(1, 2, 0))
	assert.Equal(t, 1, pi.MidInt(2, 0, 1))
	assert.Equal(t, 1, pi.MidInt(0, 2, 1))
	assert.Equal(t, -1, pi.MidInt(0, -1, -2))
}

func TestMid(t *testing.T) {
	assert.Equal(t, 0.0, pi.Mid(0, 0, 0))
	assert.Equal(t, 1.0, pi.Mid(0, 1, 2))
	assert.Equal(t, 1.0, pi.Mid(2, 1, 0))
	assert.Equal(t, 1.0, pi.Mid(1, 0, 2))
	assert.Equal(t, 1.0, pi.Mid(1, 2, 0))
	assert.Equal(t, 1.0, pi.Mid(2, 0, 1))
	assert.Equal(t, 1.0, pi.Mid(0, 2, 1))
	assert.Equal(t, -1.0, pi.Mid(0, -1, -2))

	assertNaN(t, pi.Mid(math.NaN(), math.NaN(), math.NaN()))
	assertInf(t, pi.Mid(math.Inf(1), math.Inf(1), math.Inf(1)), 1)
	assert.Equal(t, 1.0, pi.Mid(1.0, math.NaN(), 2.0)) // NaNs always go to the beginning
	assertNaN(t, pi.Mid(math.NaN(), math.NaN(), 1.0))
	assertNaN(t, pi.Mid(1.0, math.NaN(), math.NaN()))
	assert.Equal(t, 1.0, pi.Mid(1.0, 2.0, math.NaN()))
	assert.Equal(t, 1.0, pi.Mid(math.Inf(1), math.Inf(-1), 1.0))
}

func assertNaN(t *testing.T, v float64) {
	assert.True(t, math.IsNaN(v))
}

func assertInf(t *testing.T, v float64, sign int) {
	assert.True(t, math.IsInf(v, sign))
}
