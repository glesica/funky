package funky

import "sync/atomic"

// Parallel applies a bound on the number of parallel called to
// Next() will be run in parallel, even if the calls originate
// from different goroutines. Passing 0 for n will cause
// execution to occur serially. Passing any other value will
// result in that number of additional simultaneous executions.
func Parallel[T any](it *Iter[T], n uint32) *Iter[T] {
	queue := make(chan interface{}, n+1)
	done := &atomic.Bool{}

	return &Iter[T]{
		next: func() (Elem[T], bool) {
			if done.Load() {
				return DoneElem[T]()
			}

			// Claim a spot in the pool
			queue <- nil

			result := make(chan Elem[T])

			go func() {
				defer close(result)

				elem, valid := it.Next()
				if !valid {
					return
				}

				result <- elem
			}()

			elem, ok := <-result
			if !ok {
				return DoneElem[T]()
			}

			// Free up our spot in the pool
			<-queue

			return elem, true
		},
		close: func() {
			done.Store(true)
		},
	}
}
