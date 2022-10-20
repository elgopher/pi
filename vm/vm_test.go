// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package vm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi/vm"
)

func TestRGB_String(t *testing.T) {
	tests := map[string]vm.RGB{
		"#000000": {},
		"#010203": {1, 2, 3},
		"#102030": {0x10, 0x20, 0x30},
	}

	for expected, rgb := range tests {
		assert.Equal(t, expected, rgb.String())
	}
}
