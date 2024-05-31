package funky

import (
	"fmt"
	"sync"
)

// A Func is an iterator implemented as a function. Each call to the
// iterator provides a value, a boolean that will
// be false after the iterator is exhausted and true otherwise, and
// any error that occurred during production of the value.
//
// Iterators are not required, or assumed, to be thread-safe. If you
// want to use an iterator from multiple Goroutines, or if you want
// to use the Parallel method, you should arrange for synchronization.
//
// Note that an iterator may continue to provide values after an
// error has been returned, but it is not required to do so. Either
// way, it must communicate this through the boolean return value.
//
// Once false has been returned for the boolean, subsequent calls
// to the iterator must be allowed, and must also return false with
// a nil error. The Done() helper function can be used here.
type Func[T any] func() (T, bool, error)

// Apply transforms each input value but keeps the type the same.
//
// For example:
//
//	{1, 2, 3}.Apply(x -> 2 * x) -> {2, 4, 6}
func (p Func[T]) Apply(op Applier[T]) Func[T] {
	return func() (T, bool, error) {
		value, more, err := p()
		if !more {
			return Done[T](err)
		}
		if err != nil {
			return Value(value, err)
		}

		newValue, err := op(value)

		return Value(newValue, err)
	}
}

// todo: add a method to copy values into a pair but keep one iterator
// {1, 2, 3}.<whatever>() -> {Pair(1, 1), Pair(2, 2), Pair(3, 3)}

// todo: add a method to pair off adjacent values
// do we still want this despite having Chunk available?
// {1, 2, 3}.<whatever>() -> {Pair(1, 2), Pair(3, 0)}

// todo: add a ticker that calls a function with the value but passes it on
// {1, 2, 3}.<whatever>((v) -> print(v)) -> {1, 2, 3}
// this is already easily possible with apply, it's just a convenience

// Split creates two iterators that will produce the same values
// as the target. The iterators are synchronized, so only one
// value can be produced by the first until the second has also
// produced the same value, and vice versa.
func (p Func[T]) Split() (Func[T], Func[T]) {
	mut := &sync.Mutex{}

	c1 := make(chan errorPair[T], 1)
	c2 := make(chan errorPair[T], 1)

	return p.splitBranch(mut, c1, c2), p.splitBranch(mut, c2, c1)
}

// splitBranch arranges for one of a pair of split iterators to
// fetch new values when required. Note that the two channels MUST
// have a buffer of 1 or the iterators will deadlock if they are
// both accessed from the same goroutine.
func (p Func[T]) splitBranch(mut *sync.Mutex, ours chan errorPair[T], theirs chan errorPair[T]) Func[T] {
	return func() (T, bool, error) {
		select {
		case pair, more := <-ours:
			// We had a value waiting for us, fetched by them, so
			// we can go ahead and provide it. If our channel has
			// been closed, we know there are no further values.
			if !more {
				return Done[T](nil)
			}
			return Value(pair.value, pair.err)
		default:
			if mut.TryLock() {
				// We got the lock, so we are responsible for fetching
				// and providing the next value.
				defer mut.Unlock()

				value, more, err := p()
				if !more {
					// There are no further values, so we close both
					// channels since we will be the last one in this
					// position.
					close(ours)
					close(theirs)
				} else {
					// There was a new value available, so prep it and
					// send it over both channels!
					p := errorPair[T]{value: value, err: err}
					ours <- p
					theirs <- p
				}
			}

			// We didn't get the lock, or we did and sent ourselves
			// a new value (or close), so we know a value or
			// channel close will be waiting for us, so just
			// grab the next value and handle it.
			pair, more := <-ours
			if !more {
				return Done[T](nil)
			}
			return Value(pair.value, pair.err)
		}
	}
}

// Buffered reads ahead in the iterator by up to n elements and returns
// an iterator that pulls from this buffer.
//
// This is useful, for example, if each element in the iterator
// takes a lot of time to produce.
func (p Func[T]) Buffered(n int) Func[T] {
	c := p.bufferedChan(n)
	return Chan(c)
}

