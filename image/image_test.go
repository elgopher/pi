package image_test

import (
	"testing"

	"github.com/elgopher/pi/image"
	"github.com/stretchr/testify/assert"
)

func TestRGB_String(t *testing.T) {
	tests := map[string]image.RGB{
		"#000000": {},
		"#010203": {1, 2, 3},
		"#102030": {0x10, 0x20, 0x30},
	}

	for expected, rgb := range tests {
		assert.Equal(t, expected, rgb.String())
	}
}
