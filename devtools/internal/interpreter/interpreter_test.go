// (c) 2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build !js

package interpreter_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elgopher/pi/devtools/internal/interpreter"
	"github.com/elgopher/pi/devtools/internal/test"
)

func TestEval(t *testing.T) {
	t.Run("should evaluate command", func(t *testing.T) {
		tests := map[string]struct {
			code           string
			expectedOutput string
			expectedResult interpreter.EvalResult
		}{
			"simple expression": {
				code:           "1",
				expectedOutput: "int: 1\n",
				expectedResult: interpreter.GoCodeExecuted,
			},
			"run statement": {
				code:           "pi.Spr(0,0,0)",
				expectedResult: interpreter.GoCodeExecuted,
			},
			"var declaration": {
				code:           "var a string",
				expectedOutput: "string: \n",
				expectedResult: interpreter.GoCodeExecuted,
			},
			"expression returning a struct": {
				code:           "struct{}{}",
				expectedOutput: "struct {}: {}\n",
				expectedResult: interpreter.GoCodeExecuted,
			},
			"line with single curly bracket": {
				code:           "{",
				expectedResult: interpreter.Continued,
			},
			"two lines with curly bracket": {
				code:           "{\n{",
				expectedResult: interpreter.Continued,
			},
			"string literal not terminated": {
				code:           "` ",
				expectedResult: interpreter.Continued,
			},
			"undo": {
				code:           "undo",
				expectedResult: interpreter.Undoed,
			},
			"undo with space in the beginning": {
				code:           " undo",
				expectedResult: interpreter.Undoed,
			},
			"undo with space in the end": {
				code:           "undo ",
				expectedResult: interpreter.Undoed,
			},
			"pause": {
				code:           "pause",
				expectedResult: interpreter.Paused,
			},
			"pause with space in the beginning": {
				code:           " pause",
				expectedResult: interpreter.Paused,
			},
			"pause with space in the end": {
				code:           "pause ",
				expectedResult: interpreter.Paused,
			},
			"resume": {
				code:           "resume",
				expectedResult: interpreter.Resumed,
			},
			"resume with space in the beginning": {
				code:           " resume",
				expectedResult: interpreter.Resumed,
			},
			"resume with space in the end": {
				code:           "resume ",
				expectedResult: interpreter.Resumed,
			},
			"next": {
				code:           "next",
				expectedResult: interpreter.NextFrameRequested,
			},
			"next with space in the beginning": {
				code:           " next",
				expectedResult: interpreter.NextFrameRequested,
			},
			"next with space in the end": {
				code:           "next ",
				expectedResult: interpreter.NextFrameRequested,
			},
			"n": {
				code:           "n",
				expectedResult: interpreter.NextFrameRequested,
			},
		}
		for name, testCase := range tests {
			t.Run(name, func(t *testing.T) {
				swapper := test.SwapStdout(t)
				// when
				result, err := newInterpreterInstance(t).Eval(testCase.code)
				// then
				swapper.BringStdoutBack()
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedResult, result)
				assert.Equal(t, testCase.expectedOutput, swapper.ReadOutput(t))
			})
		}
	})

	t.Run("should return error on compilation error", func(t *testing.T) {
		swapper := test.SwapStdout(t)
		// when
		result, err := newInterpreterInstance(t).Eval("1 === 1")
		// then
		swapper.BringStdoutBack()
		assert.Error(t, err)
		assert.Equal(t, interpreter.GoCodeExecuted, result)
		assert.Equal(t, "", swapper.ReadOutput(t))
	})

	t.Run("should return error when script panics", func(t *testing.T) {
		swapper := test.SwapStdout(t)
		// when
		result, err := newInterpreterInstance(t).Eval(`panic("panic")`)
		// then
		swapper.BringStdoutBack()
		assert.Error(t, err)
		assert.Equal(t, interpreter.GoCodeExecuted, result)
		assert.Equal(t, "", swapper.ReadOutput(t))
	})

	t.Run("should return error when Yaegi panics", func(t *testing.T) {
		swapper := test.SwapStdout(t)
		// when
		result, err := newInterpreterInstance(t).Eval(`fmt=1`)
		// then
		swapper.BringStdoutBack()
		assert.Error(t, err)
		assert.Equal(t, interpreter.GoCodeExecuted, result)
		assert.Equal(t, "", swapper.ReadOutput(t))
	})

	t.Run("script should read exported variable", func(t *testing.T) {
		instance := newInterpreterInstance(t)

		err := instance.Export("variable", 10)
		require.NoError(t, err)

		swapper := test.SwapStdout(t)
		// when
		_, err = instance.Eval("variable")
		// then
		swapper.BringStdoutBack()
		require.NoError(t, err)
		assert.Equal(t, "int: 10\n", swapper.ReadOutput(t))
	})

	t.Run("script should update exported variable", func(t *testing.T) {
		instance := newInterpreterInstance(t)

		var variable int
		err := instance.Export("variable", &variable)
		require.NoError(t, err)

		swapper := test.SwapStdout(t)
		// when
		_, err = instance.Eval("*variable=1")
		// then
		swapper.BringStdoutBack()
		require.NoError(t, err)
		assert.Equal(t, 1, variable)
	})

	t.Run("should run exported func", func(t *testing.T) {
		instance := newInterpreterInstance(t)

		var functionExecuted bool
		err := instance.Export("fun", func() {
			functionExecuted = true
		})
		require.NoError(t, err)

		swapper := test.SwapStdout(t)
		// when
		_, err = instance.Eval("fun()")
		// then
		swapper.BringStdoutBack()
		require.NoError(t, err)
		assert.True(t, functionExecuted)
	})

	t.Run("should update exported func", func(t *testing.T) {
		instance := newInterpreterInstance(t)

		fun := func() int { return 0 }

		err := instance.Export("fun", &fun)
		require.NoError(t, err)

		swapper := test.SwapStdout(t)
		// when
		_, err = instance.Eval(`*fun = func() int { return 1 }`)
		// then
		swapper.BringStdoutBack()
		require.NoError(t, err)
		assert.Equal(t, 1, fun())
	})

	t.Run("should use exported type", func(t *testing.T) {
		instance := newInterpreterInstance(t)

		type customType struct{}
		err := interpreter.ExportType[customType](instance)
		require.NoError(t, err)

		swapper := test.SwapStdout(t)
		// when
		_, err = instance.Eval(`interpreter_test.customType{}`)
		// then
		swapper.BringStdoutBack()
		require.NoError(t, err)
		assert.Equal(t, "interpreter_test.customType: {}\n", swapper.ReadOutput(t))
	})

	t.Run("should use another exported type from the same package", func(t *testing.T) {
		instance := newInterpreterInstance(t)

		type anotherType struct{}
		type customType struct{}
		err := interpreter.ExportType[anotherType](instance)
		require.NoError(t, err)
		err = interpreter.ExportType[customType](instance)
		require.NoError(t, err)

		swapper := test.SwapStdout(t)
		// when
		_, err = instance.Eval(`interpreter_test.anotherType{}`)
		// then
		swapper.BringStdoutBack()
		require.NoError(t, err)
		assert.Equal(t, "interpreter_test.anotherType: {}\n", swapper.ReadOutput(t))
	})

	t.Run("should convert error message to exclude invalid character number", func(t *testing.T) {
		tests := map[string]struct {
			source             string
			expectedErrMessage string
		}{
			"undefined": {
				source:             "undefinedVar",
				expectedErrMessage: "1: undefined: undefinedVar",
			},
			"undefined in second line": {
				source:             "{ \n undefinedVar \n }",
				expectedErrMessage: "2: undefined: undefinedVar",
			},
			"invalid number literal": {
				source:             "a := 123a123",
				expectedErrMessage: "1: expected ';', found a123 (and 1 more errors)",
			},
			"invalid number literal in second line": {
				source:             "{ \n a := 123a123 \n }",
				expectedErrMessage: "2: expected ';', found a123 (and 1 more errors)",
			},
		}
		for name, testCase := range tests {
			t.Run(name, func(t *testing.T) {
				_, err := newInterpreterInstance(t).Eval(testCase.source)
				require.Error(t, err)
				assert.Equal(t, testCase.expectedErrMessage, err.Error())
			})
		}
	})
}

