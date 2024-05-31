package funky

import (
	"errors"
	"fmt"
)

// Chunk breaks an iterator into slices of the given size. If there
// are errors, they are joined together and passed along with the
// assembled chunk.
//
// The last chunk may be shorter than size if the number of values
// produced by the underlying iterator is indivisible by size.
//
// This function can be composed with, for example, Each and Parallel,
// to improve the efficiency of parallel computations by performing
// operations on several inputs in one goroutine.
func Chunk[T any](iterator Func[T], size int) Func[[]T] {
	return func() ([]T, bool, error) {
		next := make([]T, size)
		errs := make([]error, size)

		found := 0
		for i := 0; i < size; i++ {
			value, more, err := iterator()
			if !more {
				// There are no additional values, so we're done,
				// this means the previous chunk was full and the
				// number of values was divisible by size.
				if i == 0 {
					return Done[[]T](err)
				}

				// The number of values was not divisible by size,
				// so we just allow an incomplete chunk.
				break
			}

			found++

			next[i] = value
			errs[i] = err
		}

		return Value(next[:found], errors.Join(errs[:found]...))
	}
}

// Unchunk reverses the Chunk operation and properly handles the
// fact that the final chunk may have a different length than
// the others. In fact, all chunks may be arbitrary lengths.
//
// For now, the entire error for a given chunk is passed with
// each of its values. In the future, this function may attempt
// to pair wrapped errors with their corresponding value.
//
// For example:
//
//	Unchunk([]{{1, 2}, {3, 4}}) -> {1, 2, 3, 4}
func Unchunk[T any](iterator Func[[]T]) Func[T] {
	var next []T
	var errs error

	return func() (T, bool, error) {
		if len(next) == 0 {
			value, more, err := iterator()
			if !more {
				next = nil
				errs = nil

				return Done[T](err)
			}

			next = value
			errs = err
		}

		value := next[0]
		next = next[1:]

		return Value(value, errs)
	}
}

// Concat creates a single iterator that will produce the values
// from each of the provided iterators, in order.
//
// Note, if an infinite iterator is passed in any but the final
// position, then iterators that follow it will never be used.
//
// For example:
//
//	Concat({1, 2}, {3, 4}) -> {1, 2, 3, 4}
func Concat[T any](iterators ...Func[T]) Func[T] {
	return func() (T, bool, error) {
		for {
			if len(iterators) == 0 {
				iterators = nil
				return Done[T](nil)
			}

			value, more, err := iterators[0]()
			if !more {
				iterators = iterators[1:]
				continue
			}

			return value, more, err
		}
	}
}

// Each applies the given function to each element produced by the
// iterator.
//
// If the function returns Stop as its error, iteration will cease
// and no further values will be processed.
//
// Note that this function is not thread-safe unless the iterator
// implements its own synchronization.
//
// For example:
//
//	Each({1, 2, 3}, (v) -> print(v))
func Each[T any](iterator Func[T], f func(T) error) error {
	for {
		value, more, err := iterator()
		if err != nil {
			return fmt.Errorf("each: iterator error: %w", err)
		}

		if !more {
			break
		}

		err = f(value)
		if errors.Is(err, Stop) {
			break
		}
		if err != nil {
			return fmt.Errorf("each: function error: %w", err)
		}
	}

	return nil
}

// Map returns an iterator that produces the values from the given
// iterator, with the given callback applied to them.
//
// If the callback returns Stop as its error, iteration will cease
// and no more values will be produced.
//
// Note that this function is not thread-safe, even if the iterator
// implements its own synchronization.
//
// todo: make map() thread-safe for thread-safe iterators
// todo: must solve partial reads problem for Stop to work
func Map[T, R any](iterator Func[T], mapper Mapper[T, R]) Func[R] {
	stopped := false

	return func() (R, bool, error) {
		if stopped {
			return Done[R](nil)
		}

		inValue, more, err := iterator()
		if err != nil {
			return Done[R](fmt.Errorf("map: iterator error: %w", err))
		}

		if !more {
			return Done[R](nil)
		}

		outValue, err := mapper(inValue)
		if errors.Is(err, Stop) {
			return Done[R](nil)
		}
		if err != nil {
			return Done[R](fmt.Errorf("map: callback error: %w", err))
		}

		return outValue, more, nil
	}
}

// Reduce runs a reduce operation on the iterator, calling the
// reducer function repeatedly until the iterator has been exhausted
// and keeping an accumulated value around until it is returned.
//
// todo: must solve partial reads problem for this to work with infinite iterators
func Reduce[T any, R any](iterator Func[T], reducer Reducer[T, R]) (R, error) {
	return *new(R), nil
}

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
func Zip[L any, R any](left Func[L], right Func[R]) Func[Pair[L, R]] {
	return nil
}
