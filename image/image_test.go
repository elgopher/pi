// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package image_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi/image"
)

func TestRGB_String(t *testing.T) {
	tests := map[string]image.RGB{
		"#000000": {},
		"#FFFFFF": {0xFF, 0xFF, 0xFF},
		"#012345": {0x01, 0x23, 0x45},
		"#6789AB": {0x67, 0x89, 0xAB},
		"#CDEF01": {0xCD, 0xEF, 0x01},
	}

	for expected, rgb := range tests {
		assert.Equal(t, expected, rgb.String())
	}
}
