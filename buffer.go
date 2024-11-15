package funky

import "sync"

// Buffer produces an iterator that will pre-fetch up to size values
// from its input so that when a value is requested, it may already
// be ready. This is useful for values that take some time to produce,
// such as those involving network requests. Ordering is maintained
// and evaluation is single-threaded, though it takes place off the
// calling goroutine. The buffer begins filling immediately.
//
// Example: remoteFiles := Buffer(httpRequests, 10)
func Buffer[T any](it *Iter[T], size uint32) *Iter[T] {
	// Elements from the source iterator, will be closed when
	// the Iter stops, but may still contain previously generated
	// values which can still be read out with Next(). These
	// values are preserved so a buffered pipeline doesn't "leak".
	elements := make(chan Elem[T], size)

	// Requests for new elements to be buffered after the buffer
	// has been filled initially. These correspond to calls to
	// Next(), first trigger a back-fill of the buffer, then wait
	// for the first available ready element.
	requests := make(chan interface{}, 1)

	// We use this to keep Close() from returning until all
	// previously buffered values have been retrieved with Next().
	var elemGroup sync.WaitGroup

	loadOne := func() {
		elem, valid := it.Next()
		if !valid {
			return
		}

		elemGroup.Add(1)
		elements <- elem
	}

	// Fetch elements, one at a time, from the source iterator
	// in a goroutine so we don't block creating the buffered
	// iterator itself. This is the only place we can add to
	// the elements channel.
	go func() {
		defer close(elements)

		for i := uint32(0); i < size; i++ {
			// We don't rely on accessing this here, but we do
			// want to bail if we're closed before the initial
			// buffer is even filled.
			if requests == nil {
				return
			}
			loadOne()
		}

		for _ = range requests {
			loadOne()
		}
	}()

	// We mess around with the requests channel in two places,
	// in next() and stop(), so we guard it to prevent race
	// conditions.
	var reqLock sync.Mutex

	return &Iter[T]{
		next: func() (Elem[T], bool) {
			elem, more := <-elements
			if !more {
				return DoneElem[T]()
			}

			// We took an element, so decrement the wait group
			// to move Close() closer to being able to return.
			elemGroup.Done()

			// Ask for another element to be buffered
			// since we just took one, but only if we
			// haven't been stopped.
			go func() {
				if requests == nil {
					return
				}

				reqLock.Lock()
				defer reqLock.Unlock()

				if requests != nil {
					requests <- nil
				}
			}()

			return elem, true
		},
		close: func() {
			if requests != nil {
				reqLock.Lock()
				defer reqLock.Unlock()

				if requests != nil {
					close(requests)
					requests = nil
				}
			}

			// Wait for all outstanding elements to be retrieved
			// with Next().
			elemGroup.Wait()
		},
	}
}