func TestExport(t *testing.T) {
	t.Run("should return error when name is not a Go identifier", func(t *testing.T) {
		invalidIdentifiers := []string{
			"", "1", "1a", "a-",
		}
		for _, invalidIdentifier := range invalidIdentifiers {
			t.Run(invalidIdentifier, func(t *testing.T) {
				// when
				err := newInterpreterInstance(t).Export(invalidIdentifier, "v")
				// then
				target := &interpreter.ErrInvalidIdentifier{}
				assert.ErrorAs(t, err, target)
			})
		}
	})

	t.Run("shout not return error when name is a Go identifier", func(t *testing.T) {
		validIdentifiers := []string{
			"a", "A", "ab", "AB", "a1", "abc", "a_",
		}
		for _, validIdentifier := range validIdentifiers {
			t.Run(validIdentifier, func(t *testing.T) {
				// when
				err := newInterpreterInstance(t).Export(validIdentifier, "v")
				// then
				assert.NoError(t, err)
			})
		}
	})

	t.Run("should return error when name clashes with imported package name", func(t *testing.T) {
		err := newInterpreterInstance(t).Export("fmt", "v")
		target := &interpreter.ErrInvalidIdentifier{}
		assert.ErrorAs(t, err, target)
		fmt.Println(target.Error())
	})

	t.Run("should return error when name clashes with command name", func(t *testing.T) {
		var allCommandNames = []string{"h", "help", "p", "pause", "r", "resume", "u", "undo", "n", "next"}

		for _, cmd := range allCommandNames {
			t.Run(cmd, func(t *testing.T) {
				// when
				err := newInterpreterInstance(t).Export(cmd, "v")
				// then
				target := &interpreter.ErrInvalidIdentifier{}
				assert.ErrorAs(t, err, target)
			})
		}
	})
}

