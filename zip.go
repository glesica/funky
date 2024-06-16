package funky

import "iter"

// Zip creates an iterator from two other iterators that produces
// as its elements pairs of values, one from each of the original
// iterators. If the iterators produce different number of values,
// the "extras" will have the other field set to the type's zero
// value.
//
// For example, Zip({1, 2}, {"a", "b"}) would produce two values:
//
//  1. {Left: 1, Right: "a"}
//  2. {Left: 2, Right: "b"}
func Zip[L any, R any](left Iter[L], right Iter[R]) Iter[Pair[L, R]] {
	return func(yield func(Pair[L, R]) bool) {
		nextLeft, stopLeft := iter.Pull(left.Seq())
		nextRight, stopRight := iter.Pull(right.Seq())

		for {
			leftVal, leftValid := nextLeft()
			rightVal, rightValid := nextRight()

			if !leftValid && !rightValid {
				return
			}

			more := yield(Pair[L, R]{
				Left:  leftVal,
				Right: rightVal,
			})
			if !more {
				stopLeft()
				stopRight()
				return
			}
		}
	}
}
