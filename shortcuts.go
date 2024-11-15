package funky

import (
	"errors"
	"sync"
)

// A Pair is just two values, of potentially different types, that
// go together somehow.
type Pair[L any, R any] struct {
	Left  L
	Right R
}

type Predicate[T any] func(T) bool

// NoError simply skips any elements that include an
// error value.
func NoError[T any](it *Iter[T]) *Iter[T] {
	return &Iter[T]{
		next: func() (Elem[T], bool) {
			for {
				elem, valid := it.Next()
				if !valid {
					return DoneElem[T]()
				}

				if elem.err == nil {
					return elem, true
				}
			}
		},
		close: func() {
			it = nil
		},
	}
}

func Take[T any](it *Iter[T], n uint64) *Iter[T] {
	var lock sync.Mutex
	var count uint64

	return &Iter[T]{
		next: func() (Elem[T], bool) {
			if count >= n {
				return DoneElem[T]()
			}

			lock.Lock()
			if count >= n {
				lock.Unlock()
				return DoneElem[T]()
			}
			count++
			lock.Unlock()

			return it.Next()
		},
		close: func() {
			it = nil
		},
	}
}

func Where[T any](it *Iter[T], keep Predicate[T]) *Iter[T] {
	return &Iter[T]{
		next: func() (Elem[T], bool) {
			for {
				elem, valid := it.Next()
				if !valid {
					return DoneElem[T]()
				}

				if elem.err != nil {
					return elem, true
				}

				if keep(elem.val) {
					return elem, true
				}
			}
		},
		close: func() {
			it = nil
		},
	}
}

// While delivers values from the given iterator as long as the given
// predicate returns true. Note that the values are only guaranteed
// to be tested in order if only requested from a single goroutine.
// For example, elements that would result in {true, true, false}
// when tested in a single goroutine could be interpreted as
// {true, false true} when called from multiple goroutines.
//
// Example:
//
// While({1, 2, 3, 2, 1}, x -> x < 3) -> {1, 2}
func While[T any](it *Iter[T], still Predicate[T]) *Iter[T] {
	var done bool
	var lock sync.Mutex

	return &Iter[T]{
		next: func() (Elem[T], bool) {
			// Short circuit if we're already done. No need for a
			// lock here because once we're done, we're done forever.
			// We check again after we acquire a read lock.
			if done {
				return DoneElem[T]()
			}

			elem, valid := it.Next()
			if !valid {
				// We don't need synchronization here because the done
				// boolean is just a shortcut. If some calls miss it
				// and end up calling Next again, it's fine, because
				// Next will continue to return invalid values.
				done = true
				return DoneElem[T]()
			}

			// If we got an error, we can't test it, so we assume the
			// caller wants to see it whether it would have passed or not.
			if elem.err != nil {
				return elem, true
			}

			// Only test one value at a time. There's no global
			// ordering of the values themselves if they are requested
			// from multiple goroutines, but we can impose an order
			// once they are delivered and ready to be tested. The
			// first one tested that fails will put us into a done
			// state and all subsequent values will return as done,
			// regardless of whether they would have passed.
			lock.Lock()
			defer lock.Unlock()

			if done {
				return DoneElem[T]()
			}

			if still(elem.val) {
				return elem, true
			}

			// The value failed, so we want to bail on all future calls.
			done = true

			return DoneElem[T]()
		},
	}
}

// Zip creates an iterator from two other iterators that produces
// as its elements pairs of values, one from each of the original
// iterators. If the iterators produce different number of values,
// the resulting iterator will end once one of the inputs ends.
//
// For example, Zip({1, 2}, {"a", "b"}) would produce two values:
//
//  1. {Left: 1, Right: "a"}
//  2. {Left: 2, Right: "b"}
//
// If one iterator ends early, like Zip({1}, {"a", "b"}), then
// iteration will terminate and the last value will be ignored.
//
//  1. {Left: 1, Right: "a"}
//
// In this case, consuming additional values will continue to
// drain the non-empty iterator, but the zipped iterator will
// still signal that it is empty.
//
// Note that accesses to the left and right iterators are not
// synchronized, so if another goroutine is using one or both
// of them, the values produced by Zip may be out of order.
func Zip[L, R any](left *Iter[L], right *Iter[R]) *Iter[Pair[L, R]] {
	return &Iter[Pair[L, R]]{
		next: func() (Elem[Pair[L, R]], bool) {
			leftElem, leftValid := left.Next()
			rightElem, rightValid := right.Next()

			if !leftValid || !rightValid {
				return DoneElem[Pair[L, R]]()
			}

			err := errors.Join(leftElem.err, rightElem.err)
			if err != nil {
				return ErrElem[Pair[L, R]](err)
			}

			return ValElem(Pair[L, R]{leftElem.val, rightElem.val})
		},
		close: func() {
			left = nil
			right = nil
		},
	}
}