func TestInstance_SetUpdate(t *testing.T) {
	t.Run("script should run pi.Update function", func(t *testing.T) {
		instance := newInterpreterInstance(t)

		var updateExecuted bool
		update := func() {
			updateExecuted = true
		}
		// when
		err := instance.SetUpdate(&update)
		// then
		require.NoError(t, err)
		_, err = instance.Eval("pi.Update()")
		require.NoError(t, err)
		assert.True(t, updateExecuted)
	})

	t.Run("script should replace pi.Update function", func(t *testing.T) {
		instance := newInterpreterInstance(t)

		var updateExecuted bool
		update := func() {
			updateExecuted = true
		}
		// when
		err := instance.SetUpdate(&update)
		// then
		require.NoError(t, err)
		_, err = instance.Eval("pi.Update = func() { }")
		require.NoError(t, err)
		assert.False(t, updateExecuted)
	})
}

func TestInstance_SetDraw(t *testing.T) {
	t.Run("script should run pi.Draw function", func(t *testing.T) {
		instance := newInterpreterInstance(t)

		var drawExecuted bool
		draw := func() {
			drawExecuted = true
		}
		// when
		err := instance.SetDraw(&draw)
		// then
		require.NoError(t, err)
		_, err = instance.Eval("pi.Draw()")
		require.NoError(t, err)
		assert.True(t, drawExecuted)
	})

	t.Run("script should replace pi.Draw function", func(t *testing.T) {
		instance := newInterpreterInstance(t)

		var drawExecuted bool
		draw := func() {
			drawExecuted = true
		}
		// when
		err := instance.SetDraw(&draw)
		// then
		require.NoError(t, err)
		_, err = instance.Eval("pi.Draw = func() { }")
		require.NoError(t, err)
		assert.False(t, drawExecuted)
	})

	t.Run("should print help", func(t *testing.T) {
		tests := map[string]string{
			"h":             "",
			" h":            "",
			"h ":            "",
			"help":          "",
			" help":         "",
			"help ":         "",
			"help pi":       "pi",
			"help pi ":      "pi",
			"help  pi ":     "pi",
			"help pi.Spr":   "pi.Spr",
			"help pi.Spr ":  "pi.Spr",
			"help  pi.Spr ": "pi.Spr",
			"h pi":          "pi",
			"h pi.Spr":      "pi.Spr",
		}
		for cmd, expectedTopic := range tests {
			var actualTopic string
			printHelp := func(topic string) error {
				actualTopic = topic
				return nil
			}

			instance, err := interpreter.New(printHelp)
			require.NoError(t, err)
			// when
			result, err := instance.Eval(cmd)
			// then
			require.NoError(t, err)
			assert.Equal(t, interpreter.HelpPrinted, result)
			assert.Equal(t, expectedTopic, actualTopic)
		}
	})
}

func newInterpreterInstance(t *testing.T) interpreter.Instance {
	instance, err := interpreter.New(noopPrintHelp)
	require.NoError(t, err)

	return instance
}

func noopPrintHelp(topic string) error { return nil }