// Parallel returns an iterator that produces the same values as
// the target, but does so in parallel using up to n workers. The
// values will not necessarily appear in the same order in the
// resulting iterator.
//
// Note that the target must be thread-safe!
func (p Func[T]) Parallel(n int) Func[T] {
	// We don't buffer this because the same effect can be achieved
	// by following Parallel with a call to Buffered, which will
	// effectively allow Parallel to keep going even if values are
	// not being consumed.
	out := make(chan errorPair[T])

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			for {
				value, more, err := p()
				if !more {
					wg.Done()
					return
				}

				out <- errorPair[T]{
					value: value,
					err:   err,
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return func() (T, bool, error) {
		pair, more := <-out
		if !more {
			return Done[T](nil)
		}

		return Value(pair.value, pair.err)
	}
}

// DropN returns an iterator that produces the same values as the
// target, but skips the first n values.
func (p Func[T]) DropN(n int) Func[T] {
	count := 0
	return func() (T, bool, error) {
		for {
			value, more, err := p()
			if !more {
				return Done[T](err)
			}

			if count >= n {
				return value, more, err
			}

			count++
		}
	}
}

// Drop returns an iterator that skips values for which f returns
// true, then, once it has encountered a value for which f returns
// false, all following values are produced.
//
// For example:
//
//	{1, 0, 2}.Drop(x -> x > 0) -> {0, 2}
func (p Func[T]) Drop(f Predicate[T]) Func[T] {
	dropDone := false
	return func() (T, bool, error) {
		for {
			value, more, err := p()
			if !more {
				return Done[T](err)
			}

			if dropDone {
				return Value(value, err)
			}

			drop, err := f(value)
			if err != nil {
				return Value(value, err)
			}

			if drop {
				continue
			} else {
				dropDone = true
			}

			return Value(value, err)
		}
	}
}

// TakeN returns an iterator that produces only the first n values
// from the target.
func (p Func[T]) TakeN(n int) Func[T] {
	count := 0
	takeDone := false
	return func() (T, bool, error) {
		if takeDone {
			return Done[T](nil)
		}

		if count >= n {
			takeDone = true
			return Done[T](nil)
		}

		value, more, err := p()
		if !more {
			takeDone = true
			return Done[T](err)
		}

		count++

		return value, more, err
	}
}

// Take returns an iterator that produces values from the target
// for as long as the callback returns true, then stops. Therefore,
// calling f on any of the values produced will return true.
func (p Func[T]) Take(f Predicate[T]) Func[T] {
	done := false
	return func() (T, bool, error) {
		if done {
			return Done[T](nil)
		}

		value, more, err := p()
		if !more {
			done = true
			return Done[T](nil)
		}

		take, err := f(value)
		if err != nil {
			return Value(value, err)
		}

		if !take {
			done = true
			return Done[T](nil)
		}

		return Value(value, err)
	}
}

// Where returns an iterator that produces all the values from
// the target for which f returns true.
//
// For example:
//
//	{1, 0, 2}.Where(x -> x > 0) -> {1, 2}
func (p Func[T]) Where(f Predicate[T]) Func[T] {
	return func() (T, bool, error) {
		for {
			value, more, err := p()
			if !more {
				return Done[T](err)
			}

			if err != nil {
				return Value(value, err)
			}

			ok, err := f(value)
			if err != nil {
				return Value(value, err)
			}

			if ok {
				return Value(value, nil)
			}
		}
	}
}

// Chan returns a channel that produces one value at a time until
// it has exhausted the iterator, at which point the channel is
// closed. The channel must be fully consumed, or it will never
// be closed.
//
// If the iterator produces an error, it will be ignored and the
// channel will be closed, as though the iterator had been
// exhausted.
func (p Func[T]) Chan() <-chan T {
	return p.bufferedChan(0)
}

// Coalesce accumulates the values produced by the target into
// a single value of the same type. Think of it like Reduce, but
// the accumulator has the same type as the inputs.
//
// For example:
//
//	{1, 2, 3}.Coalesce(a, v -> a + v) -> 6
func (p Func[T]) Coalesce(c Coalescer[T]) (T, error) {
	index := -1
	acc := *new(T)
	for {
		index++

		value, more, err := p()
		if err != nil {
			return acc, fmt.Errorf("coalesce: iterator error on index %d: %w", index, err)
		}

		if !more {
			break
		}

		acc, err = c(acc, value)
		if err != nil {
			return acc, fmt.Errorf("coalesce: function error on index %d: %w", index, err)
		}
	}

	return acc, nil
}

// Slice resolves the entire iterator to a slice of the value type.
// If an error is encountered during iteration, the slice up to that
// point and an error will be returned.
//
// Note that this will result in a crash in the iterator is infinite
// or otherwise too big for all of its elements to fit in memory!
func (p Func[T]) Slice() ([]T, error) {
	var values []T
	for {
		value, more, err := p()
		if err != nil {
			return values, fmt.Errorf("slice: iterator error on index %d: %w", len(values), err)
		}

		if !more {
			break
		}

		values = append(values, value)
	}

	return values, nil
}

// bufferedChan creates a buffered channel that will produce the
// values produced by the iterator. It is used behind the scenes.
func (p Func[T]) bufferedChan(n int) <-chan T {
	// Reduce the buffer size by 1 since the actual buffer for
	// our purposes will be the values in the channel buffer plus
	// one additional value that is blocked.
	if n > 0 {
		n--
	}

	c := make(chan T, n)

	go func() {
		defer close(c)

		for {
			value, more, err := p()
			// todo: pass errors out somehow, maybe take a callback, or wrap value and error?
			if err != nil {
				break
			}

			if !more {
				break
			}

			c <- value
		}
	}()

	return c
}
