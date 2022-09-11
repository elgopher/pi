// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package key_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi/key"
)

func TestButton_String(t *testing.T) {
	t.Run("should convert key to string", func(t *testing.T) {
		tests := map[key.Button]string{
			key.Button('!'): "!",
			key.Digit0:      "0",
			key.A:           "A",
			key.Button('a'): "a",
			key.Button(31):  "?",
			key.Button(127): "?",
			key.Space:       "Space",
			key.Shift:       "Shift",
			key.Ctrl:        "Ctrl",
			key.Alt:         "Alt",
			key.Cap:         "Cap",
			key.Back:        "Back",
			key.Tab:         "Tab",
			key.Enter:       "Enter",
			key.Left:        "Left",
			key.Up:          "Up",
			key.Right:       "Right",
			key.Down:        "Down",
			key.Esc:         "Esc",
			key.F1:          "F1",
			key.F2:          "F2",
			key.F12:         "F12",
		}
		for btn, str := range tests {
			assert.Equal(t, btn.String(), str)
		}
	})
}
