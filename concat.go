package funky

import "sync"

// Concat creates a single iterator that will produce the values
// from each of the provided iterators, in order.
//
// Note, if an infinite iterator is passed in any but the final
// position, then iterators that follow it will never be used.
//
// For example (in pseudocode):
//
//	Concat({1, 2}, {3, 4}) -> {1, 2, 3, 4}
func Concat[T any](its ...*Iter[T]) *Iter[T] {
	if len(its) == 0 {
		return &Iter[T]{}
	}

	index := 0
	mut := sync.Mutex{}

	return &Iter[T]{
		next: func() (Elem[T], bool) {
			for {
				myIndex := index

				if myIndex >= len(its) {
					return DoneElem[T]()
				}

				elem, valid := its[myIndex].Next()
				if !valid {
					// Try to get the mutex. If we get it, then we
					// need to update myIndex. Otherwise, we know
					// someone else is already doing that, so we can
					// just try again immediately. Worst case scenario
					// here is that we busy wait for a couple iterations,
					// but updating the index is fast so we won't wait
					// very long.

					if !mut.TryLock() {
						continue
					}

					// Since it's possible that someone hit the lock
					// after we had already updated index and released
					// it, we still need to make sure the value hasn't
					// changed before we update it. If it is different,
					// then we know someone got to it before us.

					if myIndex == index {
						index++
					}

					mut.Unlock()
					continue
				}

				return elem, true
			}
		},
		close: func() {
			its = nil
		},
	}
}
