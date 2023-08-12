// (c) 2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build !js

package help_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi/devtools/internal/help"
	"github.com/elgopher/pi/devtools/internal/test"
)

func TestPrintHelp(t *testing.T) {
	t.Run("should return error when trying to print help for not imported packages", func(t *testing.T) {
		topics := []string{
			"io", "io.Writer",
		}
		for _, topic := range topics {
			t.Run(topic, func(t *testing.T) {
				// when
				err := help.PrintHelp(topic)
				// then
				assert.ErrorIs(t, err, help.NotFound)
			})
		}
	})

	t.Run("should return error when trying to print help for non-existent symbol", func(t *testing.T) {
		err := help.PrintHelp("pi.NonExistent")
		assert.ErrorIs(t, err, help.NotFound)
	})

	t.Run("should print help for", func(t *testing.T) {
		tests := map[string]struct {
			topic    string
			expected string
		}{
			"package": {
				topic:    "pi",
				expected: `Package pi`,
			},
			"function": {
				topic:    "pi.Spr",
				expected: `func Spr(n, x, y int)`,
			},
			"struct": {
				topic:    "pi.PixMap",
				expected: `type PixMap struct {`,
			},
		}
		for testName, testCase := range tests {
			t.Run(testName, func(t *testing.T) {
				swapper := test.SwapStdout(t)
				// when
				err := help.PrintHelp(testCase.topic)
				// then
				swapper.BringStdoutBack()
				assert.NoError(t, err)
				output := swapper.ReadOutput(t)
				assert.Contains(t, output, testCase.expected)
			})
		}
	})

	t.Run("should show help for image.Image from github.com/elgopher/pi package, not from stdlib", func(t *testing.T) {
		topics := []string{
			"image", "image.Image",
		}
		for _, topic := range topics {
			t.Run(topic, func(t *testing.T) {
				swapper := test.SwapStdout(t)
				// when
				err := help.PrintHelp("image.Image")
				// then
				swapper.BringStdoutBack()
				assert.NoError(t, err)
				output := swapper.ReadOutput(t)
				assert.Contains(t, output, `// import "github.com/elgopher/pi/image"`)
			})
		}
	})

	t.Run("should show detailed help for pi.Button", func(t *testing.T) {
		tests := map[string]string{
			"pi.Button":      "Keyboard mappings",
			"pi.MouseButton": "MouseRight  MouseButton = 2",
			"key.Button":     "func (b Button) String() string",
		}
		for topic, expected := range tests {
			t.Run(topic, func(t *testing.T) {
				swapper := test.SwapStdout(t)
				// when
				err := help.PrintHelp(topic)
				// then
				swapper.BringStdoutBack()
				assert.NoError(t, err)
				output := swapper.ReadOutput(t)
				assert.Contains(t, output, expected)
			})
		}
	})
}