// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/internal/state"
)

const stateName = "stateName"

func TestStateSave(t *testing.T) {
	t.Cleanup(state.Cleanup)

	testStateName(t, func(name string) error {
		return pi.StateSave(name, "")
	})

	t.Run("should store data", func(t *testing.T) {
		testStateSave(t, "string", "value")
		testStateSave(t, "int", 1)
		testStateSave(t, "float64", 2.0)
		testStateSave(t, "slice", []string{"0", "1"})
		testStateSave(t, "struct", struct{ String string }{String: "value"})
	})

	t.Run("should not panic when trying to store 10MB", func(t *testing.T) {
		t.Cleanup(state.Cleanup)
		assert.NotPanics(t, func() {
			_ = pi.StateSave("too-big", strings.Repeat(" ", 10*1024*1024))
		})
	})

	t.Run("should return error when trying to save unmarshalable struct", func(t *testing.T) {
		err := pi.StateSave("unmarshalable", unmarshalableStruct{})
		assert.ErrorIs(t, err, pi.ErrStateMarshalFailed)
	})
}

type unmarshalableStruct struct{}

func (s unmarshalableStruct) MarshalJSON() ([]byte, error) {
	return nil, errors.New("json marshaling failed")
}

func testStateSave[T any](t *testing.T, testName string, expected T) {
	t.Run(testName, func(t *testing.T) {
		// when
		require.NoError(t, pi.StateSave(stateName, expected))
		// then
		var actual T
		require.NoError(t, pi.StateLoad(stateName, &actual))
		assert.Equal(t, expected, actual)
	})
}

func TestStateDelete(t *testing.T) {
	t.Cleanup(state.Cleanup)

	testStateName(t, pi.StateDelete)

	t.Run("should permanently delete state", func(t *testing.T) {
		require.NoError(t, pi.StateSave(stateName, "value"))
		// when
		err := pi.StateDelete(stateName)
		// then
		require.NoError(t, err)
		var out string
		err = pi.StateLoad(stateName, &out)
		assert.ErrorIs(t, err, pi.ErrStateNotFound)
	})

	t.Run("should not return error when state is not found", func(t *testing.T) {
		err := pi.StateDelete("missing")
		assert.NoError(t, err)
	})
}

func TestStateLoad(t *testing.T) {
	t.Cleanup(state.Cleanup)

	testStateName(t, func(name string) error {
		_ = pi.StateSave(name, "")
		var out string
		return pi.StateLoad(name, &out)
	})

	t.Run("should return error when state does not exist", func(t *testing.T) {
		var out string
		err := pi.StateLoad("missing", &out)
		assert.ErrorIs(t, err, pi.ErrStateNotFound)
	})

	t.Run("should return error when out is nil", func(t *testing.T) {
		require.NoError(t, pi.StateSave(stateName, ""))
		var nilOut *string
		err := pi.StateLoad(stateName, nilOut)
		assert.ErrorIs(t, err, pi.ErrNilStateOutput)
	})

	t.Run("should return error for incompatible types", func(t *testing.T) {
		str := ""
		require.NoError(t, pi.StateSave(stateName, str))
		var number int
		err := pi.StateLoad(stateName, &number)
		assert.ErrorIs(t, err, pi.ErrStateUnmarshalFailed)
	})

	t.Run("should not return error when types are different but compatible", func(t *testing.T) {
		float := 1.0
		require.NoError(t, pi.StateSave(stateName, float))
		var integer int
		err := pi.StateLoad(stateName, &integer)
		require.NoError(t, err)
		assert.Equal(t, 1, integer)
	})
}

func TestStates(t *testing.T) {
	t.Cleanup(state.Cleanup)

	t.Run("should return empty states", func(t *testing.T) {
		state.Cleanup()

		states, err := pi.States()
		require.NoError(t, err)

		assert.Empty(t, states)
	})

	t.Run("should return previously stored states", func(t *testing.T) {
		require.NoError(t, pi.StateSave("state1", "1"))
		require.NoError(t, pi.StateSave("state2", "2"))
		states, err := pi.States()
		require.NoError(t, err)
		assert.ElementsMatch(t, states, []string{"state1", "state2"})
	})
}

func testStateName(t *testing.T, f func(name string) error) {
	t.Run("should return error for invalid name", func(t *testing.T) {
		invalidNames := []string{"", strings.Repeat("a", 33), "/", "a/", "\\", "a\\"}
		for _, name := range invalidNames {
			t.Run(name, func(t *testing.T) {
				err := f(name)
				assert.ErrorIs(t, err, pi.ErrInvalidStateName)
			})
		}
	})

	t.Run("should not return error for valid name", func(t *testing.T) {
		validNames := []string{"-", " ", strings.Repeat("a", 32), "."}
		for _, name := range validNames {
			t.Run(name, func(t *testing.T) {
				err := f(name)
				assert.NoError(t, err)
			})
		}
	})
}
