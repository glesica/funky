package funky

import (
	"iter"
)

type Iter[T any] iter.Seq[T]

// We could wrap values when we convert to a different iterator
// types. In fact, we could wrap everything once it's "inside"
// the type framework and then provide something like "Unwrap"
// to extract actual values, and handling errors in various ways.

func (p Iter[T]) Chan() <-chan T {
	return nil
}

func (p Iter[T]) Slice() []T {
	return nil
}

// OnError applies the given function to each error produced by
// the Iter. This can be used along with Seq() to, for example,
// log errors to a side channel so that they don't disappear
// entirely. The value is also passed in case it is meaningful.
//
// We're not doing it this way right now, but we could?
func (p Iter[T]) OnError(handleErr func(T, error)) Iter[T] {
	return nil
}

// todo: additional methods to implement correctly are listed below
// Elems() (provides Elem[T] instances)

// todo: change the transformation methods to be plain functions

// todo: add a method to copy values into a pair but keep one iterator
// {1, 2, 3}.<whatever>() -> {Pair(1, 1), Pair(2, 2), Pair(3, 3)}

// todo: add a method to pair off adjacent values
// do we still want this despite having Chunk available?
// {1, 2, 3}.<whatever>() -> {Pair(1, 2), Pair(3, 0)}

// todo: add a ticker that calls a function with the value but passes it on
// {1, 2, 3}.<whatever>((v) -> print(v)) -> {1, 2, 3}
// this is already easily possible with apply, it's just a convenience

func (p Iter[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range p {
			ok := yield(v)
			if !ok {
				return
			}
		}
	}
}

// Split creates two iterators that will produce the same values
// as the target.
func (p Iter[T]) Split() (Iter[T], Iter[T]) {
	return nil, nil
}

// Buffered reads ahead in the iterator by up to n elements and returns
// an iterator that pulls from this buffer.
//
// This is useful, for example, if each element in the iterator
// takes a lot of time to produce.
func (p Iter[T]) Buffered(n int) Iter[T] {
	return nil
}

// Parallel returns an iterator that produces the same values as
// the target, but does so in parallel using up to n workers. The
// values will not necessarily appear in the same order in the
// resulting iterator.
//
// Note that the target must be thread-safe!
func (p Iter[T]) Parallel(n int) Iter[T] {
	return nil
}
