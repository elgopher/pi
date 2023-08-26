// (c) 2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elgopher/pi"
)

func TestIterators_Next(t *testing.T) {
	finishImmediately := func() bool {
		return false
	}

	t.Run("should run iterator", func(t *testing.T) {
		executionCount := 0
		i := func() bool {
			executionCount++
			return true
		}
		var iterators pi.Iterators
		iterators = append(iterators, i)
		// when
		_ = iterators.Next()
		// then
		assert.Equal(t, 1, executionCount)
	})

	t.Run("should remove finished iterator", func(t *testing.T) {
		var iterators pi.Iterators
		iterators = append(iterators, finishImmediately)
		// when
		iterators = iterators.Next()
		// then
		assert.Empty(t, iterators)
	})

	t.Run("should remove first and last iterator", func(t *testing.T) {
		neverFinish := func() bool {
			return true
		}

		var iterators pi.Iterators
		iterators = append(iterators, finishImmediately)
		iterators = append(iterators, neverFinish)
		iterators = append(iterators, finishImmediately)
		// when
		iterators = iterators.Next()
		// then
		require.Len(t, iterators, 1)
		hasNext := iterators[0]()
		assert.True(t, hasNext, "remaining iterator should be neverFinish")
	})
}
