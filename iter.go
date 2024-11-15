package funky

import (
	"iter"
	"sync"
)

// Elem is an internal iterator element that captures an
// error along with a value. In this way, errors are treated
// as first class values.
type Elem[T any] struct {
	val T
	err error
	// todo: this would let us send zero elems through channels
	ok bool
}

// DoneElem is a helper for returning an invalid iterator
// response, which is necessary once an iterator has been stopped
// or exhausted.
func DoneElem[T any]() (Elem[T], bool) {
	return Elem[T]{}, false
}

func ValElem[T any](val T) (Elem[T], bool) {
	return Elem[T]{val: val}, true
}

func ErrElem[T any](err error) (Elem[T], bool) {
	return Elem[T]{err: err}, true
}

// Iter is an internal iterator that supports all the fun
// stuff provided elsewhere.
type Iter[T any] struct {
	// next is a callback that returns the next element in the
	// iterator. This function may be called concurrently from
	// multiple goroutines, so it must handle its own
	// synchronization, if necessary.
	//
	// However, calls to this function will be
	// synchronized, so implementers can assume nothing will
	// change while the function runs, and that the iterator has
	// not been stopped.
	next func() (Elem[T], bool)

	// close provides a callback to be called when the iterator
	// is manually stopped. If there is no cleanup required, it
	// can be omitted.
	//
	// There is no need to synchronize this because the outer
	// implementation synchronizes reads and writes.
	close func()

	lock sync.RWMutex
}

// Next provides the next value from the iterator.
func (it *Iter[T]) Next() (elem Elem[T], valid bool) {
	// Short-circuit here to avoid the lock for a stopped
	// iterator (since they can't be un-stopped).
	if it.next == nil {
		return DoneElem[T]()
	}

	it.lock.RLock()
	defer it.lock.RUnlock()

	// Now that we're in the critical region, if the next()
	// callback hasn't been set to nil, we know that close
	// hasn't been called before we got here.

	if it.next != nil {
		return it.next()
	}

	return DoneElem[T]()
}

// Close tells the iterator to stop producing new values. Some
// iterators may allow Next to return additional values, such as
// when values have been buffered.
func (it *Iter[T]) Close() {
	// If we have already closed, we can just bail since
	// there is no way to un-close.
	if it.next == nil {
		return
	}

	// Call the close handler, if one was provided. This is
	// just intended for cleanup, so we don't have to worry
	// much about it. We call close before setting next to nil
	// to allow it to block if we've got anything in flight.
	if it.close != nil {
		it.close()
		it.close = nil
	}

	// Get the write lock because we're going to update the
	// state, but we want to wait for outstanding reads to finish.
	it.lock.Lock()
	defer it.lock.Unlock()

	// A stopped iterator has a nil next() callback, that's
	// how we define it, so this is all we have to do.
	it.next = nil
}

func (it *Iter[T]) ToChan() chan<- T {
	panic("not implemented")
}

func (it *Iter[T]) ToSeq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for elem, valid := it.Next(); valid; elem, valid = it.Next() {
			more := yield(elem.val)
			if !more {
				break
			}
		}
	}
}

// ToSlice turns the first n elements in the iterator into a slice.
// If there are fewer than n elements in the iterator, then the
// length of the resulting slice will be less than n.
func (it *Iter[T]) ToSlice(n uint64) []T {
	var elems []T
	for elem, valid := it.Next(); valid && uint64(len(elems)) < n; elem, valid = it.Next() {
		elems = append(elems, elem.val)
	}

	return elems
}

func FromChan[T any](c chan<- T) *Iter[T] {
	panic("not implemented")
}

func FromSeq[T any](s iter.Seq[T]) *Iter[T] {
	next, stop := iter.Pull(s)
	mut := sync.Mutex{}

	return &Iter[T]{
		next: func() (Elem[T], bool) {
			mut.Lock()
			defer mut.Unlock()

			val, valid := next()
			if !valid {
				return DoneElem[T]()
			}

			return ValElem(val)
		},
		close: func() {
			mut.Lock()
			defer mut.Unlock()

			stop()
		},
	}
}

func FromSlice[T any](s []T) *Iter[T] {
	i := 0
	mut := sync.Mutex{}

	return &Iter[T]{
		next: func() (Elem[T], bool) {
			// We can check before interacting with the lock
			// because once we're done, we're done, so there's
			// no dangerous race condition.
			if i >= len(s) {
				return DoneElem[T]()
			}

			mut.Lock()
			k := i
			i++
			mut.Unlock()

			if k >= len(s) {
				return DoneElem[T]()
			}

			return Elem[T]{
				val: s[k],
				err: nil,
			}, true
		},
		close: func() {
			mut.Lock()
			i = len(s)
			mut.Unlock()
		},
	}
}

func FromVals[T any](vals ...T) *Iter[T] {
	return FromSlice(vals)
}
