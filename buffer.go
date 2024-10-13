package funky

// Buffer produces an iterator that will pre-fetch up to size values
// from its input so that when a value is requested, it may already
// be ready. This is useful for values that take some time to produce,
// such as those involving network requests. Ordering is maintained
// and evaluation is single-threaded, though it takes place off the
// calling goroutine. The buffer begins filling immediately.
//
// Example: remoteFiles := Buffer(httpRequests, 10)
func Buffer[T any](it *Iter[T], size uint32) *Iter[T] {
	// Results from the source iterator, will be closed when
	// the Iter stops, but may still contain previously generated
	// values which can still be read out with Next(). These
	// values are preserved so a buffered pipeline doesn't "leak".
	results := make(chan Elem[T], size-1)

	// We flip this to indicate to the producer goroutine that
	// stop has been called. There's no (problematic) race
	// condition here because the main thread only writes the
	// value and the goroutine only reads the value.
	stop := false

	// After we flip stop, we wait on this chan to close so
	// we know there are no further elements in flight. The
	// buffer could still hold up to size elements, but no
	// new elements will be added.
	done := make(chan interface{})

	go func() {
		defer close(results)
		defer close(done)

		for elem, valid := it.Next(); valid; elem, valid = it.Next() {
			results <- elem

			if stop {
				break
			}
		}
	}()

	return &Iter[T]{
		next: func() (Elem[T], bool) {
			select {
			case r := <-results:
				return r, true
			default:
				return DoneElem[T]()
			}
		},
		stop: func() {
			stop = true
			<-done
		},
	}
}
