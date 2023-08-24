// (c) 2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

// Iterator is a function that performs the next iteration step. Function
// returns true if there are still steps to be performed. Function
// returns false if the iterator has finished.
type Iterator func() bool

// Iterators is a slice of iterators that can be run in bulk.
type Iterators []Iterator

// Next runs the next step of all iterators. Deletes those that have ended.
// The new iterator slice is returned from the function.
func (r Iterators) Next() Iterators {
	return sliceDelete(r, func(i Iterator) bool {
		return !i()
	})
}

// Following function is taken from slices packages, from Go 1.21 stdlib: slices.DeleteFunc. Please use the original version when Go in Pi is 1.21.
func sliceDelete(iterators Iterators, del func(i Iterator) bool) Iterators {
	// Don't start copying elements until we find one to delete.
	for i, v := range iterators {
		if del(v) {
			j := i
			for i++; i < len(iterators); i++ {
				v = iterators[i]
				if !del(v) {
					iterators[j] = v
					j++
				}
			}
			return iterators[:j]
		}
	}
	return iterators
}

// Sequence is hard to debug. I cant put a breakpoint!
// TODO REMOVE IT. This thing is too much of an abstraction!
func Sequence(iterators ...Iterator) Iterator {
	iteratorIdx := 0

	return func() bool {
		if len(iterators) == iteratorIdx {
			return false
		}

		hasNext := iterators[iteratorIdx]()
		if !hasNext {
			iteratorIdx++
		}

		return len(iterators) != iteratorIdx
	}
}
